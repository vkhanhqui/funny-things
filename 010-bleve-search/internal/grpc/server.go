package indexgrpc

import (
	queryfunc "bleve-proj/internal/blevefunc"
	"bleve-proj/internal/utils"
	"bleve-proj/pkg/log"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	proto "bleve-proj/proto"

	"github.com/blevesearch/bleve/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type server struct {
	proto.UnimplementedIndexServiceServer
	indexes  map[string]bleve.Index
	mutex    *sync.RWMutex
	indexDir string
	logger   log.ILogger
}

func NewServer(indexDir string) *server {
	server := &server{
		indexes:  make(map[string]bleve.Index),
		mutex:    &sync.RWMutex{},
		indexDir: indexDir,
		logger:   log.NewLogger(),
	}
	if err := server.LoadIndexes(); err != nil {
		server.logger.Errorw("NewServer", "could not load indexes", err)
	}
	return server
}

func (s *server) CreateIndex(ctx context.Context, req *proto.CreateIndexRequest) (*proto.CreateIndexResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Convert structpb.Struct to map[string]interface{}
	indexMappingMap := req.IndexMapping.AsMap()

	// Marshal the map to JSON
	jsonData, err := json.Marshal(indexMappingMap)
	if err != nil {
		s.logger.Errorw("Create Index", "Failed to marshal index mapping:", err)
		return nil, fmt.Errorf("failed to marshal index mapping: %w", err)
	}

	indexMapping := bleve.NewIndexMapping()
	if err := indexMapping.UnmarshalJSON(jsonData); err != nil {
		s.logger.Errorw("Create Index", "Failed to unmarshal index mapping JSON:", err)
		return nil, fmt.Errorf("failed to unmarshal index mapping JSON: %w", err)
	}

	if _, exists := s.indexes[req.Partition]; exists {
		s.logger.Errorw("Create Index", "Index for partition already exists:", req.Partition)
		return nil, fmt.Errorf("index for partition %s already exists", req.Partition)
	}

	indexPath := filepath.Join(s.indexDir, req.Partition)
	newIndex, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		s.logger.Errorw("Create Index", "Failed to create new Bleve index:", err)
		return nil, fmt.Errorf("failed to create new Bleve index: %s", req.Partition)
	}

	s.indexes[req.Partition] = newIndex
	return &proto.CreateIndexResponse{Message: "Index created successfully"}, nil
}

func (s *server) LoadIndexes() error {
	if _, err := os.Stat(s.indexDir); os.IsNotExist(err) {
		s.logger.Errorw("LoadIndexes", "indexDir is not exist", s.indexDir, "err", err)
		return nil
	}

	indexDirEntries, err := os.ReadDir(s.indexDir)
	if err != nil {
		s.logger.Errorw("LoadIndexes", "error reading index directory", err)
		return fmt.Errorf("error reading index directory")
	}
	for _, entry := range indexDirEntries {
		if entry.IsDir() {
			indexName := entry.Name()
			indexPath := filepath.Join(s.indexDir, indexName)
			index, err := bleve.Open(indexPath)
			if err != nil {
				continue
			}
			s.mutex.Lock()
			s.indexes[indexName] = index
			s.mutex.Unlock()
		}
	}

	return nil
}

func (s *server) IndexDocument(ctx context.Context, req *proto.IndexDocumentRequest) (*proto.IndexDocumentResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Logging the received request

	index, exists := s.indexes[req.Partition]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Index partition %s not found", req.Partition)
	}
	var docID string
	// Convert the document from Struct to a map
	if req.GetDocId() != "" {
		docID = req.GetDocId()
	} else {
		docID = utils.GenerateUniqueID() // Generate or handle as required
	}
	if err := index.Index(docID, req.Document.AsMap()); err != nil {
		s.logger.Errorw("IndexDocument", "failed to index document", req.Document, "err", err)
		return nil, status.Errorf(codes.Internal, "failed to index document: %v", err)
	}

	return &proto.IndexDocumentResponse{
		Message: "Document indexed successfully",
	}, nil
}

