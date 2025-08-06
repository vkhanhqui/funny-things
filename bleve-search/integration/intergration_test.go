package integration_test

import (
	"bleve-proj/internal/repositories"
	"bleve-proj/internal/routes"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"

	// Adjust the import path as necessary

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var TestIntegrationDir = "integration_tests"

// generateRandomString generates a random alphanumeric string of given length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func setupApp(dir string) *fiber.App {
	// Initialize Fiber app
	app := fiber.New()

	// Initialize index repository with index directory
	repositories.IndexRepo = repositories.NewIndexRepository(TestIntegrationDir + "/" + dir)

	// Define routes
	routes.InitRoutes(app)

	return app
}

func TestCreateIndexEndpoint(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)

	// Initialize the Fiber app
	app := setupApp(partition)

	// Define your basic index mapping JSON
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

	// Create an HTTP request for CreateIndex endpoint with the index mapping JSON payload
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create request")

	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // No timeout
	assert.NoError(t, err, "Failed to make CreateIndex request")
	defer respCreate.Body.Close()

	// Assert CreateIndex response status code
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode, "Expected status 201 for CreateIndex")

	// Read CreateIndex response body
	bodyCreate, err := io.ReadAll(respCreate.Body)
	assert.NoError(t, err, "Failed to read CreateIndex response body")
	fmt.Printf("string(bodyCreate): %v\n", string(bodyCreate))
	// Unmarshal CreateIndex response body
	var createIndexResponse map[string]interface{}
	err = json.Unmarshal(bodyCreate, &createIndexResponse)
	assert.NoError(t, err, "Failed to unmarshal CreateIndex response body")

	// Assert CreateIndex response success
	assert.True(t, createIndexResponse["success"].(bool), "Expected success to be true in CreateIndex response")

	// Assert CreateIndex response message
	data := createIndexResponse["data"].([]interface{})
	assert.Len(t, data, 1, "Expected data to have exactly one element")
	assert.Equal(t, "Index created successfully", data[0], "Expected success message in CreateIndex response")

	// Clean up: Delete the index after testing
	reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
	assert.NoError(t, err, "Failed to create request for DeleteIndex")

	// Perform the DeleteIndex request
	respDelete, err := app.Test(reqDelete, -1) // No timeout
	assert.NoError(t, err, "Failed to make DeleteIndex request")

	// Assert DeleteIndex response status code
	assert.Equal(t, http.StatusOK, respDelete.StatusCode, "Expected status 200 for DeleteIndex")
}

func TestDeleteIndexEndpoint(t *testing.T) {
	// Initialize the Fiber app
	partition := generateRandomString(9)
	app := setupApp(partition)

	// Create an index for testing
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

	// Create an HTTP request for CreateIndex endpoint with the index mapping JSON payload
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create index request")
	reqCreate.Header.Set("Content-Type", "application/json")
	respCreate, err := app.Test(reqCreate, -1) // No timeout
	assert.NoError(t, err, "Failed to make create index request")
	defer respCreate.Body.Close()
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode, "Expected status 201 for CreateIndex")

	// Create an HTTP request for DeleteIndex endpoint
	reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
	assert.NoError(t, err, "Failed to create delete index request")

	// Perform the DeleteIndex request
	respDelete, err := app.Test(reqDelete, -1) // No timeout
	assert.NoError(t, err, "Failed to make delete index request")
	defer respDelete.Body.Close()

	// Assert DeleteIndex response status code
	assert.Equal(t, http.StatusOK, respDelete.StatusCode, "Expected status 200 for DeleteIndex")

	// Read DeleteIndex response body
	bodyDelete, err := io.ReadAll(respDelete.Body)
	assert.NoError(t, err, "Failed to read DeleteIndex response body")

	// Unmarshal DeleteIndex response body
	var deleteIndexResponse map[string]interface{}
	err = json.Unmarshal(bodyDelete, &deleteIndexResponse)
	assert.NoError(t, err, "Failed to unmarshal DeleteIndex response body")

	// Assert DeleteIndex response success
	assert.True(t, deleteIndexResponse["success"].(bool), "Expected success to be true in DeleteIndex response")

	// Assert DeleteIndex response message (optional if data is not returned)
	// Example: You might not have data returned in DeleteIndex success response

	// Clean up: Ensure the index is deleted after testing
	// Optionally, you can check for an error response if the index was not found
	// assert.Equal(t, http.StatusNotFound, respDelete.StatusCode, "Expected status 404 if index was not found")
}

