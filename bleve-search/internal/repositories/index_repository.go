package repositories

import (
	queryfunc "bleve-proj/internal/blevefunc"
	"bleve-proj/internal/utils"
	"bleve-proj/pkg/log"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/blevesearch/bleve/v2"
)

type IndexRepository interface {
	LoadIndexes() error
	CreateIndex(partition string, indexMappingJSON map[string]interface{}) (bleve.Index, error)
	OpenIndex(partition string) (bleve.Index, error)
	DeleteIndex(partition string) error
	IndexDocument(partition string, docId string, document interface{}) error
	Search(partition string, query string) (*bleve.SearchResult, error)
	GetDocumentCount(partition string) (uint64, error)
	BulkLoadDocuments(partition string, scanner *bufio.Scanner) error
	CheckIndexContents(partition string) (map[string]interface{}, error)
	PrintTermDictionary(partition, field string) (map[string]uint64, error)
	DumpIndexContentsWithLimit(partition string, page, size int) ([]interface{}, error)
	ListIndexFields(partition string) ([]string, error)
	PrintIndexMapping(partition string) (map[string]interface{}, error)
	GetPartitions() []string
	Close()
	Indexes() map[string]bleve.Index
}

type indexRepository struct {
	logger   log.ILogger
	indexes  map[string]bleve.Index
	mutex    *sync.RWMutex
	indexDir string
}

var (
	IndexRepo        IndexRepository
	ErrIndexNotFound = errors.New("index not found")
	ErrInvalidPage   = errors.New("invalid page")
	ErrInvalidSize   = errors.New("invalid size")
)

func NewIndexRepository(indexDir string) IndexRepository {
	logger := log.NewLogger()
	repo := &indexRepository{
		logger:   logger,
		indexes:  make(map[string]bleve.Index),
		indexDir: indexDir,
		mutex:    &sync.RWMutex{},
	}
	if err := repo.LoadIndexes(); err != nil {
		logger.Fatalf("Failed to load indexes: %v", err)
	}
	return repo
}

func (ir *indexRepository) CreateIndex(partition string, indexMappingJSON map[string]interface{}) (bleve.Index, error) {
	ir.mutex.Lock()
	defer ir.mutex.Unlock()

	// Convert map to JSON byte array
	jsonData, err := json.Marshal(indexMappingJSON)
	if err != nil {
		ir.logger.Errorw("CreateIndex", "failed to marshal index mapping JSON", "error", err)
		return nil, fmt.Errorf("failed to marshal index mapping JSON")
	}

	// Initialize Bleve index mapping
	indexMapping := bleve.NewIndexMapping()
	if err := indexMapping.UnmarshalJSON(jsonData); err != nil {
		ir.logger.Errorw("CreateIndex", "failed to unmarshal index mapping JSON", "error", err)
		return nil, fmt.Errorf("failed to unmarshal index mapping JSON")
	}

	if _, exists := ir.indexes[partition]; exists {
		ir.logger.Errorw("CreateIndex", "index for partition already exists", partition)
		return nil, fmt.Errorf("index for partition %s already exists", partition)
	}

	indexPath := filepath.Join(ir.indexDir, partition)
	newIndex, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		ir.logger.Errorw("CreateIndex", "failed to create new Bleve index", "error", err)
		return nil, fmt.Errorf("failed to create new Bleve index: %s", partition)
	}

	ir.indexes[partition] = newIndex
	return newIndex, nil
}
func (ir *indexRepository) OpenIndex(partition string) (bleve.Index, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
		return nil, ErrIndexNotFound
	}

	return index, nil
}

func (ir *indexRepository) DeleteIndex(partition string) error {
	ir.mutex.Lock()
	defer ir.mutex.Unlock()

	index, exists := ir.indexes[partition]
	if !exists {
		return ErrIndexNotFound
	}

	if err := index.Close(); err != nil {
		ir.logger.Errorw("DeleteIndex", "failed to close index", err)
		return fmt.Errorf("failed to close index")
	}

	indexPath := filepath.Join(ir.indexDir, partition)
	if err := os.RemoveAll(indexPath); err != nil {
		ir.logger.Errorw("DeleteIndex", "failed to delete index files", err)
		return fmt.Errorf("failed to delete index files")
	}

	delete(ir.indexes, partition)
	return nil
}

func (ir *indexRepository) IndexDocument(partition, docId string, document interface{}) error {
	ir.mutex.Lock()
	defer ir.mutex.Unlock()

	index, exists := ir.indexes[partition]
	if !exists {
		return ErrIndexNotFound
	}
	if err := index.Index(docId, document); err != nil {
		ir.logger.Errorw("IndexDocument", "failed to index document", document, "err", err)
		return fmt.Errorf("failed to index document")
	}
	return nil
}

