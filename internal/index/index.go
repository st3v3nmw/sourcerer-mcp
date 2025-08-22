package index

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/philippgille/chromem-go"
	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

const (
	minSimilarity = 0.3
	maxResults    = 30
)

type ChunkMetadata struct {
	Type     string // chunk type (src, docs, etc)
	Path     string // hierarchical path: Class::method
	ParsedAt int64  // when chunk was parsed
}

type Index struct {
	workspaceRoot string
	collection    *chromem.Collection

	cache   map[string][]*ChunkMetadata
	cacheMu sync.RWMutex
}

func New(ctx context.Context, workspaceRoot string) (*Index, error) {
	db, err := chromem.NewPersistentDB(".sourcerer/db", false)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector db: %w", err)
	}

	collection, err := db.GetOrCreateCollection("code-chunks", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector db collection: %w", err)
	}

	idx := &Index{
		workspaceRoot: workspaceRoot,
		collection:    collection,
		cache:         map[string][]*ChunkMetadata{},
	}

	idx.loadCache(ctx)

	return idx, nil
}

func (idx *Index) loadCache(ctx context.Context) {
	idx.cacheMu.Lock()
	defer idx.cacheMu.Unlock()

	fileChunks := map[string][]*ChunkMetadata{}
	chunkIDs := idx.collection.ListIDs(ctx)
	for _, chunkID := range chunkIDs {
		parts := strings.SplitN(chunkID, "::", 2)
		if len(parts) != 2 {
			continue
		}

		filePath := parts[0]

		// Check if file still exists
		_, err := os.Stat(filePath)
		if err != nil {
			where := map[string]string{"file": filePath}
			idx.collection.Delete(ctx, where, nil)
			continue
		}

		chunk, err := idx.GetChunk(ctx, chunkID)
		if err != nil {
			continue
		}

		fileChunks[filePath] = append(fileChunks[filePath], &ChunkMetadata{
			Type:     chunk.Type,
			Path:     chunk.Path,
			ParsedAt: chunk.ParsedAt,
		})
	}

	idx.cache = fileChunks
}

func (idx *Index) IsStale(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return true
	}

	idx.cacheMu.RLock()
	defer idx.cacheMu.RUnlock()

	chunks, exists := idx.cache[filePath]
	if !exists || len(chunks) == 0 {
		return true
	}

	// Find the maximum parsedAt timestamp among all chunks
	var maxParsedAt int64
	for _, chunk := range chunks {
		if chunk.ParsedAt > maxParsedAt {
			maxParsedAt = chunk.ParsedAt
		}
	}

	return fileInfo.ModTime().Unix() > maxParsedAt
}

func (idx *Index) Index(ctx context.Context, file *parser.File) error {
	err := idx.Remove(ctx, file.Path)
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

	chunkMetadata := make([]*ChunkMetadata, 0, len(file.Chunks))
	for _, chunk := range file.Chunks {
		chunkMetadata = append(chunkMetadata, &ChunkMetadata{
			Type:     chunk.Type,
			Path:     chunk.Path,
			ParsedAt: chunk.ParsedAt,
		})
	}
	idx.cache[file.Path] = chunkMetadata

	return nil
}

func (idx *Index) Remove(ctx context.Context, filePath string) error {
	where := map[string]string{"file": filePath}
	err := idx.collection.Delete(ctx, where, nil)
	if err != nil {
		return fmt.Errorf("failed to remove documents from vector db: %w", err)
	}

	idx.cacheMu.Lock()
	defer idx.cacheMu.Unlock()

	delete(idx.cache, filePath)

	return nil
}

func (idx *Index) Search(ctx context.Context, query string, fileTypes []string) ([]string, error) {
	if len(fileTypes) == 0 {
		fileTypes = []string{"src", "docs"}
	}

	// chromem-go doesn't support OR filtering, for now fetch more & filter manually
	results, err := idx.collection.Query(ctx, query, len(fileTypes)*maxResults, nil, nil)
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