func TestIndexDocumentEndpointWithIndexCreation(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)

	// Initialize the Fiber app
	app := setupApp(partition)

	// Create an HTTP request to create an index for the partition
	createIndexReq := map[string]string{"partition": partition}
	createIndexReqJSON, _ := json.Marshal(createIndexReq)
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", bytes.NewBuffer(createIndexReqJSON))
	assert.NoError(t, err, "Failed to create request for CreateIndex")
	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // No timeout
	assert.NoError(t, err, "Failed to make CreateIndex request")
	defer respCreate.Body.Close()

	// Assert CreateIndex response status and body
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode, "Expected status 201 for CreateIndex")

	// Read CreateIndex response body
	var createIndexResponse map[string]interface{}
	err = json.NewDecoder(respCreate.Body).Decode(&createIndexResponse)
	assert.NoError(t, err, "Failed to decode CreateIndex response body")
	assert.True(t, createIndexResponse["success"].(bool), "Expected success to be true in CreateIndex response")

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
	respIndex, err := app.Test(reqIndex, -1) // No timeout
	assert.NoError(t, err, "Failed to make IndexDocument request")
	defer respIndex.Body.Close()

	// Assert IndexDocument response status and body
	assert.Equal(t, http.StatusCreated, respIndex.StatusCode, "Expected status 201 for IndexDocument")

	// Read IndexDocument response body
	var indexResponse map[string]interface{}
	err = json.NewDecoder(respIndex.Body).Decode(&indexResponse)
	assert.NoError(t, err, "Failed to decode IndexDocument response body")
	assert.True(t, indexResponse["success"].(bool), "Expected success to be true in IndexDocument response")

	// Clean up: Delete the index after testing
	reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
	assert.NoError(t, err, "Failed to create request for DeleteIndex")

	// Perform the DeleteIndex request
	respDelete, err := app.Test(reqDelete, -1) // No timeout
	assert.NoError(t, err, "Failed to make DeleteIndex request")
	defer respDelete.Body.Close()

	// Assert DeleteIndex response status code
	assert.Equal(t, http.StatusOK, respDelete.StatusCode, "Expected status 200 for DeleteIndex")
}
func TestSearchDocumentEndpointWithIndexCreation(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)

	// Initialize the Fiber app
	app := setupApp(partition)

	// Create an HTTP request to create an index for the partition
	createIndexReq := map[string]string{"partition": partition}
	createIndexReqJSON, _ := json.Marshal(createIndexReq)
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", bytes.NewBuffer(createIndexReqJSON))
	assert.NoError(t, err, "Failed to create request for CreateIndex")
	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // No timeout
	assert.NoError(t, err, "Failed to make CreateIndex request")
	defer respCreate.Body.Close()

	// Assert CreateIndex response status and body
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode, "Expected status 201 for CreateIndex")

	// Read CreateIndex response body
	var createIndexResponse map[string]interface{}
	err = json.NewDecoder(respCreate.Body).Decode(&createIndexResponse)
	assert.NoError(t, err, "Failed to decode CreateIndex response body")
	assert.True(t, createIndexResponse["success"].(bool), "Expected success to be true in CreateIndex response")

	// Index a test document for searching
	document := map[string]interface{}{
		"id":   "123",
		"name": "john doe",
		"age":  30,
	}
	documentJSON, _ := json.Marshal(document)

	// Create an HTTP request to index a document
	reqIndex, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
	assert.NoError(t, err, "Failed to create request for IndexDocument")
	reqIndex.Header.Set("Content-Type", "application/json")

	// Perform the IndexDocument request
	respIndex, err := app.Test(reqIndex, -1) // No timeout
	assert.NoError(t, err, "Failed to make IndexDocument request")
	defer respIndex.Body.Close()

	// Assert IndexDocument response status and body
	assert.Equal(t, http.StatusCreated, respIndex.StatusCode, "Expected status 201 for IndexDocument")

	// Read IndexDocument response body
	var indexResponse map[string]interface{}
	err = json.NewDecoder(respIndex.Body).Decode(&indexResponse)
	assert.NoError(t, err, "Failed to decode IndexDocument response body")
	assert.True(t, indexResponse["success"].(bool), "Expected success to be true in IndexDocument response")

	// Perform a search request
	reqSearch, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/search?q=[{\"query\": {\"type\": \"term\",\"term\": \"john\"}}]", nil)
	assert.NoError(t, err, "Failed to create request for SearchDocument")
	reqSearch.Header.Set("Content-Type", "application/json")

	// Perform the SearchDocument request
	respSearch, err := app.Test(reqSearch, -1) // No timeout
	assert.NoError(t, err, "Failed to make SearchDocument request")
	defer respSearch.Body.Close()

	// Assert SearchDocument response status and body
	assert.Equal(t, http.StatusOK, respSearch.StatusCode, "Expected status 200 for SearchDocument")

	// Read SearchDocument response body
	var searchResponse map[string]interface{}
	err = json.NewDecoder(respSearch.Body).Decode(&searchResponse)
	assert.NoError(t, err, "Failed to decode SearchDocument response body")
	// Assert the presence of 'hits' field in searchResponse
	hits, ok := searchResponse["data"].([]interface{})
	assert.True(t, ok, "Expected 'data' field to be an array")
	assert.GreaterOrEqual(t, len(hits), 1, "Expected at least one hit in the search results")

	// Clean up: Delete the index after testing
	reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
	assert.NoError(t, err, "Failed to create request for DeleteIndex")

	// Perform the DeleteIndex request
	respDelete, err := app.Test(reqDelete, -1) // No timeout
	assert.NoError(t, err, "Failed to make DeleteIndex request")
	defer respDelete.Body.Close()

	// Assert DeleteIndex response status code
	assert.Equal(t, http.StatusOK, respDelete.StatusCode, "Expected status 200 for DeleteIndex")
}

