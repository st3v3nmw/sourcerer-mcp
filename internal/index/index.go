package index

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"

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
}

func New(workspaceRoot string) (*Index, error) {
	db, err := chromem.NewPersistentDB(".sourcerer/db", false)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector db: %w", err)
	}

	collection, err := db.GetOrCreateCollection("code", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector db collection: %w", err)
	}

	return &Index{
		workspaceRoot: workspaceRoot,
		collection:    collection,
	}, nil
}

func (idx *Index) Upsert(ctx context.Context, file *parser.File) error {
	err := idx.Remove(ctx, file.Path)
	if err != nil {
		return err
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

	return nil
}

func (idx *Index) Remove(ctx context.Context, filePath string) error {
	where := map[string]string{"file": filePath}
	err := idx.collection.Delete(ctx, where, nil)
	if err != nil {
		return fmt.Errorf("failed to remove documents from vector db: %w", err)
	}

	return nil
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

func (idx *Index) isChunkStale(chunk *parser.Chunk) bool {
	fullPath := path.Join(idx.workspaceRoot, chunk.File)

	info, err := os.Stat(fullPath)
	if err != nil {
		return true
	}

	return info.ModTime().UnixMicro() > chunk.ParsedAt
}

func (idx *Index) IsChunkStale(ctx context.Context, id string) (bool, error) {
	chunk, err := idx.GetChunk(ctx, id)
	if err != nil {
		return true, nil
	}

	return idx.isChunkStale(chunk), nil
}

func (idx *Index) GetAllChunkIDs(ctx context.Context) map[string][]string {
	ids := idx.collection.ListIDs(ctx)

	fileChunks := make(map[string][]struct {
		id        string
		startLine uint
	})

	staleChunkIDs := []string{}
	for _, id := range ids {
		chunk, err := idx.GetChunk(ctx, id)
		if err != nil {
			continue
		}

		if idx.isChunkStale(chunk) {
			staleChunkIDs = append(staleChunkIDs, chunk.ID())
			continue
		}

		fileChunks[chunk.File] = append(fileChunks[chunk.File], struct {
			id        string
			startLine uint
		}{
			id:        id,
			startLine: chunk.StartLine,
		})
	}

	result := make(map[string][]string)
	for filePath, chunks := range fileChunks {
		sort.Slice(chunks, func(i, j int) bool {
			return chunks[i].startLine < chunks[j].startLine
		})

		chunkIDs := make([]string, len(chunks))
		for i, chunk := range chunks {
			chunkIDs[i] = chunk.id
		}
		result[filePath] = chunkIDs
	}

	if len(staleChunkIDs) > 0 {
		go idx.collection.Delete(ctx, nil, nil, staleChunkIDs...)
	}

	return result
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
	staleChunkIDs := []string{}
	for i, result := range results {
		if result.Similarity < minSimilarity || i >= maxResults {
			break
		}

		chunk, err := idx.GetChunk(ctx, result.ID)
		if err != nil {
			continue
		}

		if idx.isChunkStale(chunk) {
			staleChunkIDs = append(staleChunkIDs, result.ID)
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

	if len(staleChunkIDs) > 0 {
		go idx.collection.Delete(ctx, nil, nil, staleChunkIDs...)
	}

	return paths, nil
}