func (ir *indexRepository) Search(partition string, queryString string) (*bleve.SearchResult, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
		ir.logger.Errorw("Search", "index not found", partition)
		return nil, ErrIndexNotFound
	}

	// Handle empty queryString
	if queryString == "" {
		return &bleve.SearchResult{
			Status: &bleve.SearchStatus{
				Total:      0,
				Failed:     0,
				Successful: 0,
				Errors:     nil, // Optionally handle any errors encountered
			},
			Total:    0,
			MaxScore: 0,
			Took:     0,
			Facets:   nil, // Optional: Handle facets if needed
		}, nil
	}

	var queries []map[string]interface{}
	if err := json.Unmarshal([]byte(queryString), &queries); err != nil {
		ir.logger.Errorw("Search", "failed to unmarshal query string", err)
		return nil, fmt.Errorf("failed to unmarshal query string")
	}
	searchRequest := bleve.NewSearchRequest(nil)
	// Extract and prepare search components
	if err := queryfunc.ExtractQueries(searchRequest, queries); err != nil {
		ir.logger.Errorw("Search", "failed to extract queries", err)
		return nil, fmt.Errorf("failed to extract queries")
	}

	if err := queryfunc.ExtractFacets(searchRequest, queries); err != nil {
		ir.logger.Errorw("Search", "failed to extract facets", err)
		return nil, fmt.Errorf("failed to extract facets")
	}

	if err := queryfunc.ExtractSortOrders(searchRequest, queries); err != nil {
		ir.logger.Errorw("Search", "failed to extract sort orders", err)
		return nil, fmt.Errorf("failed to extract sort orders")
	}

	queryfunc.ExtractFields(searchRequest, queries)
	queryfunc.ExtractPagination(searchRequest, queries)
	queryfunc.ExtractMiscellaneous(searchRequest, queries)

	result, err := index.Search(searchRequest)
	if err != nil {
		ir.logger.Errorw("Search", "failed to perform search on index", err)
		return nil, fmt.Errorf("failed to perform search on index")
	}

	return result, nil
}

func (ir *indexRepository) GetDocumentCount(partition string) (uint64, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()
	index, exists := ir.indexes[partition]
	if !exists {
		ir.logger.Errorw("GetDocumentCount", "index not found", partition)
		return 0, ErrIndexNotFound
	}

	return index.DocCount()
}

func (ir *indexRepository) LoadIndexes() error {
	if _, err := os.Stat(ir.indexDir); os.IsNotExist(err) {
		return nil
	}

	indexDirEntries, err := os.ReadDir(ir.indexDir)
	if err != nil {
		ir.logger.Errorw("LoadIndexes", "error reading index directory", err)
		return fmt.Errorf("error reading index directory")
	}
	for _, entry := range indexDirEntries {
		if entry.IsDir() {
			indexName := entry.Name()
			indexPath := filepath.Join(ir.indexDir, indexName)
			index, err := bleve.Open(indexPath)
			if err != nil {
				continue
			}
			ir.mutex.Lock()
			ir.indexes[indexName] = index
			ir.mutex.Unlock()
		}
	}

	return nil
}

func (ir *indexRepository) Indexes() map[string]bleve.Index {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()
	return ir.indexes
}

func (ir *indexRepository) BulkLoadDocuments(partition string, scanner *bufio.Scanner) error {
	ir.mutex.Lock()
	defer ir.mutex.Unlock()

	index, exists := ir.indexes[partition]
	if !exists {
		ir.logger.Errorw("BulkLoadDocuments", "index not found", partition)
		return ErrIndexNotFound
	}

	batch := index.NewBatch()
	for scanner.Scan() {
		var doc map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &doc); err != nil {
			ir.logger.Errorw("BulkLoadDocuments", "failed to unmarshal bytes", string(scanner.Bytes()), "err", err)
			continue
		}

		docID := doc["id"]
		if docID == nil {
			docID = utils.GenerateUniqueID()
			doc["id"] = docID.(string)
		}

		batch.Index(docID.(string), doc)
	}

	if err := scanner.Err(); err != nil {
		ir.logger.Errorw("BulkLoadDocuments", "got error from scanner", err)
		return fmt.Errorf("got error from scanner")
	}
	if err := index.Batch(batch); err != nil {
		ir.logger.Errorw("BulkLoadDocuments", "failed to batch data", err)
		return fmt.Errorf("failed to batch data")
	}
	return nil
}