func TestGetDocumentCountIntegration(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)

	// Initialize the Fiber app
	app := setupApp(partition)

	// Create an HTTP request to create an index for the partition
	createIndexReq := map[string]string{"partition": partition}
	createIndexReqJSON, _ := json.Marshal(createIndexReq)
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", bytes.NewBuffer(createIndexReqJSON))
	assert.NoError(t, err, "Failed to create request for CreateIndex")
	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // No timeout
	assert.NoError(t, err, "Failed to make CreateIndex request")
	defer respCreate.Body.Close()

	// Assert CreateIndex response status and body
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode, "Expected status 201 for CreateIndex")

	// Read CreateIndex response body
	var createIndexResponse map[string]interface{}
	err = json.NewDecoder(respCreate.Body).Decode(&createIndexResponse)
	assert.NoError(t, err, "Failed to decode CreateIndex response body")
	assert.True(t, createIndexResponse["success"].(bool), "Expected success to be true in CreateIndex response")

	// Index a document for the partition
	indexDocument := map[string]interface{}{
		"id":   "1",
		"name": "Alice",
	}
	indexDocumentJSON, _ := json.Marshal(indexDocument)
	reqIndex, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(indexDocumentJSON))
	assert.NoError(t, err, "Failed to create request for IndexDocument")
	reqIndex.Header.Set("Content-Type", "application/json")

	// Perform the IndexDocument request
	respIndex, err := app.Test(reqIndex, -1) // No timeout
	assert.NoError(t, err, "Failed to make IndexDocument request")
	defer respIndex.Body.Close()

	// Assert IndexDocument response status and body
	assert.Equal(t, http.StatusCreated, respIndex.StatusCode, "Expected status 201 for IndexDocument")

	// Read IndexDocument response body
	var indexResponse map[string]interface{}
	err = json.NewDecoder(respIndex.Body).Decode(&indexResponse)
	assert.NoError(t, err, "Failed to decode IndexDocument response body")
	assert.True(t, indexResponse["success"].(bool), "Expected success to be true in IndexDocument response")

	// Now, perform a request to get the document count
	reqCount, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/count", nil)
	assert.NoError(t, err, "Failed to create request for GetDocumentCount")

	// Perform the GetDocumentCount request
	respCount, err := app.Test(reqCount, -1) // No timeout
	assert.NoError(t, err, "Failed to make GetDocumentCount request")
	defer respCount.Body.Close()

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respCount.StatusCode, "Expected status 200 for GetDocumentCount")

	// Read response body
	bodyCount, err := io.ReadAll(respCount.Body)
	assert.NoError(t, err, "Failed to read GetDocumentCount response body")

	// Assert the document count
	var countResponse map[string]interface{}
	err = json.Unmarshal(bodyCount, &countResponse)
	assert.NoError(t, err, "Failed to unmarshal JSON response body")
	assert.Contains(t, countResponse, "data", "Expected 'data' field in response")

	data, ok := countResponse["data"].([]interface{})
	assert.True(t, ok, "Expected 'data' to be a slice")

	if len(data) > 0 {
		countValue, ok := data[0].(float64)
		assert.True(t, ok, "Expected count value to be a float64")
		assert.Equal(t, float64(1), countValue, "Expected document count to be 1")
	} else {
		t.Error("Expected non-empty 'data' array in response")
	}
}

