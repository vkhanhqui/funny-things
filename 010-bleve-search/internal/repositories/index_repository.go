package repositories

import (
	"bleve-proj/internal/utils"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	Indexes() map[string]bleve.Index
}

type indexRepository struct {
	indexes  map[string]bleve.Index
	mutex    *sync.RWMutex
	indexDir string
}

var (
	IndexRepo        IndexRepository
	ErrIndexNotFound = errors.New("index not found")
	ErrIndexExisted  = errors.New("index already exists")
	ErrInvalidPage   = errors.New("invalid page")
	ErrInvalidSize   = errors.New("invalid size")
)

func NewIndexRepository(indexDir string) IndexRepository {
	repo := &indexRepository{
		indexes:  make(map[string]bleve.Index),
		indexDir: indexDir,
		mutex:    &sync.RWMutex{},
	}
	if err := repo.LoadIndexes(); err != nil {
		log.Fatalf("Failed to load indexes: %v", err)
	}
	return repo
}

func (ir *indexRepository) CreateIndex(partition string, indexMappingJSON map[string]interface{}) (bleve.Index, error) {
	ir.mutex.Lock()
	defer ir.mutex.Unlock()
	// Convert map to JSON byte array
	jsonData, err := json.Marshal(indexMappingJSON)
	if err != nil {
		return nil, err
	}
	// Initialize Bleve index mapping
	indexMapping := bleve.NewIndexMapping()
	if err := indexMapping.UnmarshalJSON(jsonData); err != nil {
		return nil, err
	}
	if _, exists := ir.indexes[partition]; exists {
		return nil, ErrIndexExisted
	}
	indexPath := filepath.Join(ir.indexDir, partition)
	newIndex, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		return nil, err
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
		return fmt.Errorf("failed to close index: %w", err)
	}

	indexPath := filepath.Join(ir.indexDir, partition)
	if err := os.RemoveAll(indexPath); err != nil {
		return fmt.Errorf("failed to delete index files: %w", err)
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
	return index.Index(docId, document)
}

// Search performs a search on the specified partition with the given query
func (ir *indexRepository) Search(partition string, queryString string) (*bleve.SearchResult, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
		return nil, ErrIndexNotFound
	}

	// Handle empty queryString
	if queryString == "" {
		// Create an empty search result indicating no hits
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
		return nil, err
	}

	searchRequest := bleve.NewSearchRequest(nil)

	// Extract and prepare search components
	err := extractQueries(searchRequest, queries)
	if err != nil {
		return nil, err
	}

	err = extractFields(searchRequest, queries)
	if err != nil {
		return nil, err
	}

	err = extractPagination(searchRequest, queries)
	if err != nil {
		return nil, err
	}

	err = extractMiscellaneous(searchRequest, queries)
	if err != nil {
		return nil, err
	}

	err = extractFacets(searchRequest, queries)
	if err != nil {
		return nil, err
	}

	err = extractSortOrders(searchRequest, queries)
	if err != nil {
		return nil, err
	}

	result, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ir *indexRepository) GetDocumentCount(partition string) (uint64, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
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
		return fmt.Errorf("error reading index directory: %w", err)
	}
	for _, entry := range indexDirEntries {
		if entry.IsDir() {
			indexName := entry.Name()
			indexPath := filepath.Join(ir.indexDir, indexName)
			index, err := bleve.Open(indexPath)
			if err != nil {
				log.Printf("Failed to open index '%s': %v", indexName, err)
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
		return ErrIndexNotFound
	}

	batch := index.NewBatch()
	for scanner.Scan() {
		var doc map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &doc); err != nil {
			fmt.Printf("scanner.Bytes(): %v\n", string(scanner.Bytes()))
			return err
		}

		docID := doc["id"]
		if docID == nil {
			docID = utils.GenerateUniqueID()
			doc["id"] = docID.(string)
		}

		batch.Index(docID.(string), doc)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return index.Batch(batch)
}

func (ir *indexRepository) CheckIndexContents(partition string) (map[string]interface{}, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
		return nil, ErrIndexNotFound
	}
	fields, err := index.Fields()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, field := range fields {
		dict, err := index.FieldDict(field)
		if err != nil {
			return nil, err
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
		return nil, ErrIndexNotFound
	}

	termDict := make(map[string]uint64)
	dict, err := index.FieldDict(field)
	if err != nil {
		return nil, err
	}
	defer dict.Close()

	for {
		entry, err := dict.Next()
		if err != nil {
			return nil, err
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
		return nil, ErrInvalidPage
	}
	if size <= 0 {
		return nil, ErrInvalidSize
	}
	offset := (page - 1) * size
	limit := size
	index, exists := ir.indexes[partition]
	if !exists {
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
		return nil, err
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
		return nil, ErrIndexNotFound
	}

	return index.Fields()
}

func (ir *indexRepository) PrintIndexMapping(partition string) (map[string]interface{}, error) {
	ir.mutex.RLock()
	defer ir.mutex.RUnlock()

	index, exists := ir.indexes[partition]
	if !exists {
		return nil, ErrIndexNotFound
	}

	mapping := index.Mapping()
	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(mappingJSON, &result); err != nil {
		return nil, err
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