func (s *server) DeleteIndex(ctx context.Context, req *proto.DeleteIndexRequest) (*proto.DeleteIndexResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	index, exists := s.indexes[req.Partition]
	if !exists {
		s.logger.Errorw("Delete Index", "Index not found:", req.Partition)
		return nil, fmt.Errorf("index for partition %s not found", req.Partition)
	}

	if err := index.Close(); err != nil {
		s.logger.Errorw("Delete Index", "Failed to close index:", err)
		return nil, fmt.Errorf("failed to close index: %s", req.Partition)
	}

	indexPath := filepath.Join(s.indexDir, req.Partition)
	if err := os.RemoveAll(indexPath); err != nil {
		s.logger.Errorw("Delete Index", "Failed to delete index files:", err)
		return nil, fmt.Errorf("failed to delete index files: %s", req.Partition)
	}

	delete(s.indexes, req.Partition)
	return &proto.DeleteIndexResponse{Message: "Index deleted successfully"}, nil
}

func (s *server) Search(ctx context.Context, req *proto.SearchRequest) (*proto.SearchResult, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	index, exists := s.indexes[req.Partition]
	if !exists {
		s.logger.Errorw("Search", "Index not found:", req.Partition)
		return nil, fmt.Errorf("index for partition %s not found", req.Partition)
	}
	result, err := s.searchIndex(index, req.Query)
	if err != nil {
		s.logger.Errorw("Search", "Failed to search index:", err)
		return nil, fmt.Errorf("failed to search index: %s", req.Partition)
	}
	protoResult := &proto.SearchResult{
		Status: &proto.SearchStatus{
			Total:      int32(result.Status.Total),
			Failed:     int32(result.Status.Failed),
			Successful: int32(result.Status.Successful),
		},
		Hits:      make([]*proto.SearchHit, len(result.Hits)),
		TotalHits: uint64(result.Total),
		Cost:      uint64(result.Cost),
		MaxScore:  result.MaxScore,
		Took:      int32(result.Took.Milliseconds()), // Convert time.Duration to milliseconds
		Facets:    make(map[string]*proto.FacetResult),
	}

	for i, hit := range result.Hits {
		fields := make(map[string]string)
		for k, v := range hit.Fields {
			// Safely handle different data types in hit.Fields
			switch val := v.(type) {
			case string:
				fields[k] = val
			case float64:
				fields[k] = fmt.Sprintf("%f", val)
			case int, int32, int64:
				fields[k] = fmt.Sprintf("%d", val)
			// Add other types as needed
			default:
				fields[k] = fmt.Sprintf("%v", val) // Default to string representation
			}
		}
		protoResult.Hits[i] = &proto.SearchHit{
			Index:  hit.Index,
			Id:     hit.ID,
			Score:  hit.Score,
			Sort:   hit.Sort,
			Fields: fields,
		}
	}

	for facetName, facetResult := range result.Facets {
		protoFacet := &proto.FacetResult{
			Field:   facetResult.Field,
			Total:   int32(facetResult.Total),
			Missing: int32(facetResult.Missing),
			Other:   int32(facetResult.Other),
			// Terms:   make([]*proto.TermFacet, 0), // You can comment out or remove this line for now
		}

		protoResult.Facets[facetName] = protoFacet
	}

	return protoResult, nil
}

func (s *server) searchIndex(index bleve.Index, queryString string) (*bleve.SearchResult, error) {
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
		return nil, fmt.Errorf("failed to unmarshal query string")
	}

	searchRequest := bleve.NewSearchRequest(nil)

	// Extract and prepare search components
	if err := queryfunc.ExtractQueries(searchRequest, queries); err != nil {
		return nil, fmt.Errorf("failed to extract queries")
	}

	if err := queryfunc.ExtractFacets(searchRequest, queries); err != nil {
		return nil, fmt.Errorf("failed to extract facets")
	}
	if err := queryfunc.ExtractSortOrders(searchRequest, queries); err != nil {
		return nil, fmt.Errorf("failed to extract sort orders")
	}

	queryfunc.ExtractFields(searchRequest, queries)
	queryfunc.ExtractPagination(searchRequest, queries)
	queryfunc.ExtractMiscellaneous(searchRequest, queries)
	return index.Search(searchRequest)
}