func TestBulkLoadDocumentsEndpoint(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)
	// Initialize the Fiber app
	app := setupApp(partition)

	// Create an HTTP request to create an index for the partition
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

	// Create an HTTP request for CreateIndex endpoint with the index mapping JSON payload
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create request for CreateIndex")

	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make CreateIndex request")

	// Assert CreateIndex response status and body
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode, "Expected status 201 for CreateIndex")

	// Read CreateIndex response body
	bodyCreate, err := io.ReadAll(respCreate.Body)
	assert.NoError(t, err, "Failed to read CreateIndex response body")
	defer respCreate.Body.Close()

	// Assert CreateIndex response message
	var createIndexResponse map[string]interface{}
	err = json.Unmarshal(bodyCreate, &createIndexResponse)
	assert.NoError(t, err, "Failed to unmarshal CreateIndex response body")
	assert.True(t, createIndexResponse["success"].(bool), "Expected success to be true in CreateIndex response")

	// Create a string with multiple JSON documents
	documents := `
		{"id": "1", "name": "Alice", "age": 30}
		{"id": "2", "name": "Bob", "age": 25}
	`

	// Create an HTTP request to bulk load documents
	reqBulkLoad, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/bulk", bytes.NewBufferString(documents))
	assert.NoError(t, err, "Failed to create request for BulkLoadDocuments")

	reqBulkLoad.Header.Set("Content-Type", "application/json")

	// Perform the BulkLoadDocuments request
	respBulkLoad, err := app.Test(reqBulkLoad, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make BulkLoadDocuments request")

	// Assert BulkLoadDocuments response status and body
	assert.Equal(t, http.StatusOK, respBulkLoad.StatusCode, "Expected status 200 for BulkLoadDocuments")

	// Read BulkLoadDocuments response body
	bodyBulkLoad, err := io.ReadAll(respBulkLoad.Body)
	assert.NoError(t, err, "Failed to read BulkLoadDocuments response body")
	defer respBulkLoad.Body.Close()

	// Assert BulkLoadDocuments response message
	var bulkLoadResponse map[string]interface{}
	err = json.Unmarshal(bodyBulkLoad, &bulkLoadResponse)
	assert.NoError(t, err, "Failed to unmarshal BulkLoadDocuments response body")
	assert.True(t, bulkLoadResponse["success"].(bool), "Expected success to be true in BulkLoadDocuments response")

	// Clean up: Delete the index after testing
	reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
	assert.NoError(t, err, "Failed to create request for DeleteIndex")

	// Perform the DeleteIndex request
	respDelete, err := app.Test(reqDelete, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make DeleteIndex request")

	// Assert DeleteIndex response status code
	assert.Equal(t, http.StatusOK, respDelete.StatusCode, "Expected status 200 for DeleteIndex")
}
func TestCheckIndexContentsEndpoint(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)

	// Initialize the Fiber app
	app := setupApp(partition)

	// Create an HTTP request to create an index for the partition
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

	// Create an HTTP request for CreateIndex endpoint with the index mapping JSON payload
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create request for CreateIndex")
	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make CreateIndex request")
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode, "Expected status 201 for CreateIndex")

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
	respIndex, err := app.Test(reqIndex, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make IndexDocument request")
	assert.Equal(t, http.StatusCreated, respIndex.StatusCode, "Expected status 201 for IndexDocument")

	// Create an HTTP request to check the index contents
	reqCheck, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/check", nil)
	assert.NoError(t, err, "Failed to create request for CheckIndexContents")

	// Perform the CheckIndexContents request
	respCheck, err := app.Test(reqCheck, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make CheckIndexContents request")
	assert.Equal(t, http.StatusOK, respCheck.StatusCode, "Expected status 200 for CheckIndexContents")

	// Read response body
	bodyCheck, err := io.ReadAll(respCheck.Body)
	assert.NoError(t, err, "Failed to read CheckIndexContents response body")
	defer respCheck.Body.Close()

	// Assert response message
	var checkResponse map[string]interface{}
	err = json.Unmarshal(bodyCheck, &checkResponse)
	assert.NoError(t, err, "Failed to unmarshal CheckIndexContents response body")

	// Check if "data" field exists and is of type []interface{}
	data, ok := checkResponse["data"].([]interface{})
	assert.True(t, ok, "Expected 'data' to be an array of maps")

	// Check if the first item in "data" contains the "name" field
	if len(data) > 0 {
		firstItem, ok := data[0].(map[string]interface{})
		assert.True(t, ok, "Expected first item in 'data' to be a map")
		assert.Contains(t, firstItem, "name", "Expected 'name' field in check result")
	} else {
		t.Errorf("Expected 'data' to contain at least one item")
	}

	// Clean up: Delete the index after testing
	reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
	assert.NoError(t, err, "Failed to create request for DeleteIndex")

	// Perform the DeleteIndex request
	respDelete, err := app.Test(reqDelete, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make DeleteIndex request")
	assert.Equal(t, http.StatusOK, respDelete.StatusCode, "Expected status 200 for DeleteIndex")
}

func TestPrintTermDictionaryEndpoint(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)

	// Initialize the Fiber app
	app := setupApp(partition)

	// Create an HTTP request to create an index for the partition
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

	// Create an HTTP request for CreateIndex endpoint with the index mapping JSON payload
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create request for CreateIndex")
	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make CreateIndex request")
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode, "Expected status 201 for CreateIndex")

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
	respIndex, err := app.Test(reqIndex, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make IndexDocument request")
	assert.Equal(t, http.StatusCreated, respIndex.StatusCode, "Expected status 201 for IndexDocument")

	// Create an HTTP request to print the term dictionary for the field "name"
	field := "name"
	reqPrint, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/dictionary?field="+field, nil)
	assert.NoError(t, err, "Failed to create request for PrintTermDictionary")

	// Perform the PrintTermDictionary request
	respPrint, err := app.Test(reqPrint, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make PrintTermDictionary request")
	assert.Equal(t, http.StatusOK, respPrint.StatusCode, "Expected status 200 for PrintTermDictionary")

	// Read response body
	bodyPrint, err := io.ReadAll(respPrint.Body)
	assert.NoError(t, err, "Failed to read PrintTermDictionary response body")
	defer respPrint.Body.Close()

	// Assert response message
	var printResponse map[string]interface{}
	err = json.Unmarshal(bodyPrint, &printResponse)
	assert.NoError(t, err, "Failed to unmarshal PrintTermDictionary response body")

	// Ensure 'data' field exists
	data, found := printResponse["data"]
	assert.True(t, found, "Expected 'data' field in PrintTermDictionary response")

	// Assert that 'data' is an array and contains maps with expected term frequencies
	dataArray, ok := data.([]interface{})
	assert.True(t, ok, "Expected 'data' to be an array of maps")

	// Iterate over each map in 'data' array and check for expected terms and frequencies
	foundJohn := false
	for _, item := range dataArray {
		if termMap, ok := item.(map[string]interface{}); ok {
			for term, count := range termMap {
				if term == "john" {
					foundJohn = true
					assert.Equal(t, float64(1), count, "Expected term frequency of 'john' to be 1")
				}
			}
		}
	}

	assert.True(t, foundJohn, "Expected 'john' term in term dictionary for field 'name'")

	// Clean up: Delete the index after testing
	reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
	assert.NoError(t, err, "Failed to create request for DeleteIndex")

	// Perform the DeleteIndex request
	respDelete, err := app.Test(reqDelete, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make DeleteIndex request")
	assert.Equal(t, http.StatusOK, respDelete.StatusCode, "Expected status 200 for DeleteIndex")
}

func TestDumpIndexContentsWithLimitEndpoint(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)

	// Initialize the Fiber app
	app := setupApp(partition)
	if app == nil {
		t.Fatal("Failed to initialize Fiber app")
	}

	// Defer cleanup function
	defer func() {
		// Clean up: Delete the index after testing
		reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
		if err != nil {
			t.Fatalf("Failed to create request for DeleteIndex: %v", err)
		}

		// Perform the DeleteIndex request
		respDelete, err := app.Test(reqDelete, -1) // Timeout in milliseconds
		if err != nil {
			t.Fatalf("Failed to make DeleteIndex request: %v", err)
		}

		// Assert DeleteIndex response status code
		if respDelete.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for DeleteIndex, got %d", respDelete.StatusCode)
		}
	}()

	// Create an HTTP request to create an index for the partition
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

	// Create an HTTP request for CreateIndex endpoint with the index mapping JSON payload
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", strings.NewReader(indexMappingJSON))
	if err != nil {
		t.Fatalf("Failed to create request for CreateIndex: %v", err)
	}
	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // Timeout in milliseconds
	if err != nil {
		t.Fatalf("Failed to make CreateIndex request: %v", err)
	}
	if respCreate.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 for CreateIndex, got %d", respCreate.StatusCode)
	}

	// Create test documents
	documents := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
		{"id": "3", "name": "Alice", "age": 35},
	}
	for _, doc := range documents {
		documentJSON, _ := json.Marshal(doc)
		reqIndex, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
		if err != nil {
			t.Fatalf("Failed to create request for IndexDocument: %v", err)
		}
		reqIndex.Header.Set("Content-Type", "application/json")

		// Perform the IndexDocument request
		respIndex, err := app.Test(reqIndex, -1) // Timeout in milliseconds
		if err != nil {
			t.Fatalf("Failed to make IndexDocument request: %v", err)
		}
		if respIndex.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status 201 for IndexDocument, got %d", respIndex.StatusCode)
		}
	}

	// Create an HTTP request to dump the index contents with limit and page
	page := 1
	size := 2
	reqDump, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/%s/dump?page=%d&size=%d", partition, page, size), nil)
	if err != nil {
		t.Fatalf("Failed to create request for DumpIndexContentsWithLimit: %v", err)
	}

	// Perform the DumpIndexContentsWithLimit request
	respDump, err := app.Test(reqDump, -1) // Timeout in milliseconds
	if err != nil {
		t.Fatalf("Failed to make DumpIndexContentsWithLimit request: %v", err)
	}
	if respDump.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for DumpIndexContentsWithLimit, got %d", respDump.StatusCode)
	}

	// Read response body
	var dumpResponse map[string]interface{}
	err = json.NewDecoder(respDump.Body).Decode(&dumpResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal DumpIndexContentsWithLimit response body: %v", err)
	}

	// Check if the response contains expected keys or values
	if _, ok := dumpResponse["data"]; !ok {
		t.Fatalf("Expected 'data' key in DumpIndexContentsWithLimit response")
	}

	// Ensure the number of returned documents matches the requested size
	documentsList, ok := dumpResponse["data"].([]interface{})
	if !ok {
		t.Fatal("Expected 'data' field to be a list")
	}

	if len(documentsList) != size {
		t.Fatalf("Expected number of documents to match limit %d, got %d", size, len(documentsList))
	}

	// Ensure each document contains expected fields (example assertion)
	for _, doc := range documentsList {
		docMap, ok := doc.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected document to be a map[string]interface{}")
		}

		if _, ok := docMap["id"]; !ok {
			t.Fatalf("Expected 'id' field in dumped document")
		}
		if _, ok := docMap["name"]; !ok {
			t.Fatalf("Expected 'name' field in dumped document")
		}
		if _, ok := docMap["age"]; !ok {
			t.Fatalf("Expected 'age' field in dumped document")
		}
	}
}

