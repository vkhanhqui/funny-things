package repositories

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis"
	_ "github.com/blevesearch/bleve/v2/analysis/token/apostrophe"
	_ "github.com/blevesearch/bleve/v2/analysis/token/camelcase"
	_ "github.com/blevesearch/bleve/v2/analysis/token/compound"
	_ "github.com/blevesearch/bleve/v2/analysis/token/edgengram"
	_ "github.com/blevesearch/bleve/v2/analysis/token/elision"
	_ "github.com/blevesearch/bleve/v2/analysis/token/hierarchy"
	_ "github.com/blevesearch/bleve/v2/analysis/token/keyword"
	_ "github.com/blevesearch/bleve/v2/analysis/token/length"
	_ "github.com/blevesearch/bleve/v2/analysis/token/ngram"
	_ "github.com/blevesearch/bleve/v2/analysis/token/porter"
	_ "github.com/blevesearch/bleve/v2/analysis/token/reverse"
	_ "github.com/blevesearch/bleve/v2/analysis/token/shingle"
	_ "github.com/blevesearch/bleve/v2/analysis/token/snowball"
	_ "github.com/blevesearch/bleve/v2/analysis/token/stop"
	_ "github.com/blevesearch/bleve/v2/analysis/token/truncate"
	_ "github.com/blevesearch/bleve/v2/analysis/token/unicodenorm"
	_ "github.com/blevesearch/bleve/v2/analysis/token/unique"
	_ "github.com/blevesearch/bleve/v2/analysis/tokenmap"
	"github.com/blevesearch/bleve/v2/registry"
)

var globalRepo IndexRepository
var defaultIndex map[string]interface{}
var TestDir, CustomPartition string

func init() {
	TestDir = "test_index"
	CustomPartition = "customPartition"
	dictListConfig := map[string]interface{}{
		"type":   "custom",
		"tokens": []interface{}{"factor", "soft", "ball", "team"},
	}
	cache := registry.NewCache()
	tokenMap, _ := cache.DefineTokenMap("dict_test", dictListConfig)
	registry.RegisterTokenMap("dict_test", func(config map[string]interface{}, cache *registry.Cache) (analysis.TokenMap, error) {
		// Define your token map logic here
		return tokenMap, nil
	})

	articleListConfig := map[string]interface{}{
		"type":   "custom",
		"tokens": []interface{}{"ar"},
	}
	tokenMap, _ = cache.DefineTokenMap("articles_test", articleListConfig)
	registry.RegisterTokenMap("articles_test", func(config map[string]interface{}, cache *registry.Cache) (analysis.TokenMap, error) {
		// Define your token map logic here
		return tokenMap, nil
	})

}

func TestMain(m *testing.M) {
	globalRepo = NewIndexRepository("./" + TestDir)
	defaultIndex = map[string]interface{}{
		"default_analyzer": "standard",
	}
	code := m.Run()
	os.Exit(code)
}

func Test_indexRepository_CreateIndex(t *testing.T) {
	type args struct {
		partition string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Create new index successfully",
			args:    args{partition: "testPartition"},
			wantErr: false,
		},
		{
			name:    "Create index with existing partition",
			args:    args{partition: "testPartition"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := globalRepo.CreateIndex(tt.args.partition, defaultIndex)

			if (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("indexRepository.CreateIndex() expected valid index, got nil")
			}
		})
	}
	for _, tt := range tests {
		if tt.name == "Create new index successfully" {
			if err := globalRepo.DeleteIndex(tt.args.partition); err != nil {
				t.Errorf("indexRepository.DeleteIndex() error = %v during cleanup", err)
			}
		}
	}
}