func (s *server) GetDocumentCount(ctx context.Context, req *proto.GetDocumentCountRequest) (*proto.GetDocumentCountResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	index, exists := s.indexes[req.Partition]
	if !exists {
		s.logger.Errorw("GetDocumentCount", "index not found", req.Partition)
		return nil, status.Errorf(codes.NotFound, "index for partition %s not found", req.Partition)
	}

	count, err := index.DocCount()
	if err != nil {
		s.logger.Errorw("GetDocumentCount", "failed to get document count", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to get document count: %v", err)
	}

	return &proto.GetDocumentCountResponse{Count: count}, nil
}

func (s *server) BulkLoadDocuments(ctx context.Context, req *proto.BulkLoadDocumentsRequest) (*proto.BulkLoadDocumentsResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	index, exists := s.indexes[req.Partition]
	if !exists {
		s.logger.Errorw("BulkLoadDocuments", "index not found", req.Partition)
		return nil, status.Errorf(codes.NotFound, "index for partition %s not found", req.Partition)
	}

	scanner := bufio.NewScanner(bytes.NewReader(req.Documents))
	batch := index.NewBatch()

	for scanner.Scan() {
		var doc map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &doc); err != nil {
			s.logger.Errorw("BulkLoadDocuments", "failed to unmarshal bytes", string(scanner.Bytes()), "err", err)
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
		s.logger.Errorw("BulkLoadDocuments", "got error from scanner", err)
		return nil, fmt.Errorf("got error from scanner: %w", err) // Wrap the error
	}

	if err := index.Batch(batch); err != nil {
		s.logger.Errorw("BulkLoadDocuments", "failed to batch data", err)
		return nil, fmt.Errorf("failed to batch data: %w", err) // Wrap the error
	}

	return &proto.BulkLoadDocumentsResponse{Message: "Documents bulk loaded successfully"}, nil
}

func (s *server) CheckIndexContents(ctx context.Context, req *proto.CheckIndexContentsRequest) (*proto.CheckIndexContentsResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	index, exists := s.indexes[req.Partition]
	if !exists {
		s.logger.Errorw("CheckIndexContents", "index not found", req.Partition)
		return nil, status.Errorf(codes.NotFound, "index for partition %s not found", req.Partition)
	}

	fields, err := index.Fields()
	if err != nil {
		s.logger.Errorw("CheckIndexContents", "cannot extract fields from index", req.Partition, "err", err)
		return nil, fmt.Errorf("cannot extract fields from index: %w", err)
	}

	result := &proto.CheckIndexContentsResponse{
		Fields: make(map[string]*proto.FieldInfo),
	}
	for _, field := range fields {
		dict, err := index.FieldDict(field)
		if err != nil {
			s.logger.Errorw("CheckIndexContents", "cannot extract field dictionary from index", req.Partition, "field", field, "err", err)
			return nil, fmt.Errorf("cannot extract field dictionary from index: %w", err)
		}
		defer dict.Close()

		termCount := uint64(0)
		for term, err := dict.Next(); term != nil && err == nil; term, err = dict.Next() {
			termCount++
		}

		result.Fields[field] = &proto.FieldInfo{
			TermCount: termCount,
		}
	}

	return result, nil
}

func (s *server) PrintTermDictionary(ctx context.Context, req *proto.PrintTermDictionaryRequest) (*proto.PrintTermDictionaryResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	index, exists := s.indexes[req.Partition]
	if !exists {
		s.logger.Errorw("PrintTermDictionary", "index not found", req.Partition)
		return nil, status.Errorf(codes.NotFound, "index for partition %s not found", req.Partition)
	}

	termDict := make(map[string]uint64)
	dict, err := index.FieldDict(req.Field)
	if err != nil {
		s.logger.Errorw("PrintTermDictionary", "cannot extract field dictionary from index", req.Partition, "field", req.Field, "err", err)
		return nil, fmt.Errorf("cannot extract field dictionary from index: %w", err)
	}
	defer dict.Close()

	for {
		entry, err := dict.Next()
		if err != nil {
			s.logger.Errorw("PrintTermDictionary", "cannot extract the next dictionary entry", err)
			return nil, fmt.Errorf("cannot extract the next dictionary entry: %w", err)
		}
		if entry == nil {
			break
		}
		termDict[entry.Term] = entry.Count
	}

	return &proto.PrintTermDictionaryResponse{TermCounts: termDict}, nil
}

func (s *server) DumpIndexContentsWithLimit(ctx context.Context, req *proto.DumpIndexContentsWithLimitRequest) (*proto.DumpIndexContentsWithLimitResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Input validation
	if req.Page <= 0 {
		s.logger.Errorw("DumpIndexContentsWithLimit", "invalid page", req.Page)
		return nil, status.Errorf(codes.InvalidArgument, "invalid page: %d", req.Page)
	}
	if req.Size <= 0 {
		s.logger.Errorw("DumpIndexContentsWithLimit", "invalid size", req.Size)
		return nil, status.Errorf(codes.InvalidArgument, "invalid size: %d", req.Size)
	}

	index, exists := s.indexes[req.Partition]
	if !exists {
		s.logger.Errorw("DumpIndexContentsWithLimit", "index not found", req.Partition)
		return nil, status.Errorf(codes.NotFound, "index for partition %s not found", req.Partition)
	}

	// Query to fetch documents with offset and limit
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchRequest.From = int((req.Page - 1) * req.Size)
	searchRequest.Size = int(req.Size)
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		s.logger.Errorw("DumpIndexContentsWithLimit", "failed to search data to dump", err)
		return nil, fmt.Errorf("failed to search data to dump: %w", err)
	}

	// Collect documents from search result
	var documents []*structpb.Struct
	for _, hit := range searchResult.Hits {
		doc, err := structpb.NewStruct(hit.Fields)
		if err != nil {
			return nil, fmt.Errorf("failed to convert hit.Fields to structpb.Struct: %w", err)
		}
		documents = append(documents, doc)
	}

	return &proto.DumpIndexContentsWithLimitResponse{Documents: documents}, nil
}