func TestListIndexFieldsEndpoint(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)

	// Initialize the Fiber app
	app := setupApp(partition)
	if app == nil {
		t.Fatal("Failed to initialize Fiber app")
	}

	// Defer cleanup function
	defer func() {
		// Clean up: Delete the index after testing
		reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
		if err != nil {
			t.Fatalf("Failed to create request for DeleteIndex: %v", err)
		}

		// Perform the DeleteIndex request
		respDelete, err := app.Test(reqDelete, -1) // Timeout in milliseconds
		if err != nil {
			t.Fatalf("Failed to make DeleteIndex request: %v", err)
		}

		// Assert DeleteIndex response status code
		if respDelete.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for DeleteIndex, got %d", respDelete.StatusCode)
		}
	}()

	// Create an HTTP request to create an index for the partition
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

	// Create an HTTP request for CreateIndex endpoint with the index mapping JSON payload
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", strings.NewReader(indexMappingJSON))
	if err != nil {
		t.Fatalf("Failed to create request for CreateIndex: %v", err)
	}
	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // Timeout in milliseconds
	if err != nil {
		t.Fatalf("Failed to make CreateIndex request: %v", err)
	}
	if respCreate.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 for CreateIndex, got %d", respCreate.StatusCode)
	}

	// Create test documents
	documents := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
	}
	for _, doc := range documents {
		documentJSON, _ := json.Marshal(doc)
		reqIndex, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
		if err != nil {
			t.Fatalf("Failed to create request for IndexDocument: %v", err)
		}
		reqIndex.Header.Set("Content-Type", "application/json")

		// Perform the IndexDocument request
		respIndex, err := app.Test(reqIndex, -1) // Timeout in milliseconds
		if err != nil {
			t.Fatalf("Failed to make IndexDocument request: %v", err)
		}
		if respIndex.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status 201 for IndexDocument, got %d", respIndex.StatusCode)
		}
	}

	// Create an HTTP request to list the index fields
	reqListFields, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/fields", nil)
	if err != nil {
		t.Fatalf("Failed to create request for ListIndexFields: %v", err)
	}

	// Perform the ListIndexFields request
	respListFields, err := app.Test(reqListFields, -1) // Timeout in milliseconds
	if err != nil {
		t.Fatalf("Failed to make ListIndexFields request: %v", err)
	}
	if respListFields.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for ListIndexFields, got %d", respListFields.StatusCode)
	}

	// Read response body
	var fieldsResponse struct {
		Code    int      `json:"code"`
		Success bool     `json:"success"`
		Data    []string `json:"data"`
	}
	bodyListFields, err := io.ReadAll(respListFields.Body)
	if err != nil {
		t.Fatalf("Failed to read ListIndexFields response body: %v", err)
	}
	defer respListFields.Body.Close()

	// Unmarshal response body into struct
	err = json.Unmarshal(bodyListFields, &fieldsResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal ListIndexFields response body: %v", err)
	}

	// Assert fieldsResponse contains expected fields
	expectedFields := []string{"id", "name", "age", "_id", "_all"}
	assert.ElementsMatch(t, expectedFields, fieldsResponse.Data, "Expected fields to match indexed document fields")
}
func TestPrintIndexMappingEndpoint(t *testing.T) {
	// Generate a random partition name
	partition := generateRandomString(9)

	// Initialize the Fiber app
	app := setupApp(partition)
	if app == nil {
		t.Fatal("Failed to initialize Fiber app")
	}

	// Defer cleanup function
	defer func() {
		// Clean up: Delete the index after testing
		reqDelete, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition+"/delete", nil)
		if err != nil {
			t.Fatalf("Failed to create request for DeleteIndex: %v", err)
		}

		// Perform the DeleteIndex request
		respDelete, err := app.Test(reqDelete, -1) // Timeout in milliseconds
		if err != nil {
			t.Fatalf("Failed to make DeleteIndex request: %v", err)
		}

		// Assert DeleteIndex response status code
		if respDelete.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for DeleteIndex, got %d", respDelete.StatusCode)
		}

		os.RemoveAll(TestIntegrationDir)
	}()

	// Create an HTTP request to create an index for the partition
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

	// Create an HTTP request for CreateIndex endpoint with the index mapping JSON payload
	reqCreate, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create request for CreateIndex")

	reqCreate.Header.Set("Content-Type", "application/json")

	// Perform the CreateIndex request
	respCreate, err := app.Test(reqCreate, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make CreateIndex request")
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode, "Expected status 201 for CreateIndex")

	// Create test documents
	documents := []map[string]interface{}{
		{"id": "1", "name": "John Doe", "age": 30},
		{"id": "2", "name": "Jane Doe", "age": 25},
	}
	for _, doc := range documents {
		documentJSON, _ := json.Marshal(doc)
		reqIndex, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition+"/index", bytes.NewBuffer(documentJSON))
		assert.NoError(t, err, "Failed to create request for IndexDocument")

		reqIndex.Header.Set("Content-Type", "application/json")

		// Perform the IndexDocument request
		respIndex, err := app.Test(reqIndex, -1) // Timeout in milliseconds
		assert.NoError(t, err, "Failed to make IndexDocument request")
		assert.Equal(t, http.StatusCreated, respIndex.StatusCode, "Expected status 201 for IndexDocument")
	}

	// Create an HTTP request to print the index mapping
	reqPrintMapping, err := http.NewRequest(http.MethodGet, "/api/v1/"+partition+"/mapping", nil)
	assert.NoError(t, err, "Failed to create request for PrintIndexMapping")

	// Perform the PrintIndexMapping request
	respPrintMapping, err := app.Test(reqPrintMapping, -1) // Timeout in milliseconds
	assert.NoError(t, err, "Failed to make PrintIndexMapping request")
	assert.Equal(t, http.StatusOK, respPrintMapping.StatusCode, "Expected status 200 for PrintIndexMapping")

	// Read response body
	bodyPrintMapping, err := io.ReadAll(respPrintMapping.Body)
	assert.NoError(t, err, "Failed to read PrintIndexMapping response body")
	defer respPrintMapping.Body.Close()

	// Assert response message
	var mappingResponse map[string]interface{}
	err = json.Unmarshal(bodyPrintMapping, &mappingResponse)
	assert.NoError(t, err, "Failed to unmarshal PrintIndexMapping response body")

	// Ensure the response contains expected keys and values
	expectedMapping := map[string]interface{}{
		"code":    float64(200), // Ensure code is treated as float64
		"success": true,
		"data": []interface{}{
			map[string]interface{}{
				"analysis":                map[string]interface{}{},
				"default_analyzer":        "standard",
				"default_datetime_parser": "dateTimeOptional",
				"default_field":           "defaultField",
				"default_mapping": map[string]interface{}{
					"dynamic": true,
					"enabled": true,
					"properties": map[string]interface{}{
						"defaultField": map[string]interface{}{
							"dynamic": true,
							"enabled": true,
						},
					},
				},
				"default_type":      "_default",
				"docvalues_dynamic": false,
				"index_dynamic":     true,
				"store_dynamic":     true,
				"type_field":        "_type",
				"types": map[string]interface{}{
					"document": map[string]interface{}{
						"dynamic": true,
						"enabled": true,
						"properties": map[string]interface{}{
							"author": map[string]interface{}{
								"dynamic": true,
								"enabled": true,
							},
							"date": map[string]interface{}{
								"dynamic": true,
								"enabled": true,
							},
							"title": map[string]interface{}{
								"dynamic": true,
								"enabled": true,
							},
						},
					},
				},
			},
		},
	}

	// Compare the mappingResponse with expectedMapping
	assert.Equal(t, expectedMapping, mappingResponse, "Mapping response does not match expected structure")
}