func Test_indexRepository_DeleteIndex(t *testing.T) {
	type args struct {
		partition string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Delete existing index successfully",
			args:    args{partition: "testPartition"},
			wantErr: false,
		},
		{
			name:    "Attempt to delete non-existent index",
			args:    args{partition: "nonExistentPartition"},
			wantErr: true,
		},
	}

	// Create an index for the first test case
	_, err := globalRepo.CreateIndex("testPartition", defaultIndex)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := globalRepo.DeleteIndex(tt.args.partition); (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_indexRepository_OpenIndex(t *testing.T) {
	type args struct {
		partition string
	}
	tests := []struct {
		name    string
		args    args
		want    bleve.Index
		wantErr bool
	}{
		{
			name: "Open existing index",
			args: args{
				partition: "existingPartition",
			},
			wantErr: false,
		},
		{
			name: "Open non-existent index",
			args: args{
				partition: "nonExistentPartition",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		// Create the index for the "Open existing index" test case
		if tt.name == "Open existing index" {
			if _, err := globalRepo.CreateIndex(tt.args.partition, defaultIndex); err != nil {
				t.Fatalf("Failed to create test index: %v", err)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if _, err := globalRepo.OpenIndex(tt.args.partition); (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.OpenIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Cleanup: delete the created index using the repository method
			if !tt.wantErr {
				if err := globalRepo.DeleteIndex(tt.args.partition); err != nil {
					t.Errorf("indexRepository.DeleteIndex() error = %v during cleanup", err)
				}
			}
		})
	}
}

func Test_indexRepository_IndexDocument(t *testing.T) {
	type args struct {
		partition string
		docID     string
		document  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Index new document successfully",
			args: args{
				partition: "testPartition",
				docID:     "doc1",
				document: map[string]interface{}{
					"field1": "value1",
					"field2": 123,
				},
			},
			wantErr: false,
		},
		{
			name: "Index document with existing ID",
			args: args{
				partition: "testPartition",
				docID:     "doc2",
				document: map[string]interface{}{
					"field1": "value2",
					"field2": 456,
				},
			},
			wantErr: false,
		},
		{
			name: "Index document with missing partition",
			args: args{
				partition: "nonExistingPartition",
				docID:     "doc3",
				document: map[string]interface{}{
					"field1": "value3",
					"field2": 789,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure index exists for the partition
			if !tt.wantErr {
				if _, err := globalRepo.CreateIndex(tt.args.partition, defaultIndex); err != nil {
					t.Fatalf("Failed to create index for partition %s: %v", tt.args.partition, err)
				}
				defer func() {
					if err := globalRepo.DeleteIndex(tt.args.partition); err != nil {
						t.Errorf("Failed to delete index for partition %s: %v", tt.args.partition, err)
					}
				}()
			}

			// Index the document based on the test case arguments using globalRepo
			err := globalRepo.IndexDocument(tt.args.partition, tt.args.docID, tt.args.document)
			// Check if error matches expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("globalRepo.IndexDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_indexRepository_Search(t *testing.T) {
	// Create the test partition before running the tests
	partition := "testPartition"
	_, err := globalRepo.CreateIndex(partition, defaultIndex)
	if err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		if err := globalRepo.DeleteIndex(partition); err != nil {
			t.Fatalf("Failed to create test index: %v", err)
		}
	}()

	// Index some test documents into the partition
	docs := []map[string]interface{}{
		{"id": "1", "name": "john doe", "age": 30},
		{"id": "2", "name": "jand doe", "age": 25},
		{"id": "3", "name": "jake smith", "age": 40},
	}
	for _, doc := range docs {
		err := globalRepo.IndexDocument(partition, doc["id"].(string), doc)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	tests := []struct {
		name       string
		partition  string
		query      string
		wantResult []string
		wantErr    bool
	}{
		{
			name:       "Search for john",
			partition:  "testPartition",
			query:      `[{"query":{"type": "term", "term": "john"}}]`,
			wantResult: []string{"1"},
			wantErr:    false,
		},
		{
			name:       "Search for doe",
			partition:  "testPartition",
			query:      `[{"query":{"type": "term", "term": "doe"}}]`,
			wantResult: []string{"1", "2"},
			wantErr:    false,
		},
		{
			name:       "Search for smith",
			partition:  "testPartition",
			query:      `[{"query":{"type": "term", "term": "smith"}}]`,
			wantResult: []string{"3"},
			wantErr:    false,
		},
		{
			name:       "Search for john with size and from",
			partition:  "testPartition",
			query:      `[{"query":{"type": "term", "term": "john"},"size": 1, "from": 0}]`,
			wantResult: []string{"1"},
			wantErr:    false,
		},
		{
			name:       "Search for doe with size and from",
			partition:  "testPartition",
			query:      `[{"query":{"type": "term", "term": "doe"},"size": 2, "from": 0}]`,
			wantResult: []string{"1", "2"},
			wantErr:    false,
		},
		{
			name:       "Search for smith with size and from",
			partition:  "testPartition",
			query:      `[{"query":{"type": "term", "term": "smith"},"size": 1, "from": 0}]`,
			wantResult: []string{"3"},
			wantErr:    false,
		},
		{
			name:       "Search for non-existing partition",
			partition:  "nonExistingPartition",
			query:      `[{"type": "term", "term": "John Doe"}]`,
			wantResult: nil,
			wantErr:    true,
		},

		{
			name:       "Search with empty query",
			partition:  "testPartition",
			query:      "",
			wantResult: []string{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := globalRepo.Search(tt.partition, tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				gotResult := []string{}
				for _, hit := range got.Hits {
					gotResult = append(gotResult, hit.ID)
				}
				if !reflect.DeepEqual(gotResult, tt.wantResult) {
					t.Errorf("indexRepository.Search() = %v, want %v", gotResult, tt.wantResult)
				}
			}
		})
	}
}
func Test_indexRepository_GetDocumentCount(t *testing.T) {
	partition := "testPartition"
	// Create the test partition before running the tests
	if _, err := globalRepo.CreateIndex(partition, defaultIndex); err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		if err := globalRepo.DeleteIndex(partition); err != nil {
			t.Fatalf("Failed to create test index: %v", err)
		}
	}()

	// Index some test documents into the partition
	docs := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
		{"id": "3", "name": "Jake Smith", "age": 40},
	}
	for _, doc := range docs {
		err := globalRepo.IndexDocument("testPartition", doc["id"].(string), doc)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	// Test cases for GetDocumentCount
	tests := []struct {
		name      string
		partition string
		want      uint64
		wantErr   bool
	}{
		{
			name:      "Get document count for existing partition",
			partition: "testPartition",
			want:      3,
			wantErr:   false,
		},
		{
			name:      "Get document count for non-existing partition",
			partition: "nonExistingPartition",
			want:      0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := globalRepo.GetDocumentCount(tt.partition)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.GetDocumentCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("indexRepository.GetDocumentCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexRepository_BulkLoadDocuments(t *testing.T) {
	partition := "testPartition"
	// Create the test partition before running the tests
	if _, err := globalRepo.CreateIndex(partition, defaultIndex); err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		if err := globalRepo.DeleteIndex(partition); err != nil {
			t.Fatalf("Failed to create test index: %v", err)
		}
	}()
	// Prepare mock data for the scanner
	mockData := `{"id": "1", "name": "John Doe", "age": 30}
{"id": "2", "name": "Jane Smith", "age": 25}
{"id": "3", "name": "Jake Johnson", "age": 40}`

	scanner := bufio.NewScanner(strings.NewReader(mockData))

	// Test cases for BulkLoadDocuments
	tests := []struct {
		name      string
		partition string
		scanner   *bufio.Scanner
		wantErr   bool
	}{
		{
			name:      "Bulk load documents into existing partition",
			partition: "testPartition",
			scanner:   scanner,
			wantErr:   false,
		},
		{
			name:      "Bulk load documents into non-existing partition",
			partition: "nonExistingPartition",
			scanner:   scanner,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := globalRepo.BulkLoadDocuments(tt.partition, tt.scanner)
			if err != nil && !tt.wantErr {
				t.Errorf("indexRepository.BulkLoadDocuments() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.wantErr {
				t.Errorf("indexRepository.BulkLoadDocuments() expected error, got nil")
			}

			// Optionally, verify the document count after bulk loading
			if !tt.wantErr {
				count, err := globalRepo.GetDocumentCount(tt.partition)
				if err != nil {
					t.Errorf("Failed to get document count after bulk load: %v", err)
				}
				expectedCount := uint64(3) // Since we have 3 documents in mockData
				if count != expectedCount {
					t.Errorf("Document count mismatch after bulk load: got %d, want %d", count, expectedCount)
				}
			}
		})
	}
}

func Test_indexRepository_CheckIndexContents(t *testing.T) {
	partition := "testPartition"
	// Create the test partition before running the tests
	if _, err := globalRepo.CreateIndex(partition, defaultIndex); err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		if err := globalRepo.DeleteIndex(partition); err != nil {
			t.Fatalf("Failed to create test index: %v", err)
		}
	}()

	// Index some test documents into the partition
	docs := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
		{"id": "3", "name": "Jake Smith", "age": 40},
	}
	for _, doc := range docs {
		err := globalRepo.IndexDocument("testPartition", doc["id"].(string), doc)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	tests := []struct {
		name      string
		partition string
		want      map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "Check contents of existing partition",
			partition: "testPartition",
			wantErr:   false,
			want: map[string]interface{}{
				"_all": map[string]int{
					"termCount": 51,
				},
				"_id": map[string]int{
					"termCount": 3,
				},
				"age": map[string]int{
					"termCount": 43,
				},
				"name": map[string]int{
					"termCount": 5,
				},
				"id": map[string]int{
					"termCount": 3,
				},
			},
		},
		{
			name:      "Check contents of non-existing partition",
			partition: "nonExistingPartition",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := globalRepo.CheckIndexContents(tt.partition)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.CheckIndexContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				// If we expect an error, got should be nil
				if got != nil {
					t.Errorf("indexRepository.CheckIndexContents() expected nil result for non-existing partition, got %v", got)
				} else {
					// Compare the contents of the maps
					if !mapsEqual(got, tt.want) {
						t.Errorf("indexRepository.CheckIndexContents() = %v, want %v", got, tt.want)
					}
				}

			}
		})
	}
}

func mapsEqual(got, want map[string]interface{}) bool {
	if len(got) != len(want) {
		return false
	}

	for key, wantValue := range want {
		gotValue, ok := got[key]
		if !ok {
			return false
		}

		wantMap, wantIsMap := wantValue.(map[string]int)
		gotMap, gotIsMap := gotValue.(map[string]interface{})
		if wantIsMap && gotIsMap {
			// Both are maps, check each key-value pair
			if !intMapsEqual(gotMap, wantMap) {
				return false
			}
		} else {
			// Otherwise, directly compare values
			if gotValue != wantValue {
				return false
			}
		}
	}

	return true
}

func intMapsEqual(got map[string]interface{}, want map[string]int) bool {
	if len(got) != len(want) {
		return false
	}

	for key, wantValue := range want {
		gotValue, ok := got[key]
		if !ok {
			return false
		}

		gotInt, gotIsInt := gotValue.(int)
		if !gotIsInt || gotInt != wantValue {
			return false
		}
	}

	return true
}

func Test_indexRepository_PrintTermDictionary(t *testing.T) {
	partition := "testPartition"
	// Create the test partition before running the tests
	if _, err := globalRepo.CreateIndex(partition, defaultIndex); err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		if err := globalRepo.DeleteIndex(partition); err != nil {
			t.Fatalf("Failed to create test index: %v", err)
		}
	}()
	// Index some test documents into the partition
	docs := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
		{"id": "3", "name": "Jake Smith", "age": 40},
	}
	for _, doc := range docs {
		err := globalRepo.IndexDocument("testPartition", doc["id"].(string), doc)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}

	tests := []struct {
		partition string
		name      string
		want      map[string]uint64
		wantErr   bool
	}{
		{
			partition: "testPartition",
			name:      "Retrieve term dictionary for existing field",
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
			partition: "nonExistingPartition",
			name:      "Retrieve term dictionary for non-existing partition",
			wantErr:   true,
		},
		{
			partition: "testPartition",
			name:      "Retrieve term dictionary for non-existing field",
			want: map[string]uint64{
				"john":  1,
				"doe":   2,
				"jane":  1,
				"smith": 1,
				"jake":  1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := globalRepo.PrintTermDictionary(tt.partition, "name")
			if (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.PrintTermDictionary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if got != nil {
					t.Errorf("indexRepository.PrintTermDictionary() expected nil result for error case, got %v", got)
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("indexRepository.PrintTermDictionary() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
func Test_indexRepository_DumpIndexContentsWithLimit(t *testing.T) {

	partition := "testPartition"
	// Create the test partition before running the tests
	if _, err := globalRepo.CreateIndex(partition, defaultIndex); err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		if err := globalRepo.DeleteIndex(partition); err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()

	// Index some test documents into the partition
	docs := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
		{"id": "3", "name": "Jake Smith", "age": 40},
	}
	for _, doc := range docs {
		err := globalRepo.IndexDocument(partition, doc["id"].(string), doc)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}
	compareMapsSlice := func(got, want []interface{}) bool {
		if len(got) != len(want) {
			return false
		}

		for i := range got {
			gotMap, gotOK := got[i].(map[string]interface{})
			wantMap, wantOK := want[i].(map[string]interface{})
			if !gotOK || !wantOK {
				return false
			}

			// Sort keys for both maps
			gotKeys := make([]string, 0, len(gotMap))
			wantKeys := make([]string, 0, len(wantMap))

			for key := range gotMap {
				gotKeys = append(gotKeys, key)
			}
			sort.Strings(gotKeys)

			for key := range wantMap {
				wantKeys = append(wantKeys, key)
			}
			sort.Strings(wantKeys)

			// Compare sorted key-value pairs
			if !reflect.DeepEqual(gotKeys, wantKeys) {
				return false
			}

			for _, key := range gotKeys {
				gotValue := fmt.Sprintf("%v", gotMap[key])   // Convert to string representation
				wantValue := fmt.Sprintf("%v", wantMap[key]) // Convert to string representation

				if gotValue != wantValue {
					return false
				}
			}
		}

		return true
	}
	type args struct {
		partition string
		page      int
		size      int
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "Dump contents with limit 2",
			args: args{
				partition: partition,
				page:      1,
				size:      2,
			},
			wantErr: false,
			want: []interface{}{
				docs[2], // {"id": "1", "name": "John Doe", "age": 30}
				docs[0], // {"id": "2", "name": "Jane Doe", "age": 25}
			},
		},
		{
			name: "Dump contents with limit 3",
			args: args{
				partition: partition,
				page:      1,
				size:      3,
			},
			wantErr: false,
			want: []interface{}{
				docs[2], // {"id": "3", "name": "Jake Smith", "age": 40}
				docs[0], // {"id": "1", "name": "John Doe", "age": 30}
				docs[1], // {"id": "2", "name": "Jane Doe", "age": 25}
			},
		},
		{
			name: "Dump contents with limit 1 and offset 2",
			args: args{
				partition: partition,
				page:      2,
				size:      1,
			},
			wantErr: false,
			want: []interface{}{
				docs[0], // {"id": "1", "name": "John Doe", "age": 30}
			},
		},
		{
			name: "Dump contents with invalid offset",
			args: args{
				partition: partition,
				page:      -1,
				size:      10,
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Dump contents with invalid limit",
			args: args{
				partition: partition,
				page:      0,
				size:      -1,
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Dump contents with non-existing partition",
			args: args{
				partition: "nonExistingPartition",
				page:      0,
				size:      10,
			},
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {

		// Calculate offset and limit
		t.Run(tt.name, func(t *testing.T) {
			// Calculate offset and limit
			got, err := globalRepo.DumpIndexContentsWithLimit(tt.args.partition, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.DumpIndexContentsWithLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareMapsSlice(got, tt.want) {
				t.Errorf("indexRepository.DumpIndexContentsWithLimit() results mismatch: got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexRepository_ListIndexFields(t *testing.T) {
	// Test cases
	tests := []struct {
		name      string
		partition string
		want      []string
		wantErr   bool
	}{
		{
			name:      "List fields of existing partition",
			partition: "testPartition",
			want:      []string{"_id", "_all", "age", "id", "name"},
			wantErr:   false,
		},
		{
			name:      "List fields of non-existing partition",
			partition: "nonExistingPartition",
			want:      nil,
			wantErr:   true,
		},
	}

	partition := "testPartition"
	// Create the test partition before running the tests
	if _, err := globalRepo.CreateIndex(partition, defaultIndex); err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		if err := globalRepo.DeleteIndex(partition); err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
	}()
	// Index some test documents into the partition
	docs := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
		{"id": "3", "name": "Jake Smith", "age": 40},
	}
	for _, doc := range docs {
		err := globalRepo.IndexDocument(partition, doc["id"].(string), doc)
		if err != nil {
			t.Fatalf("Failed to index document: %v", err)
		}
	}
	compareList := func(list1, list2 []string) bool {
		if len(list1) != len(list2) {
			return false
		}

		// Count occurrences of each string in list1
		count1 := make(map[string]int)
		for _, s := range list1 {
			count1[s]++
		}

		// Count occurrences of each string in list2
		count2 := make(map[string]int)
		for _, s := range list2 {
			count2[s]++
		}

		// Compare counts to check if they are equal
		return reflect.DeepEqual(count1, count2)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Call the method under test
			got, err := globalRepo.ListIndexFields(tt.partition)

			// Check error condition
			if (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.ListIndexFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !compareList(got, tt.want) {
				t.Errorf("indexRepository.ListIndexFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_indexRepository_PrintIndexMapping(t *testing.T) {
	// Assuming globalRepo is your instance of indexRepository
	type args struct {
		partition string
	}
	tests := []struct {
		name         string
		args         args
		want         map[string]interface{}
		wantErr      bool
		indexMapping map[string]interface{}
	}{
		{
			name:    "Get mapping for existing partition with custom mappings",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"to_lower": map[string]interface{}{
							"type": "to_lower",
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"to_lower": map[string]interface{}{
							"type": "to_lower",
						},
					},
				},
			},
		},
		{
			name:    "Get mapping for existing partition with 'apostrophe' token filter",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"apostrophe": map[string]interface{}{
							"type": "apostrophe",
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"apostrophe": map[string]interface{}{
							"type": "apostrophe",
						},
					},
				},
			},
		},
		{
			name:    "Get mapping for existing partition with 'camelCase' token filter",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"camelCase": map[string]interface{}{
							"type": "camelCase",
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"camelCase": map[string]interface{}{
							"type": "camelCase",
						},
					},
				},
			},
		},
		{
			name:    "Get mapping for existing partition with 'dict_compound' token filter",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"dict_compound": map[string]interface{}{
							"type":               "dict_compound",
							"dict_token_map":     "dict_test",
							"min_word_size":      5,
							"min_subword_size":   2,
							"max_subword_size":   15,
							"only_longest_match": false,
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"dict_compound": map[string]interface{}{
							"type":               "dict_compound",
							"dict_token_map":     "dict_test", // Replace with your actual token map name
							"min_word_size":      5,
							"min_subword_size":   2,
							"max_subword_size":   15,
							"only_longest_match": false,
						},
					},
				},
			},
		},
		{
			name:    "Front edge ngrams",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"edge_ngram": map[string]interface{}{
							"type": "edge_ngram",
							"back": false,
							"min":  1,
							"max":  3,
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"edge_ngram": map[string]interface{}{
							"type": "edge_ngram",
							"back": false,
							"min":  1,
							"max":  3,
						},
					},
				},
			},
		},
		{
			name:    "Back edge ngrams",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"edge_ngram": map[string]interface{}{
							"type": "edge_ngram",
							"back": true,
							"min":  1,
							"max":  3,
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"edge_ngram": map[string]interface{}{
							"type": "edge_ngram",
							"back": true,
							"min":  1,
							"max":  3,
						},
					},
				},
			},
		},
		{
			name:    "Remove apostrophe prefixed word",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"elision": map[string]interface{}{
							"type":               "elision",
							"articles_token_map": "articles_test",
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"elision": map[string]interface{}{
							"type":               "elision",
							"articles_token_map": "articles_test",
						},
					},
				},
			},
		},
		{
			name:    "Remove apostrophe prefixed word",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"elision": map[string]interface{}{
							"type":               "elision",
							"articles_token_map": "articles_test",
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"elision": map[string]interface{}{
							"type":               "elision",
							"articles_token_map": "articles_test",
						},
					},
				},
			},
		},
		{
			name:    "Single token a/b/c, delimiter /, split",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"hierarchy": map[string]interface{}{
							"type":        "hierarchy",
							"delimiter":   "/",
							"max_levels":  10,
							"split_input": true,
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"hierarchy": map[string]interface{}{
							"type":        "hierarchy",
							"delimiter":   "/",
							"max_levels":  10,
							"split_input": true,
						},
					},
				},
			},
		},
		{
			name:    "Multiple tokens a b c, delimiter /, no split",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"hierarchy": map[string]interface{}{
							"type":        "hierarchy",
							"delimiter":   "/",
							"max_levels":  10,
							"split_input": false,
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"hierarchy": map[string]interface{}{
							"type":        "hierarchy",
							"delimiter":   "/",
							"max_levels":  10,
							"split_input": false,
						},
					},
				},
			},
		},
		{
			name:    "Get mapping for existing partition with custom mappings",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"keyword_marker": map[string]interface{}{
							"type":               "keyword_marker",
							"keywords_token_map": "dict_test",
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"keyword_marker": map[string]interface{}{
							"type":               "keyword_marker",
							"keywords_token_map": "dict_test",
						},
					},
				},
			},
		},
		{
			name:    "Length filter with min and max constraints",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"length": map[string]interface{}{
							"type": "length",
							"min":  3,
							"max":  6,
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"length": map[string]interface{}{
							"type": "length",
							"min":  3,
							"max":  6,
						},
					},
				},
			},
		},
		{
			name:    "Front edge ngrams",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"ngram": map[string]interface{}{
							"type": "ngram",
							"min":  1,
							"max":  3,
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"ngram": map[string]interface{}{
							"type": "ngram",
							"min":  1,
							"max":  3,
						},
					},
				},
			},
		},
		{
			name:    "Porter Stemmer",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"stemmer_porter": map[string]interface{}{
							"type": "stemmer_porter",
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"stemmer_porter": map[string]interface{}{
							"type": "stemmer_porter",
						},
					},
				},
			},
		},
		{
			name:    "Shingle Filter Test Case",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"shingle": map[string]interface{}{
							"type": "shingle",
							// Add other shingle filter specific parameters here
							"min": 2,
							"max": 4,
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"shingle": map[string]interface{}{
							"type": "shingle",
							// Add other shingle filter specific parameters here
							"min": 2,
							"max": 4,
						},
					},
				},
			},
		},
		{
			name:    "Snowball Stemmer Test Case",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"stemmer_snowball": map[string]interface{}{
							"type":     "stemmer_snowball",
							"language": "english", // Adjust language if needed
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"stemmer_snowball": map[string]interface{}{
							"type":     "stemmer_snowball",
							"language": "english", // Adjust language if needed
						},
					},
				},
			},
		},
		{
			name:    "Stop Tokens Filter Test Case",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"stop_tokens": map[string]interface{}{
							"type":           "stop_tokens",
							"stop_token_map": "dict_test", // Adjust if needed
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"stop_tokens": map[string]interface{}{
							"type":           "stop_tokens",
							"stop_token_map": "dict_test", // Adjust if needed
						},
					},
				},
			},
		},
		{
			name:    "Truncate Token Filter Test Case",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"truncate_token": map[string]interface{}{
							"type":   "truncate_token",
							"length": 5, // Adjusted length value
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"truncate_token": map[string]interface{}{
							"type":   "truncate_token",
							"length": 5, // Adjusted length value
						},
					},
				},
			},
		},
		{
			name:    "Unicode Normalize Filter Test Case",
			args:    args{partition: CustomPartition},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"normalize_unicode": map[string]interface{}{
							"type": "normalize_unicode",
							"form": "nfc", // Adjusted form value
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"normalize_unicode": map[string]interface{}{
							"type": "normalize_unicode",
							"form": "nfc", // Adjusted form value
						},
					},
				},
			},
		},
		{
			name:    "Unique Term Filter Test Case",
			args:    args{},
			wantErr: false,
			want: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"unique": map[string]interface{}{
							"type": "unique",
						},
					},
				},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"analysis": map[string]interface{}{
					"token_filters": map[string]interface{}{
						"unique": map[string]interface{}{
							"type": "unique",
						},
					},
				},
			},
		},
		{
			name:    "Unique Term Filter Test Case",
			args:    args{},
			wantErr: false,
			want: map[string]interface{}{
				"analysis":                map[string]interface{}{},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "_all",
				"default_mapping": map[string]interface{}{
					"default_analyzer": "standard",
					"dynamic":          true,
					"enabled":          true,
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"dynamic": true,
							"enabled": true,
							"fields": []map[string]interface{}{
								{
									"analyzer":                   "standard",
									"dims":                       1,
									"docvalues":                  true,
									"include_in_all":             true,
									"include_term_vectors":       true,
									"index":                      true,
									"name":                       "test_field",
									"similarity":                 "test_field",
									"skip_freq_norm":             true,
									"store":                      true,
									"type":                       "text",
									"vector_index_optimized_for": "test_field",
								},
							},
						},
					},
				},
				"default_type":      "_default",
				"docvalues_dynamic": true,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
			},
			indexMapping: map[string]interface{}{
				"default_mapping": map[string]interface{}{
					"enabled":          true,
					"dynamic":          true,
					"default_analyzer": "standard",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"fields": []map[string]interface{}{
								{
									"name":                       "test_field",
									"type":                       "text",
									"analyzer":                   "standard",
									"store":                      true,
									"index":                      true,
									"include_term_vectors":       true,
									"include_in_all":             true,
									"docvalues":                  true,
									"skip_freq_norm":             true,
									"dims":                       1,
									"similarity":                 "test_field",
									"vector_index_optimized_for": "test_field",
								},
							},
						},
					},
				},
			},
		},
		{
			name:         "Non-existing partition",
			args:         args{partition: "nonExistingPartition"},
			wantErr:      true,
			want:         nil, // No mapping expected if partition doesn't exist
			indexMapping: nil,
		},
	}
	compareFn := func(got, want map[string]interface{}) bool {
		gotBytes, err := json.Marshal(got)
		if err != nil {
			return false
		}

		wantBytes, err := json.Marshal(want)
		if err != nil {
			return false
		}

		return reflect.DeepEqual(gotBytes, wantBytes)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.indexMapping != nil {
				// Create the test partition before running the tests
				if _, err := globalRepo.CreateIndex(tt.args.partition, tt.indexMapping); err != nil {
					t.Fatalf("Failed to create test index: %v", err)
				}
				defer func() {
					if err := globalRepo.DeleteIndex(tt.args.partition); err != nil {
						t.Fatalf("Failed to delete test index: %v", err)
					}
				}()
			}

			got, err := globalRepo.PrintIndexMapping(tt.args.partition)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexRepository.PrintIndexMapping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareFn(got, tt.want) {
				t.Errorf("indexRepository.PrintIndexMapping() = %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_indexRepository_GetPartitions(t *testing.T) {
	// Initialize a mock index repository with some test partitions
	partition := "testPartition"
	// Create the test partition before running the tests
	if _, err := globalRepo.CreateIndex(partition, defaultIndex); err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}

	partition1 := "testPartition1"
	// Create the test partition before running the tests
	if _, err := globalRepo.CreateIndex(partition1, defaultIndex); err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	defer func() {
		if err := globalRepo.DeleteIndex(partition); err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
		if err := globalRepo.DeleteIndex(partition1); err != nil {
			t.Fatalf("Failed to delete test index: %v", err)
		}
		os.RemoveAll(TestDir)
	}()

	// Define the expected result
	expectedPartitions := []string{"testPartition", "testPartition1"}

	// Run the test cases
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "Get all partitions from mock repository",
			want: expectedPartitions,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := globalRepo.GetPartitions()
			if len(got) != len(tt.want) {
				t.Errorf("Length mismatch: got %v, want %v", got, tt.want)
				return
			}
			for _, expected := range tt.want {
				found := false
				for _, actual := range got {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Partition %s not found in %v", expected, got)
				}
			}
		})
	}
}
