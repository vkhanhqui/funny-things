package integration_test

import (
	"bleve-proj/internal/controllers"
	"bleve-proj/internal/repositories"
	"bleve-proj/internal/routes"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// MockIndexRepository implements IndexRepository with mock behaviors.
type MockIndexRepository struct{}

func (m *MockIndexRepository) LoadIndexes() error {
	fmt.Println("Mock: Loading indexes")
	return nil
}

func (m *MockIndexRepository) CreateIndex(partition string, index map[string]interface{}) (bleve.Index, error) {
	fmt.Printf("Mock: Creating index for partition %s\n", partition)
	return nil, nil
}

func (m *MockIndexRepository) OpenIndex(partition string) (bleve.Index, error) {
	fmt.Printf("Mock: Opening index for partition %s\n", partition)
	return nil, nil
}

func (m *MockIndexRepository) DeleteIndex(partition string) error {
	fmt.Printf("Mock: Deleting index for partition %s\n", partition)
	return nil
}

func (m *MockIndexRepository) IndexDocument(partition string, docId string, document interface{}) error {
	fmt.Printf("Mock: Indexing document with ID %s in partition %s\n", docId, partition)
	return nil
}

func (m *MockIndexRepository) Search(partition string, query string) (*bleve.SearchResult, error) {
	fmt.Printf("Mock: Searching in partition %s with query %s\n", partition, query)
	return nil, nil
}

func (m *MockIndexRepository) GetDocumentCount(partition string) (uint64, error) {
	fmt.Printf("Mock: Getting document count in partition %s\n", partition)
	return 1, nil
}

func (m *MockIndexRepository) BulkLoadDocuments(partition string, scanner *bufio.Scanner) error {
	fmt.Printf("Mock: Bulk loading documents in partition %s\n", partition)
	return nil
}

func (m *MockIndexRepository) CheckIndexContents(partition string) (map[string]interface{}, error) {
	fmt.Printf("Mock: Checking index contents in partition %s\n", partition)
	return nil, nil
}

func (m *MockIndexRepository) PrintTermDictionary(partition, field string) (map[string]uint64, error) {
	fmt.Printf("Mock: Printing term dictionary for field %s in partition %s\n", field, partition)
	return nil, nil
}

func (m *MockIndexRepository) DumpIndexContentsWithLimit(partition string, offset, limit int) ([]interface{}, error) {
	fmt.Printf("Mock: Dumping index contents with limit in partition %s, offset %d, limit %d\n", partition, offset, limit)
	return nil, nil
}

func (m *MockIndexRepository) ListIndexFields(partition string) ([]string, error) {
	fmt.Printf("Mock: Listing index fields in partition %s\n", partition)
	return nil, nil
}

func (m *MockIndexRepository) PrintIndexMapping(partition string) (map[string]interface{}, error) {
	fmt.Printf("Mock: Printing index mapping in partition %s\n", partition)
	return nil, nil
}

func (m *MockIndexRepository) Indexes() map[string]bleve.Index {
	fmt.Println("Mock: Getting indexes")
	return nil
}

func (m *MockIndexRepository) GetPartitions() []string {
	fmt.Println("Mock: Getting partitions")
	return nil
}

func (m *MockIndexRepository) Close() {
	fmt.Println("Mock: Closing indexes")
}

func setupAppWithMockRepo() *fiber.App {
	// Create a new Fiber instance
	app := fiber.New()

	// Replace the real repository with the mock repository
	repositories.IndexRepo = &MockIndexRepository{}

	// Define routes
	routes.InitRoutes(app)

	return app
}

func TestCreateIndexEndpoint_POC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()
	indexMappingJSON := `
	{
		"types": {
			"document": {
				"properties": {
					"title": {
						"type": "text"
					},
					"author": {
						"type": "keyword"
					},
					"date": {
						"type": "datetime"
					}
				}
			}
		},
		"default_mapping": {
			"properties": {
				"defaultField": {
					"type": "text"
				}
			}
		},
		"type_field": "_type",
		"default_type": "_default",
		"default_analyzer": "standard",
		"default_datetime_parser": "dateTimeOptional",
		"default_field": "defaultField",
		"store_dynamic": true,
		"index_dynamic": true,
		"docvalues_dynamic": false
	}`
	// Create an HTTP request for CreateIndex endpoint
	req, err := http.NewRequest(http.MethodPost, "/api/v1/"+generateRandomString(9)+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create request")

	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make request")

	// Assert response status
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 201 for CreateIndex")

	// Read response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")
	defer resp.Body.Close()

	var response controllers.Response
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err, "Failed to unmarshal JSON response body")

	// Assert the response structure and content
	assert.Equal(t, http.StatusCreated, response.Code, "Expected response code to match")
	assert.True(t, response.Success, "Expected response success to be true")
	assert.Nil(t, response.Error, "Expected no error in response")
	assert.NotNil(t, response.Data, "Expected data in response")

	// Optionally, you can check specific fields in the data if needed
	// For example, if your success response includes a message:
	// assert.Equal(t, "Index created successfully", response.Data["message"])

	// Clean up: Perform any necessary cleanup steps after testing
	// For example, delete the created index or reset the test environment
}

