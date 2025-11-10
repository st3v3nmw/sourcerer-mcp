package index

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"

	"github.com/philippgille/chromem-go"
	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

const (
	minSimilarity = 0.3
	maxResults    = 30
)

type Index struct {
	workspaceRoot string
	collection    *chromem.Collection

	cache   map[string]int64 // filePath -> max parsedAt timestamp
	cacheMu sync.RWMutex

	initOnce sync.Once
	initErr  error
}

func New(ctx context.Context, workspaceRoot string) (*Index, error) {
	idx := &Index{
		workspaceRoot: workspaceRoot,
		cache:         map[string]int64{},
	}

	go idx.ensureInitialized(ctx)

	return idx, nil
}

func (idx *Index) ensureInitialized(ctx context.Context) error {
	idx.initOnce.Do(func() {
		db, err := chromem.NewPersistentDB(".sourcerer/db", false)
		if err != nil {
			idx.initErr = fmt.Errorf("failed to create vector db: %w", err)
			return
		}

		collection, err := db.GetOrCreateCollection("code-chunks", nil, nil)
		if err != nil {
			idx.initErr = fmt.Errorf("failed to create vector db collection: %w", err)
			return
		}

		idx.collection = collection
		idx.loadCache(ctx)
	})

	return idx.initErr
}

func (idx *Index) loadCache(ctx context.Context) {
	idx.cacheMu.Lock()
	defer idx.cacheMu.Unlock()

	docs, err := idx.collection.ListDocumentsShallow(ctx)
	if err != nil {
		return
	}

	fileMaxParsed := make(map[string]int64)
	for _, doc := range docs {
		filePath := doc.Metadata["file"]
		_, exists := fileMaxParsed[filePath]
		if exists {
			continue
		}

		parsedAt, err := strconv.ParseInt(doc.Metadata["parsedAt"], 10, 64)
		if err != nil {
			continue
		}

		fileMaxParsed[filePath] = parsedAt
	}

	idx.cache = fileMaxParsed
}

func (idx *Index) IsStale(ctx context.Context, filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return true
	}

	idx.cacheMu.RLock()
	defer idx.cacheMu.RUnlock()

	maxParsedAt, exists := idx.cache[filePath]
	if !exists {
		return true
	}

	return fileInfo.ModTime().Unix() > maxParsedAt
}

func (idx *Index) Index(ctx context.Context, file *parser.File) error {
	err := idx.ensureInitialized(ctx)
	if err != nil {
		return err
	}

	err = idx.Remove(ctx, file.Path)
	if err != nil {
		return err
	}

	if len(file.Chunks) == 0 {
		return nil
	}

	docs := []chromem.Document{}
	for _, chunk := range file.Chunks {
		doc := chromem.Document{
			ID: chunk.ID(),
			Metadata: map[string]string{
				"file":        file.Path,
				"type":        chunk.Type,
				"path":        chunk.Path,
				"summary":     chunk.Summary,
				"startLine":   strconv.Itoa(int(chunk.StartLine)),
				"startColumn": strconv.Itoa(int(chunk.StartColumn)),
				"endLine":     strconv.Itoa(int(chunk.EndLine)),
				"endColumn":   strconv.Itoa(int(chunk.EndColumn)),
				"parsedAt":    strconv.FormatInt(chunk.ParsedAt, 10),
			},
			Content: chunk.Source,
		}

		docs = append(docs, doc)
	}

	err = idx.collection.AddDocuments(ctx, docs, runtime.NumCPU())
	if err != nil {
		return fmt.Errorf("failed to add documents to vector db: %w", err)
	}

	idx.cacheMu.Lock()
	defer idx.cacheMu.Unlock()

	if len(file.Chunks) > 0 {
		idx.cache[file.Path] = file.Chunks[0].ParsedAt
	}

	return nil
}

func (idx *Index) Remove(ctx context.Context, filePath string) error {
	err := idx.ensureInitialized(ctx)
	if err != nil {
		return err
	}

	where := map[string]string{"file": filePath}
	err = idx.collection.Delete(ctx, where, nil)
	if err != nil {
		return fmt.Errorf("failed to remove documents from vector db: %w", err)
	}

	idx.cacheMu.Lock()
	defer idx.cacheMu.Unlock()

	delete(idx.cache, filePath)

	return nil
}

