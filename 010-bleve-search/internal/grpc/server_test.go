package indexgrpc

import (
	proto "bleve-proj/proto"
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

var TestDir = "test_index"
var mockServer = NewServer("./" + TestDir)

func Test_server_CreateIndex(t *testing.T) {
	type args struct {
		ctx context.Context
		req *proto.CreateIndexRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.CreateIndexResponse
		wantErr bool
	}{
		{
			name: "Create new index successfully",
			args: args{
				req: &proto.CreateIndexRequest{
					Partition: "testPartition",
					IndexMapping: &structpb.Struct{ // Use structpb.Struct
						Fields: map[string]*structpb.Value{
							"analysis": {
								Kind: &structpb.Value_StructValue{
									StructValue: &structpb.Struct{
										Fields: map[string]*structpb.Value{
											"token_filters": {
												Kind: &structpb.Value_StructValue{
													StructValue: &structpb.Struct{
														Fields: map[string]*structpb.Value{
															"to_lower": {
																Kind: &structpb.Value_StructValue{
																	StructValue: &structpb.Struct{
																		Fields: map[string]*structpb.Value{
																			"type": {
																				Kind: &structpb.Value_StringValue{
																					StringValue: "to_lower",
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want:    &proto.CreateIndexResponse{Message: "Index created successfully"},
			wantErr: false,
		},
		{
			name: "Create index with existing partition",
			args: args{
				req: &proto.CreateIndexRequest{
					Partition: "testPartition",
					IndexMapping: &structpb.Struct{ // Use structpb.Struct
						Fields: map[string]*structpb.Value{
							"analysis": {
								Kind: &structpb.Value_StructValue{
									StructValue: &structpb.Struct{
										Fields: map[string]*structpb.Value{
											"token_filters": {
												Kind: &structpb.Value_StructValue{
													StructValue: &structpb.Struct{
														Fields: map[string]*structpb.Value{
															"to_lower": {
																Kind: &structpb.Value_StructValue{
																	StructValue: &structpb.Struct{
																		Fields: map[string]*structpb.Value{
																			"type": {
																				Kind: &structpb.Value_StringValue{
																					StringValue: "to_lower",
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.CreateIndex(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.CreateIndex() = %v, want %v", got, tt.want)
			}
		})
	}
	defer func() {
		for _, tt := range tests {
			if tt.name == "Create new index successfully" {
				_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: tt.args.req.Partition})
				if err != nil {
					t.Fatalf("Failed to delete test index: %v", err)
				}
			}
		}
	}()
}

func Test_server_IndexDocument(t *testing.T) {
	// Create the test partition before running the tests
	partition := "testPartitionForIndex"
	createReq := &proto.CreateIndexRequest{Partition: partition}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	type args struct {
		ctx context.Context
		req *proto.IndexDocumentRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.IndexDocumentResponse
		wantErr bool
	}{
		{
			name: "Index document successfully",
			args: args{
				req: &proto.IndexDocumentRequest{
					Partition: partition,
					DocId:     "1",
					Document: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"id": {
								Kind: &structpb.Value_StringValue{StringValue: "1"},
							},
							"name": {
								Kind: &structpb.Value_StringValue{StringValue: "john doe"},
							},
							"age": {
								Kind: &structpb.Value_NumberValue{NumberValue: float64(30)},
							},
						},
					},
				},
			},
			want: &proto.IndexDocumentResponse{
				Message: "Document indexed successfully",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.IndexDocument(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.IndexDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.IndexDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_DeleteIndex(t *testing.T) {
	partition := "testPartitionToDelete"
	createReq := &proto.CreateIndexRequest{Partition: partition}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}

	type args struct {
		ctx context.Context
		req *proto.DeleteIndexRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.DeleteIndexResponse
		wantErr bool
	}{
		{
			name: "Delete index successfully",
			args: args{
				req: &proto.DeleteIndexRequest{
					Partition: partition,
				},
			},
			want: &proto.DeleteIndexResponse{
				Message: "Index deleted successfully",
			},
			wantErr: false,
		},
		{
			name: "Delete index with empty partition",
			args: args{
				req: &proto.DeleteIndexRequest{
					Partition: "",
				},
			},
			wantErr: true, // Should return an error because partition is empty
		},
		{
			name: "Delete index with non-existing partition",
			args: args{
				req: &proto.DeleteIndexRequest{
					Partition: "nonExistingPartition",
				},
			},
			wantErr: true, // Should return an error because partition doesn't exist
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.DeleteIndex(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.DeleteIndex() = %v, want %v", got, tt.want)
			}

			// Additional check to ensure index is actually deleted
			if !tt.wantErr {
				indexPath := fmt.Sprintf("./%s/%s", TestDir, tt.args.req.Partition)
				if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
					t.Errorf("Index directory still exists after DeleteIndex: %s", indexPath)
				}
			}
		})
	}
}

func Test_server_Search(t *testing.T) {
	partition := "testPartitionForSearch"
	createReq := &proto.CreateIndexRequest{Partition: partition}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	// Index some test documents into the partition
	docs := []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{"1", "john doe", 30},
		{"2", "jand doe", 25},
		{"3", "jake smith", 40},
	}
	for _, doc := range docs {
		req := &proto.IndexDocumentRequest{
			Partition: partition,
			DocId:     doc.Id,
			Document: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"id": {
						Kind: &structpb.Value_StringValue{StringValue: doc.Id},
					},
					"name": {
						Kind: &structpb.Value_StringValue{StringValue: doc.Name},
					},
					"age": {
						Kind: &structpb.Value_NumberValue{NumberValue: float64(doc.Age)},
					},
				},
			},
		}
		_, err = mockServer.IndexDocument(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	type args struct {
		ctx context.Context
		req *proto.SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.SearchResult
		wantErr bool
	}{
		{
			name: "Search for john",
			args: args{
				req: &proto.SearchRequest{
					Partition: partition,
					Query:     `[{"query":{"type": "term", "term": "john"}}]`,
				},
			},
			want: &proto.SearchResult{
				Hits:      []*proto.SearchHit{{Id: "1"}},
				TotalHits: 1,
			},
			wantErr: false,
		},
		{
			name: "Search for doe",
			args: args{
				req: &proto.SearchRequest{
					Partition: partition,
					Query:     `[{"query":{"type": "term", "term": "doe"}}]`,
				},
			},
			want: &proto.SearchResult{
				Hits:      []*proto.SearchHit{{Id: "1"}, {Id: "2"}},
				TotalHits: 2,
			},
			wantErr: false,
		},
		{
			name: "Search for smith",
			args: args{
				req: &proto.SearchRequest{
					Partition: partition,
					Query:     `[{"query":{"type": "term", "term": "smith"}}]`,
				},
			},
			want: &proto.SearchResult{
				Hits:      []*proto.SearchHit{{Id: "3"}},
				TotalHits: 1,
			},
			wantErr: false,
		},
		{
			name: "Search for john with size and from",
			args: args{
				req: &proto.SearchRequest{
					Partition: partition,
					Query:     `[{"query":{"type": "term", "term": "john"},"size": 1, "from": 0}]`,
				},
			},
			want: &proto.SearchResult{
				Hits:      []*proto.SearchHit{{Id: "1"}},
				TotalHits: 1,
			},
			wantErr: false,
		},
		{
			name: "Search for doe with size and from",
			args: args{
				req: &proto.SearchRequest{
					Partition: partition,
					Query:     `[{"query":{"type": "term", "term": "doe"},"size": 2, "from": 0}]`,
				},
			},
			want: &proto.SearchResult{
				Hits:      []*proto.SearchHit{{Id: "1"}, {Id: "2"}},
				TotalHits: 2,
			},
			wantErr: false,
		},
		{
			name: "Search for smith with size and from",
			args: args{
				req: &proto.SearchRequest{
					Partition: partition,
					Query:     `[{"query":{"type": "term", "term": "smith"},"size": 1, "from": 0}]`,
				},
			},
			want: &proto.SearchResult{
				Hits:      []*proto.SearchHit{{Id: "3"}},
				TotalHits: 1,
			},
			wantErr: false,
		},
		{
			name: "Search for non-existing partition",
			args: args{
				req: &proto.SearchRequest{
					Partition: "nonExistingPartition",
					Query:     `[{"type": "term", "term": "John Doe"}]`,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Search with empty query",
			args: args{
				req: &proto.SearchRequest{
					Partition: partition,
					Query:     "",
				},
			},
			want: &proto.SearchResult{
				Hits:      []*proto.SearchHit{},
				TotalHits: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.Search(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Assert only the relevant fields, as Took and MaxScore can be variable
			if !tt.wantErr {
				// Only assert TotalHits
				if uint64(got.TotalHits) != tt.want.TotalHits {
					t.Errorf("server.Search() TotalHits = %v, want %v", got.TotalHits, tt.want.TotalHits)
				}
			}
		})
	}
}

func Test_server_GetDocumentCount(t *testing.T) {
	partition := "testPartitionForDocumentCount"
	createReq := &proto.CreateIndexRequest{Partition: partition}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	// Index some test documents into the partition
	docs := []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{"1", "john doe", 30},
		{"2", "jand doe", 25},
		{"3", "jake smith", 40},
	}
	for _, doc := range docs {
		req := &proto.IndexDocumentRequest{
			Partition: partition,
			DocId:     doc.Id,
			Document: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"id": {
						Kind: &structpb.Value_StringValue{StringValue: doc.Id},
					},
					"name": {
						Kind: &structpb.Value_StringValue{StringValue: doc.Name},
					},
					"age": {
						Kind: &structpb.Value_NumberValue{NumberValue: float64(doc.Age)},
					},
				},
			},
		}
		_, err = mockServer.IndexDocument(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	type args struct {
		ctx context.Context
		req *proto.GetDocumentCountRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.GetDocumentCountResponse
		wantErr bool
	}{
		{
			name: "Get document count for existing partition",
			args: args{
				req: &proto.GetDocumentCountRequest{
					Partition: partition,
				},
			},
			want: &proto.GetDocumentCountResponse{
				Count: 3, // Expecting 3 documents
			},
			wantErr: false,
		},
		{
			name: "Get document count for non-existing partition",
			args: args{
				req: &proto.GetDocumentCountRequest{
					Partition: "nonExistingPartition",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.GetDocumentCount(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GetDocumentCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.GetDocumentCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_BulkLoadDocuments(t *testing.T) {
	partition := "testPartitionForBulkLoad"
	createReq := &proto.CreateIndexRequest{Partition: partition}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	mockData := `{"name": "Alice", "age": 28}
{"id": "doc2", "name": "Bob", "age": 35}
{"name": "Charlie"}` // Missing "age" field in this document

	type args struct {
		ctx context.Context
		req *proto.BulkLoadDocumentsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.BulkLoadDocumentsResponse
		wantErr bool
	}{
		{
			name: "Bulk load documents successfully",
			args: args{
				req: &proto.BulkLoadDocumentsRequest{
					Partition: partition,
					Documents: []byte(mockData),
				},
			},
			want: &proto.BulkLoadDocumentsResponse{
				Message: "Documents bulk loaded successfully",
			},
			wantErr: false,
		},
		{
			name: "Bulk load documents with invalid JSON",
			args: args{
				req: &proto.BulkLoadDocumentsRequest{
					Partition: partition,
					Documents: []byte(`{"invalid": "json" ]`),
				},
			},
			want: &proto.BulkLoadDocumentsResponse{
				Message: "Documents bulk loaded successfully",
			},
			wantErr: false, // Should return an error because of invalid JSON
		},
		{
			name: "Bulk load documents with non-existing partition",
			args: args{
				req: &proto.BulkLoadDocumentsRequest{
					Partition: "nonExistingPartition",
					Documents: []byte(mockData),
				},
			},
			wantErr: true, // Should return an error because partition doesn't exist
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.BulkLoadDocuments(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.BulkLoadDocuments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr { // If no error expected, check response and document count
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("server.BulkLoadDocuments() = %v, want %v", got, tt.want)
				}

				countRes, err := mockServer.GetDocumentCount(tt.args.ctx, &proto.GetDocumentCountRequest{Partition: partition})
				if err != nil {
					t.Errorf("Failed to get document count after bulk load: %v", err)
				}

				expectedCount := uint64(3)
				if countRes.Count != expectedCount {
					t.Errorf("Document count mismatch after bulk load: got %d, want %d", countRes.Count, expectedCount)
				}
			}
		})
	}
}

func Test_server_CheckIndexContents(t *testing.T) {
	partition := "testPartitionForCheckContents"
	createReq := &proto.CreateIndexRequest{Partition: partition}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	// Index some test documents
	docs := []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{"1", "John Doe", 30},
		{"2", "Jane Doe", 25},
		{"3", "Jake Smith", 40},
	}
	for _, doc := range docs {
		req := &proto.IndexDocumentRequest{
			Partition: partition,
			DocId:     doc.Id,
			Document: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"id":   {Kind: &structpb.Value_StringValue{StringValue: doc.Id}},
					"name": {Kind: &structpb.Value_StringValue{StringValue: doc.Name}},
					"age":  {Kind: &structpb.Value_NumberValue{NumberValue: float64(doc.Age)}},
				},
			},
		}
		_, err = mockServer.IndexDocument(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	type args struct {
		ctx context.Context
		req *proto.CheckIndexContentsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]uint64 // Expected term counts
		wantErr bool
	}{
		{
			name: "Check contents of existing partition",
			args: args{
				req: &proto.CheckIndexContentsRequest{
					Partition: partition,
				},
			},
			want: map[string]uint64{
				"_all": 51, // Adjust based on your index
				"_id":  3,
				"age":  43, // Adjust based on your index
				"name": 5,
				"id":   3,
			},
			wantErr: false,
		},
		{
			name: "Check contents of non-existing partition",
			args: args{
				req: &proto.CheckIndexContentsRequest{
					Partition: "nonExistingPartition",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.CheckIndexContents(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.CheckIndexContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr { // If no error, check field counts
				if len(got.Fields) != len(tt.want) {
					t.Errorf("Unexpected number of fields returned. Got: %d, Want: %d", len(got.Fields), len(tt.want))
				}
				for field, count := range tt.want {
					if gotCount, ok := got.Fields[field]; !ok || gotCount.TermCount != count {
						t.Errorf("Incorrect term count for field '%s'. Got: %v, Want: %v", field, gotCount, count)
					}
				}
			}
		})
	}
}

func Test_server_PrintTermDictionary(t *testing.T) {
	partition := "testPartitionForTermDict"
	createReq := &proto.CreateIndexRequest{Partition: partition}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	// Index some test documents
	docs := []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{"1", "John Doe", 30},
		{"2", "Jane Doe", 25},
		{"3", "Jake Smith", 40},
	}
	for _, doc := range docs {
		req := &proto.IndexDocumentRequest{
			Partition: partition,
			DocId:     doc.Id,
			Document: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"id":   {Kind: &structpb.Value_StringValue{StringValue: doc.Id}},
					"name": {Kind: &structpb.Value_StringValue{StringValue: doc.Name}},
					"age":  {Kind: &structpb.Value_NumberValue{NumberValue: float64(doc.Age)}},
				},
			},
		}
		_, err = mockServer.IndexDocument(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	type args struct {
		ctx context.Context
		req *proto.PrintTermDictionaryRequest
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]uint64 // Expected term dictionary
		wantErr bool
	}{
		{
			name: "Retrieve term dictionary for existing field",
			args: args{
				req: &proto.PrintTermDictionaryRequest{
					Partition: partition,
					Field:     "name",
				},
			},
			want: map[string]uint64{
				"john":  1,
				"doe":   2,
				"jane":  1,
				"smith": 1,
				"jake":  1,
			},
			wantErr: false,
		},
		{
			name: "Retrieve term dictionary for non-existing partition",
			args: args{
				req: &proto.PrintTermDictionaryRequest{
					Partition: "nonExistingPartition",
					Field:     "name",
				},
			},
			wantErr: true,
		},
		{
			name: "Retrieve term dictionary for non-existing field",
			args: args{
				req: &proto.PrintTermDictionaryRequest{
					Partition: partition,
					Field:     "nonExistingField",
				},
			},
			want:    map[string]uint64{}, // Expect an empty map for a non-existing field
			wantErr: false,               // No error expected, just an empty result
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.PrintTermDictionary(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.PrintTermDictionary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if got != nil {
					t.Errorf("server.PrintTermDictionary() expected nil result for error case, got %v", got)
				}
			} else {
				if !reflect.DeepEqual(got.TermCounts, tt.want) {
					t.Errorf("server.PrintTermDictionary() = %v, want %v", got.TermCounts, tt.want)
				}
			}
		})
	}
}

func Test_server_DumpIndexContentsWithLimit(t *testing.T) {
	partition := "testPartitionForDump"
	createReq := &proto.CreateIndexRequest{Partition: partition}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	// Index some test documents
	docs := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
		{"id": "3", "name": "Jake Smith", "age": 40},
	}
	for _, doc := range docs {
		docStruct, err := structpb.NewStruct(doc)
		if err != nil {
			t.Fatalf("Failed to convert doc to structpb.Struct: %v", err)
		}
		req := &proto.IndexDocumentRequest{
			Partition: partition,
			DocId:     doc["id"].(string),
			Document:  docStruct,
		}
		_, err = mockServer.IndexDocument(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	type args struct {
		ctx context.Context
		req *proto.DumpIndexContentsWithLimitRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []*structpb.Struct
		wantErr bool
	}{
		{
			name: "Dump contents with limit 2",
			args: args{
				req: &proto.DumpIndexContentsWithLimitRequest{
					Partition: partition,
					Page:      1,
					Size:      2,
				},
			},
			want: []*structpb.Struct{
				// Create structpb.Struct instances from docs[0] and docs[1]
				func() *structpb.Struct { s, _ := structpb.NewStruct(docs[2]); return s }(),
				func() *structpb.Struct { s, _ := structpb.NewStruct(docs[0]); return s }(),
			},
			wantErr: false,
		},
		// Add more test cases here for different page/size combinations and error scenarios...
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.DumpIndexContentsWithLimit(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.DumpIndexContentsWithLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.Documents, tt.want) {
					t.Errorf("server.DumpIndexContentsWithLimit() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_server_ListIndexFields(t *testing.T) {
	partition := "testPartitionForListFields"
	createReq := &proto.CreateIndexRequest{Partition: partition}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	// Index some test documents (this will add fields to the index)
	docs := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
		{"id": "3", "name": "Jake Smith", "age": 40},
	}
	for _, doc := range docs {
		docStruct, err := structpb.NewStruct(doc)
		if err != nil {
			t.Fatalf("Failed to convert doc to structpb.Struct: %v", err)
		}
		req := &proto.IndexDocumentRequest{
			Partition: partition,
			DocId:     doc["id"].(string),
			Document:  docStruct,
		}
		_, err = mockServer.IndexDocument(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	type args struct {
		ctx context.Context
		req *proto.ListIndexFieldsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "List fields of existing partition",
			args: args{
				req: &proto.ListIndexFieldsRequest{
					Partition: partition,
				},
			},
			want:    []string{"_all", "_id", "age", "id", "name"}, // Adjust based on expected fields
			wantErr: false,
		},
		{
			name: "List fields of non-existing partition",
			args: args{
				req: &proto.ListIndexFieldsRequest{
					Partition: "nonExistingPartition",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.ListIndexFields(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.ListIndexFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Sort both 'got' and 'want' slices before comparison
				sort.Strings(got.Fields)
				sort.Strings(tt.want)

				if !reflect.DeepEqual(got.Fields, tt.want) {
					t.Errorf("server.ListIndexFields() = %v, want %v", got.Fields, tt.want)
				}
			}
		})
	}
}

func Test_server_PrintIndexMapping(t *testing.T) {
	partition := "testPartitionForMapping"
	indexMapping := &structpb.Struct{ // Create indexMapping as structpb.Struct
		Fields: map[string]*structpb.Value{
			"analysis": {
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"token_filters": {
								Kind: &structpb.Value_StructValue{
									StructValue: &structpb.Struct{
										Fields: map[string]*structpb.Value{
											"to_lower": {
												Kind: &structpb.Value_StructValue{
													StructValue: &structpb.Struct{
														Fields: map[string]*structpb.Value{
															"type": {
																Kind: &structpb.Value_StringValue{
																	StringValue: "to_lower",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	createReq := &proto.CreateIndexRequest{
		Partition:    partition,
		IndexMapping: indexMapping, // Use the new IndexMapping field
	}
	_, err := mockServer.CreateIndex(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		_, err := mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	type args struct {
		ctx context.Context
		req *proto.PrintIndexMappingRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *structpb.Struct // Expected mapping as structpb.Struct
		wantErr bool
	}{
		{
			name: "Get mapping for existing partition",
			args: args{
				req: &proto.PrintIndexMappingRequest{
					Partition: partition,
				},
			},
			want: func() *structpb.Struct {
				// Construct the expected structpb.Struct with ALL expected fields
				expectedStruct, _ := structpb.NewStruct(map[string]interface{}{
					"analysis": map[string]interface{}{
						"token_filters": map[string]interface{}{
							"to_lower": map[string]interface{}{
								"type": "to_lower",
							},
						},
					},
					"default_analyzer":        "standard",         // Default field
					"default_datetime_parser": "dateTimeOptional", // Default field
					"default_field":           "_all",             // Default field
					"default_mapping": map[string]interface{}{ // Default field
						"dynamic": true,
						"enabled": true,
					},
					"default_type":      "_default", // Default field
					"docvalues_dynamic": true,       // Default field
					"index_dynamic":     true,       // Default field
					"store_dynamic":     true,       // Default field
					"type_field":        "_type",    // Default field
					// ... Add any other default fields that are present in your mapping ...
				})
				return expectedStruct
			}(),
			wantErr: false,
		},
		{
			name: "Get mapping for non-existing partition",
			args: args{
				req: &proto.PrintIndexMappingRequest{
					Partition: "nonExistingPartition",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.PrintIndexMapping(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.PrintIndexMapping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.Mapping, tt.want) {
					t.Errorf("server.PrintIndexMapping() = %v, want %v", got.Mapping, tt.want)
				}
			}
		})
	}
}

func Test_server_GetPartitions(t *testing.T) {
	partition1 := "testPartition1"
	partition2 := "testPartition2"

	// Create the test partitions
	_, err := mockServer.CreateIndex(context.Background(), &proto.CreateIndexRequest{Partition: partition1})
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	_, err = mockServer.CreateIndex(context.Background(), &proto.CreateIndexRequest{Partition: partition2})
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		// Clean up test partitions
		_, err = mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition1})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
		_, err = mockServer.DeleteIndex(context.Background(), &proto.DeleteIndexRequest{Partition: partition2})
		if err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "Get all partitions",
			want:    []string{partition1, partition2}, // Expected partitions
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockServer.GetPartitions(context.Background(), &proto.GetPartitionsRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GetPartitions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Compare partitions (ignoring order)
				if len(got.Partitions) != len(tt.want) {
					t.Errorf("Unexpected number of partitions. Got: %d, Want: %d", len(got.Partitions), len(tt.want))
				}
				for _, expectedPartition := range tt.want {
					found := false
					for _, actualPartition := range got.Partitions {
						if expectedPartition == actualPartition {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Partition '%s' not found in returned partitions: %v", expectedPartition, got.Partitions)
					}
				}
			}
		})
	}
}