func TestGetPartitionsEndpoint(t *testing.T) {
	// Initialize the Fiber app

	// Create some partitions for testing
	partition1 := generateRandomString(9)
	app := setupApp(partition1)

	// Create index mappings for the test partitions
	indexMappingJSON := `{}`

	// Create the first index
	reqCreate1, err := http.NewRequest(http.MethodPost, "/api/v1/"+partition1+"/create", strings.NewReader(indexMappingJSON))
	assert.NoError(t, err, "Failed to create request for partition 1")
	reqCreate1.Header.Set("Content-Type", "application/json")
	respCreate1, err := app.Test(reqCreate1, -1)
	assert.NoError(t, err, "Failed to make CreateIndex request for partition 1")
	assert.Equal(t, http.StatusCreated, respCreate1.StatusCode, "Expected status 201 for CreateIndex partition 1")

	// Test the GetPartitions endpoint
	reqGet, err := http.NewRequest(http.MethodGet, "/api/v1/", nil)
	assert.NoError(t, err, "Failed to create request for GetPartitions")

	// Perform the GetPartitions request
	respGet, err := app.Test(reqGet, -1) // No timeout
	assert.NoError(t, err, "Failed to make GetPartitions request")
	defer respGet.Body.Close()

	// Assert GetPartitions response status code
	assert.Equal(t, http.StatusOK, respGet.StatusCode, "Expected status 200 for GetPartitions")

	// Read GetPartitions response body
	bodyGet, err := io.ReadAll(respGet.Body)
	assert.NoError(t, err, "Failed to read GetPartitions response body")

	// Unmarshal GetPartitions response body
	var getPartitionsResponse map[string]interface{}
	err = json.Unmarshal(bodyGet, &getPartitionsResponse)
	assert.NoError(t, err, "Failed to unmarshal GetPartitions response body")

	// Assert GetPartitions response success
	assert.True(t, getPartitionsResponse["success"].(bool), "Expected success to be true in GetPartitions response")

	// Assert GetPartitions response data
	data := getPartitionsResponse["data"].([]interface{})
	assert.Contains(t, data, partition1, "Expected partition 1 in GetPartitions response data")

	// Clean up: Delete the test partitions
	reqDelete1, err := http.NewRequest(http.MethodDelete, "/api/v1/"+partition1+"/delete", nil)
	assert.NoError(t, err, "Failed to create request for DeleteIndex partition 1")
	respDelete1, err := app.Test(reqDelete1, -1)
	assert.NoError(t, err, "Failed to make DeleteIndex request for partition 1")
	assert.Equal(t, http.StatusOK, respDelete1.StatusCode, "Expected status 200 for DeleteIndex partition 1")
}