func (idx *Index) Search(ctx context.Context, query string, fileTypes []string) ([]string, error) {
	err := idx.ensureInitialized(ctx)
	if err != nil {
		return nil, err
	}

	if len(fileTypes) == 0 {
		fileTypes = []string{"src", "docs"}
	}

	// chromem-go doesn't support OR filtering, for now fetch more & filter manually
	nResults := min(len(fileTypes)*maxResults, idx.collection.Count())
	results, err := idx.collection.Query(ctx, query, nResults, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to perform similarity search: %w", err)
	}

	allowedTypes := make(map[string]bool)
	for _, ft := range fileTypes {
		allowedTypes[ft] = true
	}

	return idx.formatSearchResults(ctx, results, minSimilarity, maxResults, "", allowedTypes), nil
}

func (idx *Index) FindSimilarChunks(ctx context.Context, chunkID string) ([]string, error) {
	err := idx.ensureInitialized(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := idx.collection.GetByID(ctx, chunkID)
	if err != nil {
		return nil, fmt.Errorf("chunk not found: %s", chunkID)
	}

	results, err := idx.collection.QueryEmbedding(ctx, doc.Embedding, 10, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to perform similarity search: %w", err)
	}

	return idx.formatSearchResults(ctx, results, 2*minSimilarity, 10, chunkID, nil), nil
}

func (idx *Index) formatSearchResults(
	ctx context.Context,
	results []chromem.Result,
	minSimilarity float32,
	maxCount int,
	skipID string,
	typeFilter map[string]bool,
) []string {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	paths := []string{}
	for _, result := range results {
		if result.ID == skipID {
			continue
		}

		if result.Similarity < minSimilarity || len(paths) >= maxCount {
			break
		}

		chunk, err := idx.GetChunk(ctx, result.ID)
		if err != nil {
			continue
		}

		if typeFilter != nil && !typeFilter[chunk.Type] {
			continue
		}

		var lines string
		if chunk.StartLine == chunk.EndLine {
			lines = fmt.Sprintf("line %d", chunk.StartLine)
		} else {
			lines = fmt.Sprintf("lines %d-%d", chunk.StartLine, chunk.EndLine)
		}

		paths = append(
			paths,
			fmt.Sprintf("%s | %s [%s]", result.ID, chunk.Summary, lines),
		)
	}

	return paths
}

func (idx *Index) GetChunk(ctx context.Context, id string) (*parser.Chunk, error) {
	err := idx.ensureInitialized(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := idx.collection.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("chunk not found: %s", id)
	}

	startLine, _ := strconv.Atoi(doc.Metadata["startLine"])
	startColumn, _ := strconv.Atoi(doc.Metadata["startColumn"])
	endLine, _ := strconv.Atoi(doc.Metadata["endLine"])
	endColumn, _ := strconv.Atoi(doc.Metadata["endColumn"])
	parsedAt, _ := strconv.ParseInt(doc.Metadata["parsedAt"], 10, 64)

	return &parser.Chunk{
		File:        doc.Metadata["file"],
		Type:        doc.Metadata["type"],
		Path:        doc.Metadata["path"],
		Summary:     doc.Metadata["summary"],
		Source:      doc.Content,
		StartLine:   uint(startLine),
		StartColumn: uint(startColumn),
		EndLine:     uint(endLine),
		EndColumn:   uint(endColumn),
		ParsedAt:    parsedAt,
	}, nil
}

func (idx *Index) CleanupDeletedFiles(ctx context.Context) {
	err := idx.ensureInitialized(ctx)
	if err != nil {
		return
	}

	idx.cacheMu.RLock()
	defer idx.cacheMu.RUnlock()

	for filePath := range idx.cache {
		go func(path string) {
			_, err := os.Stat(filePath)
			if os.IsNotExist(err) {
				idx.Remove(ctx, filePath)
			}
		}(filePath)
	}
}