func (s *server) ListIndexFields(ctx context.Context, req *proto.ListIndexFieldsRequest) (*proto.ListIndexFieldsResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	index, exists := s.indexes[req.Partition]
	if !exists {
		s.logger.Errorw("ListIndexFields", "index not found", req.Partition)
		return nil, status.Errorf(codes.NotFound, "index for partition %s not found", req.Partition)
	}

	fields, err := index.Fields()
	if err != nil {
		s.logger.Errorw("ListIndexFields", "failed to extract index fields", req.Partition, "err", err)
		return nil, fmt.Errorf("failed to extract index fields: %w", err)
	}

	return &proto.ListIndexFieldsResponse{Fields: fields}, nil
}

func (s *server) PrintIndexMapping(ctx context.Context, req *proto.PrintIndexMappingRequest) (*proto.PrintIndexMappingResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	index, exists := s.indexes[req.Partition]
	if !exists {
		s.logger.Errorw("PrintIndexMapping", "index not found", req.Partition)
		return nil, status.Errorf(codes.NotFound, "index for partition %s not found", req.Partition)
	}

	mapping := index.Mapping()

	// Marshal the mapping to JSON
	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		s.logger.Errorw("PrintIndexMapping", "failed to marshal index mapping", req.Partition, "err", err)
		return nil, fmt.Errorf("failed to marshal index mapping: %w", err)
	}

	// Unmarshal the JSON back into a map[string]interface{}
	var mappingMap map[string]interface{}
	if err := json.Unmarshal(mappingJSON, &mappingMap); err != nil {
		s.logger.Errorw("PrintIndexMapping", "failed to unmarshal index mapping", req.Partition, "err", err)
		return nil, fmt.Errorf("failed to unmarshal index mapping: %w", err)
	}

	mappingStruct, err := structpb.NewStruct(mappingMap)
	if err != nil {
		s.logger.Errorw("PrintIndexMapping", "failed to convert mapping to structpb.Struct", req.Partition, "err", err)
		return nil, fmt.Errorf("failed to convert mapping to structpb.Struct: %w", err)
	}

	return &proto.PrintIndexMappingResponse{Mapping: mappingStruct}, nil
}

func (s *server) GetPartitions(ctx context.Context, req *proto.GetPartitionsRequest) (*proto.GetPartitionsResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	partitions := make([]string, 0, len(s.indexes))
	for partition := range s.indexes {
		partitions = append(partitions, partition)
	}

	return &proto.GetPartitionsResponse{Partitions: partitions}, nil
}