func (ir *indexRepository) CheckIndexContents(partition string) (map[string]interface{}, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
		ir.logger.Errorw("CheckIndexContents", "index not found", partition)
		return nil, ErrIndexNotFound
	}
	fields, err := index.Fields()
	if err != nil {
		ir.logger.Errorw("CheckIndexContents", "can not extract fields from index", partition, "err", err)
		return nil, fmt.Errorf("can not extract fields from index")
	}

	result := make(map[string]interface{})
	for _, field := range fields {
		dict, err := index.FieldDict(field)
		if err != nil {
			ir.logger.Errorw("CheckIndexContents", "can not extract fields dict from index", partition, "field", field, "err", err)
			return nil, fmt.Errorf("can not extract fields dict from index")
		}
		defer dict.Close()

		termCount := 0
		for term, err := dict.Next(); term != nil && err == nil; term, err = dict.Next() {
			termCount++
		}

		result[field] = map[string]interface{}{
			"termCount": termCount,
		}
	}

	return result, nil
}

func (ir *indexRepository) PrintTermDictionary(partition, field string) (map[string]uint64, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
		ir.logger.Errorw("PrintTermDictionary", "index not found", partition)
		return nil, ErrIndexNotFound
	}

	termDict := make(map[string]uint64)
	dict, err := index.FieldDict(field)
	if err != nil {
		ir.logger.Errorw("PrintTermDictionary", "can not extract fields dict from index", partition, "field", field, "err", err)
		return nil, fmt.Errorf("can not extract fields dict from index")
	}
	defer dict.Close()

	for {
		entry, err := dict.Next()
		if err != nil {
			ir.logger.Errorw("PrintTermDictionary", "can not extract the next dict", err)
			return nil, fmt.Errorf("can not extract the next dict")
		}
		if entry == nil {
			break
		}
		termDict[entry.Term] = entry.Count
	}

	return termDict, nil
}

func (ir *indexRepository) DumpIndexContentsWithLimit(partition string, page, size int) ([]interface{}, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()
	// Calculate offset and limit
	if page <= 0 {
		ir.logger.Errorw("DumpIndexContentsWithLimit", "invalid page", page)
		return nil, ErrInvalidPage
	}
	if size <= 0 {
		ir.logger.Errorw("DumpIndexContentsWithLimit", "invalid size", size)
		return nil, ErrInvalidSize
	}
	offset := (page - 1) * size
	limit := size
	index, exists := ir.indexes[partition]
	if !exists {
		ir.logger.Errorw("DumpIndexContentsWithLimit", "index not found", partition)
		return nil, ErrIndexNotFound
	}

	// Query to fetch documents with offset and limit
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchRequest.From = offset
	searchRequest.Size = limit
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		ir.logger.Errorw("DumpIndexContentsWithLimit", "failed to search data to dump", err)
		return nil, fmt.Errorf("failed to search data to dump")
	}

	// Collect documents from search result
	var documents []interface{}
	for _, hit := range searchResult.Hits {
		documents = append(documents, hit.Fields)
	}

	return documents, nil
}

func (ir *indexRepository) ListIndexFields(partition string) ([]string, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
		ir.logger.Errorw("ListIndexFields", "index not found", partition)
		return nil, ErrIndexNotFound
	}
	fields, err := index.Fields()
	if err != nil {
		ir.logger.Errorw("ListIndexFields", "failed to extract index fields", partition, "err", err)
		return nil, fmt.Errorf("failed to extract index fields")
	}
	return fields, nil
}

func (ir *indexRepository) PrintIndexMapping(partition string) (map[string]interface{}, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
		ir.logger.Errorw("PrintIndexMapping", "index not found", partition)
		return nil, ErrIndexNotFound
	}

	mapping := index.Mapping()
	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		ir.logger.Errorw("PrintIndexMapping", "failed to marshal index mapping", partition, "err", err)
		return nil, fmt.Errorf("failed to marshal index mapping")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(mappingJSON, &result); err != nil {
		ir.logger.Errorw("PrintIndexMapping", "failed to unmarshal index mapping", partition, "err", err)
		return nil, fmt.Errorf("failed to unmarshal index mapping")
	}

	return result, nil
}

func (ir *indexRepository) GetPartitions() []string {
	result := []string{}
	for key := range ir.indexes {
		result = append(result, key)
	}
	return result
}

func (ir *indexRepository) Close() {
	indexes := []string{}
	for i := range ir.indexes {
		ir.indexes[i].Close()
		indexes = append(indexes, i)
	}
	ir.logger.Infow("Close", "close all index", indexes)

}
