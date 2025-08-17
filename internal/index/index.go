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

type Index struct {
	workspaceRoot string
	collection    *chromem.Collection

	fileMetadata map[string]int64
	metadataMu   sync.RWMutex
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
		fileMetadata:  map[string]int64{},
	}

	idx.loadFileMetadata(ctx)

	return idx, nil
}

func (idx *Index) loadFileMetadata(ctx context.Context) {
	seen := map[string]bool{}
	chunkIDs := idx.collection.ListIDs(ctx)
	for _, chunkID := range chunkIDs {
		parts := strings.SplitN(chunkID, "::", 2)
		if len(parts) != 2 {
			continue
		}

		filePath := parts[0]
		_, exists := idx.fileMetadata[filePath]
		if !exists && !seen[filePath] {
			chunk, err := idx.GetChunk(ctx, chunkID)
			if err != nil {
				continue
			}

			seen[filePath] = true
			if !idx.IsStale(filePath) {
				where := map[string]string{"file": filePath}
				idx.collection.Delete(ctx, where, nil)
				continue
			}

			idx.fileMetadata[filePath] = chunk.ParsedAt
		}
	}

}

func (idx *Index) IsStale(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	idx.metadataMu.RLock()
	defer idx.metadataMu.RUnlock()

	parsedAt, exists := idx.fileMetadata[filePath]
	if !exists {
		return true
	}

	return fileInfo.ModTime().UnixMicro() > parsedAt
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

	idx.metadataMu.Lock()
	idx.fileMetadata[file.Path] = file.Chunks[0].ParsedAt
	idx.metadataMu.Unlock()

	return nil
}

func (idx *Index) Remove(ctx context.Context, filePath string) error {
	where := map[string]string{"file": filePath}
	err := idx.collection.Delete(ctx, where, nil)
	if err != nil {
		return fmt.Errorf("failed to remove documents from vector db: %w", err)
	}

	idx.metadataMu.Lock()
	defer idx.metadataMu.Unlock()

	delete(idx.fileMetadata, filePath)

	return nil
}

func (idx *Index) Search(ctx context.Context, query string) ([]string, error) {
	results, err := idx.collection.Query(ctx, query, 2*maxResults, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to perform similarity search: %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	paths := []string{}
	for i, result := range results {
		if result.Similarity < minSimilarity || i >= maxResults {
			break
		}

		chunk, err := idx.GetChunk(ctx, result.ID)
		if err != nil {
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

	return paths, nil
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