func TestDeleteIndexEndpoint_POC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Create an HTTP request for DeleteIndex endpoint
	req, err := http.NewRequest(http.MethodDelete, "/api/v1/"+generateRandomString(9)+"/delete", nil)
	assert.NoError(t, err, "Failed to create request")

	// Perform the request
	resp, err := app.Test(req, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make request")

	// Assert response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 for DeleteIndex")
}

func TestIndexDocumentEndpoint_POC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate a random partition name
	partition := generateRandomString(9)

	// Create a test document
	document := map[string]interface{}{
		"id":   "123",
		"name": "Jane Doe",
		"age":  25,
	}
	documentJSON, _ := json.Marshal(document)

	// Create an HTTP request to index a document
	req, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
	assert.NoError(t, err, "Failed to create request for IndexDocument")

	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make IndexDocument request")

	// Assert response status
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 201 for IndexDocument")

	// Read response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")
	defer resp.Body.Close()

	var response controllers.Response
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err, "Failed to unmarshal JSON response body")

	// Assert the response structure and content
	assert.Equal(t, http.StatusCreated, response.Code, "Expected response code to match")
	assert.True(t, response.Success, "Expected response success to be true")
	assert.Nil(t, response.Error, "Expected no error in response")
	assert.NotNil(t, response.Data, "Expected data in response")

	// Optionally, check specific fields in the data if needed
	// For example, if your success response includes a message:
	// assert.Equal(t, "Document indexed successfully", response.Data["message"])

	// Clean up: Perform any necessary cleanup steps after testing
	// For example, delete the created index or reset the test environment
}

