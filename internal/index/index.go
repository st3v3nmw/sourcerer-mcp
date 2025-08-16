package index

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/philippgille/chromem-go"
	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

const (
	minSimilarity = 0.3
	maxResults    = 15
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
		indexedAt := strconv.FormatInt(file.ParsedAt.UnixMicro(), 10)
		doc := chromem.Document{
			ID: file.Path + "::" + chunk.Path,
			Metadata: map[string]string{
				"file":      file.Path,
				"path":      chunk.Path,
				"summary":   chunk.Summary,
				"startByte": strconv.Itoa(int(chunk.StartByte)),
				"endByte":   strconv.Itoa(int(chunk.EndByte)),
				"indexedAt": indexedAt,
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

func (idx *Index) isDocumentStale(doc *chromem.Document) bool {
	filePath := doc.Metadata["file"]
	fullPath := path.Join(idx.workspaceRoot, filePath)

	info, err := os.Stat(fullPath)
	if err != nil {
		return true
	}

	indexedAt := doc.Metadata["indexedAt"]
	indexedAtUnixMicro, _ := strconv.ParseInt(indexedAt, 10, 64)
	return info.ModTime().UnixMicro() > indexedAtUnixMicro
}

func (idx *Index) IsChunkStale(ctx context.Context, id string, fileModTime time.Time) (bool, error) {
	chunk, err := idx.collection.GetByID(ctx, id)
	if err != nil {
		return true, nil
	}

	return idx.isDocumentStale(&chunk), nil
}

func (idx *Index) GetChunk(ctx context.Context, id string) (*parser.Chunk, error) {
	doc, err := idx.collection.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("chunk not found: %s", id)
	}

	startByte, _ := strconv.Atoi(doc.Metadata["startByte"])
	endByte, _ := strconv.Atoi(doc.Metadata["endByte"])

	return &parser.Chunk{
		Path:      doc.Metadata["path"],
		Summary:   doc.Metadata["summary"],
		Source:    doc.Content,
		StartByte: uint(startByte),
		EndByte:   uint(endByte),
	}, nil
}

func (idx *Index) GetAllChunkIDs(ctx context.Context) map[string][]string {
	ids := idx.collection.ListIDs(ctx)

	fileChunks := make(map[string][]struct {
		id        string
		startByte int
	})

	for _, id := range ids {
		doc, err := idx.collection.GetByID(ctx, id)
		if err != nil {
			continue
		}

		filePath := doc.Metadata["file"]
		startByte, _ := strconv.Atoi(doc.Metadata["startByte"])

		fileChunks[filePath] = append(fileChunks[filePath], struct {
			id        string
			startByte int
		}{
			id:        id,
			startByte: startByte,
		})
	}

	result := make(map[string][]string)
	for filePath, chunks := range fileChunks {
		sort.Slice(chunks, func(i, j int) bool {
			return chunks[i].startByte < chunks[j].startByte
		})

		chunkIDs := make([]string, len(chunks))
		for i, chunk := range chunks {
			chunkIDs[i] = chunk.id
		}
		result[filePath] = chunkIDs
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

		doc, err := idx.collection.GetByID(ctx, result.ID)
		if err != nil {
			staleChunkIDs = append(staleChunkIDs, result.ID)
			continue
		}

		if idx.isDocumentStale(&doc) {
			staleChunkIDs = append(staleChunkIDs, result.ID)
			continue
		}

		paths = append(paths, result.ID+" | "+doc.Metadata["summary"])
	}

	if len(staleChunkIDs) > 0 {
		go idx.collection.Delete(ctx, nil, nil, staleChunkIDs...)
	}

	return paths, nil
}