func TestSearchDocument_POC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate a random partition name
	partition := generateRandomString(9)

	// Create a test document
	document := map[string]interface{}{
		"id":   "123",
		"name": "Jane Doe",
		"age":  25,
	}
	documentJSON, _ := json.Marshal(document)

	// Create an HTTP request to index a document
	indexReq, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
	assert.NoError(t, err, "Failed to create request for IndexDocument")

	indexReq.Header.Set("Content-Type", "application/json")

	// Perform the IndexDocument request
	indexResp, err := app.Test(indexReq, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make IndexDocument request")

	// Assert IndexDocument response status and body
	assert.Equal(t, http.StatusCreated, indexResp.StatusCode, "Expected status 201 for IndexDocument")

	// Read IndexDocument response body
	bodyIndex, err := io.ReadAll(indexResp.Body)
	assert.NoError(t, err, "Failed to read IndexDocument response body")
	defer indexResp.Body.Close()

	// Assert IndexDocument response message
	var indexResponse map[string]interface{}
	err = json.Unmarshal(bodyIndex, &indexResponse)
	assert.NoError(t, err, "Failed to unmarshal IndexDocument response body")
	assert.Equal(t, true, indexResponse["success"].(bool), "Expected success field to be true")
	assert.Equal(t, "Document indexed successfully", indexResponse["data"].([]interface{})[0], "Expected success message in IndexDocument response")

	// Create an HTTP request to search for documents
	searchReq, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/search?q=Doe", nil)
	assert.NoError(t, err, "Failed to create request for SearchDocument")

	searchReq.Header.Set("Content-Type", "application/json")

	// Perform the SearchDocument request
	searchResp, err := app.Test(searchReq, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make SearchDocument request")

	// Assert SearchDocument response status and body
	assert.Equal(t, http.StatusOK, searchResp.StatusCode, "Expected status 200 for SearchDocument")

	// Read SearchDocument response body
	bodySearch, err := io.ReadAll(searchResp.Body)
	assert.NoError(t, err, "Failed to read SearchDocument response body")
	defer searchResp.Body.Close()

	// Assert response structure and content
	var searchResponse map[string]interface{}
	err = json.Unmarshal(bodySearch, &searchResponse)
	assert.NoError(t, err, "Failed to unmarshal SearchDocument response body")
	assert.Equal(t, 200, int(searchResponse["code"].(float64)), "Expected response code to match")
	assert.True(t, searchResponse["success"].(bool), "Expected success field to be true")
	assert.NotNil(t, searchResponse["data"], "Expected data field to be present")
	// Optionally, add further assertions on the content of searchResponse["data"]

}

func TestGetDocumentCountPOC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate a random partition name
	partition := generateRandomString(9)

	// Create a test document
	document := map[string]interface{}{
		"id":   "1",
		"name": "Alice",
	}
	documentJSON, _ := json.Marshal(document)

	// Create an HTTP request to index a document
	reqIndex, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
	assert.NoError(t, err, "Failed to create request for IndexDocument")

	reqIndex.Header.Set("Content-Type", "application/json")

	// Perform the IndexDocument request
	respIndex, err := app.Test(reqIndex, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make IndexDocument request")

	// Assert IndexDocument response status and body
	assert.Equal(t, http.StatusCreated, respIndex.StatusCode, "Expected status 201 for IndexDocument")

	// Read IndexDocument response body
	bodyIndex, err := io.ReadAll(respIndex.Body)
	assert.NoError(t, err, "Failed to read IndexDocument response body")
	defer respIndex.Body.Close()

	// Assert IndexDocument response message
	var indexResponse map[string]interface{}
	err = json.Unmarshal(bodyIndex, &indexResponse)
	assert.NoError(t, err, "Failed to unmarshal IndexDocument response body")
	assert.Contains(t, indexResponse, "data", "Expected 'data' field in IndexDocument response")

	data, ok := indexResponse["data"].([]interface{})
	assert.True(t, ok, "Expected 'data' to be an array")
	assert.Len(t, data, 1, "Expected 'data' array to have length 1")
	assert.Equal(t, "Document indexed successfully", data[0], "Expected success message in IndexDocument response")

	// Now, perform a request to get the document count
	reqCount, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/count", nil)
	assert.NoError(t, err, "Failed to create request for GetDocumentCount")

	// Perform the GetDocumentCount request
	respCount, err := app.Test(reqCount, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make GetDocumentCount request")

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respCount.StatusCode, "Expected status 200 for GetDocumentCount")

	// Read response body
	bodyCount, err := io.ReadAll(respCount.Body)
	assert.NoError(t, err, "Failed to read GetDocumentCount response body")
	defer respCount.Body.Close()

	// Assert the document count
	var countResponse map[string]interface{}
	err = json.Unmarshal(bodyCount, &countResponse)
	assert.NoError(t, err, "Failed to unmarshal JSON response body")
	assert.Contains(t, countResponse, "data", "Expected 'data' field in response")

	data, ok = countResponse["data"].([]interface{})
	assert.True(t, ok, "Expected 'data' to be an array")
	assert.Len(t, data, 1, "Expected 'data' array to have length 1")
	count, ok := data[0].(float64)
	assert.True(t, ok, "Expected count to be a float64")
	assert.Equal(t, float64(1), count, "Expected document count to be 1")
}
func TestBulkLoadDocumentsPOC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate a random partition name
	partition := generateRandomString(9)

	// Create a string with multiple JSON documents
	documents := `
		{"id": "1", "name": "Alice", "age": 30}
		{"id": "2", "name": "Bob", "age": 25}
	`

	// Create an HTTP request to bulk load documents
	req, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/bulk", bytes.NewBufferString(documents))
	assert.NoError(t, err, "Failed to create request for BulkLoadDocuments")

	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make BulkLoadDocuments request")

	// Assert response status
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 for BulkLoadDocuments")

	// Clean up: Delete the index after testing
	reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
	assert.NoError(t, err, "Failed to create request for DeleteIndex")

	// Perform the DeleteIndex request
	respDelete, err := app.Test(reqDelete, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make DeleteIndex request")

	// Assert DeleteIndex response status code
	assert.Equal(t, http.StatusOK, respDelete.StatusCode, "Expected status 200 for DeleteIndex")
}
func TestCheckIndexContentsPOC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate a random partition name
	partition := generateRandomString(9)

	// Create a test document
	document := map[string]interface{}{
		"id":   "123",
		"name": "John Doe",
		"age":  30,
	}
	documentJSON, _ := json.Marshal(document)

	// Create an HTTP request to index a document
	reqIndex, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
	assert.NoError(t, err, "Failed to create request for IndexDocument")

	reqIndex.Header.Set("Content-Type", "application/json")

	// Perform the IndexDocument request
	respIndex, err := app.Test(reqIndex, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make IndexDocument request")

	// Assert IndexDocument response status and body
	assert.Equal(t, http.StatusCreated, respIndex.StatusCode, "Expected status 201 for IndexDocument")

	// Read IndexDocument response body
	bodyIndex, err := io.ReadAll(respIndex.Body)
	assert.NoError(t, err, "Failed to read IndexDocument response body")
	defer respIndex.Body.Close()

	// Assert IndexDocument response structure
	var indexResponse map[string]interface{}
	err = json.Unmarshal(bodyIndex, &indexResponse)
	assert.NoError(t, err, "Failed to unmarshal IndexDocument response body")

	// Assert success and data content
	assert.Equal(t, true, indexResponse["success"], "Expected success to be true in IndexDocument response")
	assert.Contains(t, indexResponse["data"], "Document indexed successfully", "Expected data field to contain success message")

	// Create an HTTP request to check the index contents
	reqCheck, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/check", nil)
	assert.NoError(t, err, "Failed to create request for CheckIndexContents")

	// Perform the CheckIndexContents request
	respCheck, err := app.Test(reqCheck, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make CheckIndexContents request")

	// Assert CheckIndexContents response status and body
	assert.Equal(t, http.StatusOK, respCheck.StatusCode, "Expected status 200 for CheckIndexContents")

	// Read response body
	bodyCheck, err := io.ReadAll(respCheck.Body)
	assert.NoError(t, err, "Failed to read CheckIndexContents response body")
	defer respCheck.Body.Close()

	// Assert response structure and content
	var checkResponse map[string]interface{}
	err = json.Unmarshal(bodyCheck, &checkResponse)
	assert.NoError(t, err, "Failed to unmarshal CheckIndexContents response body")

	// Check if 'data' field is present and of type []interface{}
	data, ok := checkResponse["data"].([]interface{})
	assert.True(t, ok, "Expected 'data' field to be an array")
	assert.NotNil(t, data, "Expected 'data' field to be non-nil")

	// Example: Assert the length of 'data'
	assert.Equal(t, 1, len(data), "Expected 'data' array to contain 1 item")

	// Optionally, further assertions based on your application's response structure

	// Clean up: Delete the index after testing
	reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
	assert.NoError(t, err, "Failed to create request for DeleteIndex")

	// Perform the DeleteIndex request
	respDelete, err := app.Test(reqDelete, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make DeleteIndex request")

	// Assert DeleteIndex response status code
	assert.Equal(t, http.StatusOK, respDelete.StatusCode, "Expected status 200 for DeleteIndex")
}

func TestPrintTermDictionaryPOC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate a random partition name
	partition := generateRandomString(9)

	// Create a test document
	document := map[string]interface{}{
		"id":   "123",
		"name": "John Doe",
		"age":  30,
	}
	documentJSON, _ := json.Marshal(document)

	// Create an HTTP request to index a document
	reqIndex, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
	assert.NoError(t, err, "Failed to create request for IndexDocument")

	reqIndex.Header.Set("Content-Type", "application/json")

	// Perform the IndexDocument request
	respIndex, err := app.Test(reqIndex, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make IndexDocument request")

	// Assert IndexDocument response status
	assert.Equal(t, http.StatusCreated, respIndex.StatusCode, "Expected status 201 for IndexDocument")

	// Create an HTTP request to print the term dictionary for a specific field
	field := "name"
	reqPrint, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/dictionary?field="+field, nil)
	assert.NoError(t, err, "Failed to create request for PrintTermDictionary")

	// Perform the PrintTermDictionary request
	respPrint, err := app.Test(reqPrint, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make PrintTermDictionary request")

	// Assert PrintTermDictionary response status
	assert.Equal(t, http.StatusOK, respPrint.StatusCode, "Expected status 200 for PrintTermDictionary")

	// Read response body
	bodyPrint, err := io.ReadAll(respPrint.Body)
	assert.NoError(t, err, "Failed to read PrintTermDictionary response body")
	defer respPrint.Body.Close()

	// Assert response structure and content
	var printResponse map[string]interface{}
	err = json.Unmarshal(bodyPrint, &printResponse)
	assert.NoError(t, err, "Failed to unmarshal PrintTermDictionary response body")

	// Ensure 'data' field exists
	data, dataExists := printResponse["data"]
	assert.True(t, dataExists, "Expected 'data' field to exist")
	assert.Equal(t, data, []interface{}{nil})
}
func TestDumpIndexContentsWithLimitPOC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate a random partition name
	partition := generateRandomString(9)

	// Create test documents
	documents := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
		{"id": "3", "name": "Alice", "age": 35},
	}
	for _, doc := range documents {
		documentJSON, _ := json.Marshal(doc)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
		assert.NoError(t, err, "Failed to create request for IndexDocument")

		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := app.Test(req, 10) // Timeout in milliseconds
		assert.NoError(t, err, "Failed to make IndexDocument request")

		// Assert response status and body
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 201 for IndexDocument")

		// Read response body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, "Failed to read response body")
		defer resp.Body.Close()

		// Assert response message
		var responseBody map[string]interface{}
		err = json.Unmarshal(body, &responseBody)
		assert.NoError(t, err, "Failed to unmarshal JSON response body")
		expectedMessage := "Document indexed successfully"
		actualMessage, ok := responseBody["data"].([]interface{})[0].(string)
		assert.True(t, ok, "Expected message field in JSON response")
		assert.Equal(t, expectedMessage, actualMessage, "Expected success message in JSON response")
	}

	// Create an HTTP request to dump the index contents with limit and page
	page := 1
	size := 2
	reqDump, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/%s/dump?page=%d&size=%d", partition, page, size), nil)
	assert.NoError(t, err, "Failed to create request for DumpIndexContentsWithLimit")

	// Perform the DumpIndexContentsWithLimit request
	respDump, err := app.Test(reqDump, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make DumpIndexContentsWithLimit request")

	// Assert DumpIndexContentsWithLimit response status
	assert.Equal(t, http.StatusOK, respDump.StatusCode, "Expected status 200 for DumpIndexContentsWithLimit")

	// Read response body
	bodyDump, err := io.ReadAll(respDump.Body)
	assert.NoError(t, err, "Failed to read DumpIndexContentsWithLimit response body")
	defer respDump.Body.Close()

	// Assert response structure and content
	var dumpResponse map[string]interface{}
	err = json.Unmarshal(bodyDump, &dumpResponse)
	assert.NoError(t, err, "Failed to unmarshal DumpIndexContentsWithLimit response body")

}

func TestListIndexFieldsPOC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate a random partition name
	partition := generateRandomString(9)

	// Create test documents
	documents := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
	}
	for _, doc := range documents {
		documentJSON, _ := json.Marshal(doc)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
		assert.NoError(t, err, "Failed to create request for IndexDocument")

		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := app.Test(req, 10) // Timeout in milliseconds
		assert.NoError(t, err, "Failed to make IndexDocument request")

		// Assert response status and body
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 201 for IndexDocument")

		// Read response body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, "Failed to read response body")
		defer resp.Body.Close()

		// Assert response message
		var responseBody map[string]interface{}
		err = json.Unmarshal(body, &responseBody)
		assert.NoError(t, err, "Failed to unmarshal JSON response body")

		expectedMessage := "Document indexed successfully"
		actualMessage, ok := responseBody["data"].([]interface{})[0].(string)
		assert.True(t, ok, "Expected message field in JSON response")
		assert.Equal(t, expectedMessage, actualMessage, "Expected success message in JSON response")
	}

	// Create an HTTP request to list the index fields
	reqListFields, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/fields", nil)
	assert.NoError(t, err, "Failed to create request for ListIndexFields")

	// Perform the ListIndexFields request
	respListFields, err := app.Test(reqListFields, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make ListIndexFields request")

	// Assert ListIndexFields response status
	assert.Equal(t, http.StatusOK, respListFields.StatusCode, "Expected status 200 for ListIndexFields")

	// Read response body
	bodyListFields, err := io.ReadAll(respListFields.Body)
	assert.NoError(t, err, "Failed to read ListIndexFields response body")
	defer respListFields.Body.Close()

	// Unmarshal response body into a slice of strings
	var fieldsResponse map[string]any
	err = json.Unmarshal(bodyListFields, &fieldsResponse)
	assert.NoError(t, err, "Failed to unmarshal ListIndexFields response body")
}

func TestPrintIndexMappingPOC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate a random partition name
	partition := generateRandomString(9)

	// Create test documents
	documents := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
	}
	for _, doc := range documents {
		documentJSON, _ := json.Marshal(doc)
		req, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
		assert.NoError(t, err, "Failed to create request for IndexDocument")

		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := app.Test(req, 10) // Timeout in milliseconds
		assert.NoError(t, err, "Failed to make IndexDocument request")

		// Assert response status and body
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status 201 for IndexDocument")

		// Read response body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, "Failed to read response body")
		defer resp.Body.Close()

		// Assert response message
		var responseBody map[string]interface{}
		err = json.Unmarshal(body, &responseBody)
		assert.NoError(t, err, "Failed to unmarshal JSON response body")

		expectedMessage := "Document indexed successfully"
		actualMessage, ok := responseBody["data"].([]interface{})[0].(string)
		assert.True(t, ok, "Expected message field in JSON response")
		assert.Equal(t, expectedMessage, actualMessage, "Expected success message in JSON response")
	}

	// Create an HTTP request to print the index mapping
	reqPrintMapping, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/mapping", nil)
	assert.NoError(t, err, "Failed to create request for PrintIndexMapping")

	// Perform the PrintIndexMapping request
	respPrintMapping, err := app.Test(reqPrintMapping, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make PrintIndexMapping request")

	// Assert PrintIndexMapping response status
	assert.Equal(t, http.StatusOK, respPrintMapping.StatusCode, "Expected status 200 for PrintIndexMapping")

	// Read response body
	bodyPrintMapping, err := io.ReadAll(respPrintMapping.Body)
	assert.NoError(t, err, "Failed to read PrintIndexMapping response body")
	defer respPrintMapping.Body.Close()

	// Assert response structure and content
	var mappingResponse map[string]interface{}
	err = json.Unmarshal(bodyPrintMapping, &mappingResponse)
	assert.NoError(t, err, "Failed to unmarshal PrintIndexMapping response body")

}

func TestGetPartitionsPOC(t *testing.T) {
	// Initialize the Fiber app with mock repository
	app := setupAppWithMockRepo()

	// Generate random partition names
	partition1 := generateRandomString(9)
	partition2 := generateRandomString(9)

	// Create index mappings for the test partitions
	indexMappingJSON := `
	{
		"types": {
			"document": {
				"properties": {
					"title": {
						"type": "text"
					},
					"author": {
						"type": "keyword"
					},
					"date": {
						"type": "datetime"
					}
				}
			}
		},
		"default_mapping": {
			"properties": {
				"defaultField": {
					"type": "text"
				}
			}
		},
		"type_field": "_type",
		"default_type": "_default",
		"default_analyzer": "standard",
		"default_datetime_parser": "dateTimeOptional",
		"default_field": "defaultField",
		"store_dynamic": true,
		"index_dynamic": true,
		"docvalues_dynamic": false
	}`

	// Create the first index
	reqCreate1, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition1+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create request for partition 1")
	reqCreate1.Header.Set("Content-Type", "application/json")
	respCreate1, err := app.Test(reqCreate1, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make CreateIndex request for partition 1")
	assert.Equal(t, http.StatusCreated, respCreate1.StatusCode, "Expected status 201 for CreateIndex partition 1")

	// Create the second index
	reqCreate2, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition2+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create request for partition 2")
	reqCreate2.Header.Set("Content-Type", "application/json")
	respCreate2, err := app.Test(reqCreate2, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make CreateIndex request for partition 2")
	assert.Equal(t, http.StatusCreated, respCreate2.StatusCode, "Expected status 201 for CreateIndex partition 2")

	// Create an HTTP request to get the partitions
	reqGetPartitions, err := http.NewRequest(http.MethodGet, "/api/v1/", nil)
	assert.NoError(t, err, "Failed to create request for GetPartitions")

	// Perform the GetPartitions request
	respGetPartitions, err := app.Test(reqGetPartitions, 10) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make GetPartitions request")

	// Assert GetPartitions response status
	assert.Equal(t, http.StatusOK, respGetPartitions.StatusCode, "Expected status 200 for GetPartitions")

	// Read response body
	bodyGetPartitions, err := io.ReadAll(respGetPartitions.Body)
	assert.NoError(t, err, "Failed to read GetPartitions response body")
	defer respGetPartitions.Body.Close()

	// Assert response structure and content
	var getPartitionsResponse map[string]interface{}
	err = json.Unmarshal(bodyGetPartitions, &getPartitionsResponse)
	assert.NoError(t, err, "Failed to unmarshal GetPartitions response body")

	// Assert GetPartitions response success
	assert.True(t, getPartitionsResponse["success"].(bool), "Expected success to be true in GetPartitions response")

	// // Assert GetPartitions response data
	// data := getPartitionsResponse["data"].([]interface{})
	// assert.Contains(t, data, partition1, "Expected partition 1 in GetPartitions response data")
	// assert.Contains(t, data, partition2, "Expected partition 2 in GetPartitions response data")

	// // Clean up: Delete the test partitions
	// reqDelete1, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition1+"/delete", nil)
	// assert.NoError(t, err, "Failed to create request for DeleteIndex partition 1")
	// respDelete1, err := app.Test(reqDelete1, 10) // Timeout in milliseconds
	// assert.NoError(t, err, "Failed to make DeleteIndex request for partition 1")
	// assert.Equal(t, http.StatusOK, respDelete1.StatusCode, "Expected status 200 for DeleteIndex partition 1")

	// reqDelete2, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition2+"/delete", nil)
	// assert.NoError(t, err, "Failed to create request for DeleteIndex partition 2")
	// respDelete2, err := app.Test(reqDelete2, 10) // Timeout in milliseconds
	// assert.NoError(t, err, "Failed to make DeleteIndex request for partition 2")
	// assert.Equal(t, http.StatusOK, respDelete2.StatusCode, "Expected status 200 for DeleteIndex partition 2")
}
