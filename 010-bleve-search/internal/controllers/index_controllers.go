package controllers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"bleve-proj/internal/repositories"
	"bleve-proj/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CreateIndex creates a new index in the repository with custom mapping
//
//	@Summary				Create a new index
//	@Description.Markdown	CreateIndex
//	@Tags					Bleve
//	@Accept					json
//	@Produce				json
//	@Param					partition		path		string		true	"Partition Name"
//	@Param					indexMapping	body		object		true	"Index Mapping JSON"
//	@Success				201				{object}	Response	"Index created successfully"
//	@Failure				400				{object}	Response	"Invalid JSON format"
//	@Failure				500				{object}	Response	"Internal Server Error"
//	@Router					/{partition}/create [post]
func CreateIndex(c *fiber.Ctx) error {
	partition := c.Params("partition")

	// Assuming the index mapping JSON is provided in the request body
	var indexMappingJSON map[string]interface{}
	if err := c.BodyParser(&indexMappingJSON); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
	}

	// Create index using repository
	if _, err := repositories.IndexRepo.CreateIndex(partition, indexMappingJSON); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, http.StatusCreated, "Index created successfully")
}

// DeleteIndex deletes an index from the repository
//
//	@Summary	Delete an index
//	@Tags		Bleve
//	@Accept		json
//	@Produce	json
//	@Param		partition	path		string		true	"Partition Name"
//	@Success	200			{object}	Response	"Index deleted successfully"
//	@Failure	404			{object}	Response	"Index not found"
//	@Failure	500			{object}	Response	"Internal Server Error"
//	@Router		/{partition}/delete [delete]
func DeleteIndex(c *fiber.Ctx) error {
	partition := c.Params("partition")
	if err := repositories.IndexRepo.DeleteIndex(partition); err != nil {
		if errors.Is(err, repositories.ErrIndexNotFound) {
			return ErrorResponse(c, http.StatusNotFound, "index not found")
		}
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, http.StatusOK, "Index deleted successfully")
}

// IndexDocument indexes a document in the specified index partition
//
//	@Summary				Index a document
//	@Description.Markdown	IndexDocument
//	@Tags					Bleve
//	@Accept					json
//	@Produce				json
//	@Param					partition	path		string		true	"Partition Name"
//	@Param					doc			body		object		true	"Document JSON"
//	@Success				201			{object}	Response	"Document indexed successfully"
//	@Failure				400			{object}	Response	"Invalid JSON"
//	@Router					/{partition}/index [post]
func IndexDocument(c *fiber.Ctx) error {
	partition := c.Params("partition")
	var doc map[string]interface{}
	if err := c.BodyParser(&doc); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "Invalid JSON")
	}

	docID, ok := doc["id"].(string)
	if !ok || docID == "" {
		docID = utils.GenerateUniqueID()
		doc["id"] = docID
	}
	// IndexDocument runs as a goroutine
	if err := repositories.IndexRepo.IndexDocument(partition, docID, doc); err != nil {
		// Log or handle error as needed
		fmt.Printf("IndexDocument Error encountered: %v", err)
	}

	return SuccessResponse(c, http.StatusCreated, "Document indexed successfully")
}

// SearchDocument searches for documents in the specified index partition based on a query string.
//
//	@Summary				Search documents
//	@Description.Markdown	SearchDocument
//	@Tags					Bleve
//	@Accept					json
//	@Produce				json
//	@Param					partition	path		string		true	"Partition Name"
//	@Param					q			query		string		true	"Query string"
//	@Success				200			{object}	Response	"Search results"
//	@Failure				400			{object}	Response	"Invalid query"
//	@Failure				500			{object}	Response	"Internal Server Error"
//	@Router					/{partition}/search [get]
func SearchDocument(c *fiber.Ctx) error {
	partition := c.Params("partition")
	queryString := c.Query("q")

	result, err := repositories.IndexRepo.Search(partition, queryString)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, http.StatusOK, result)
}

// GetDocumentCount retrieves the count of documents in the specified index partition.
//
//	@Summary		Get document count
//	@Description	Retrieves the count of documents in the specified index partition.
//	@Tags			Bleve
//	@Accept			json
//	@Produce		json
//	@Param			partition	path		string		true	"Partition Name"
//	@Success		200			{object}	Response	"Document count retrieved successfully"
//	@Failure		500			{object}	Response	"Internal Server Error"
//	@Router			/{partition}/count [get]
func GetDocumentCount(c *fiber.Ctx) error {
	partition := c.Params("partition")
	count, err := repositories.IndexRepo.GetDocumentCount(partition)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, http.StatusOK, count)
}

// BulkLoadDocuments loads multiple documents in bulk into the specified index partition.
//
//	@Summary				Bulk load documents
//	@Description.Markdown	BulkDocument
//	@Tags					Bleve
//	@Accept					json
//	@Produce				json
//	@Param					partition	path		string						true	"Partition Name"
//	@Param					body		body		[]map[string]interface{}	true	"List of documents to be loaded"
//	@Success				200			{object}	Response					"Bulk documents loaded successfully"
//	@Failure				500			{object}	Response					"Internal Server Error"
//	@Router					/{partition}/bulk [post]
func BulkLoadDocuments(c *fiber.Ctx) error {
	partition := c.Params("partition")
	bodyReader := bytes.NewReader(c.Body())
	scanner := bufio.NewScanner(bodyReader)

	// BulkLoadDocuments runs as a goroutine
	go func() {
		if err := repositories.IndexRepo.BulkLoadDocuments(partition, scanner); err != nil {
			// Log or handle error as needed
			fmt.Printf("BulkLoadDocuments Error encountered: %v", err)
		}
	}()

	return SuccessResponse(c, http.StatusOK, "Bulk documents loaded successfully")
}

// CheckIndexContents checks the contents of the specified index partition
//
//	@Summary		Check index contents
//	@Description	Checks the contents of the specified index partition.
//	@Tags			Bleve
//	@Accept			json
//	@Produce		json
//	@Param			partition	path		string		true	"Partition Name"
//	@Success		200			{object}	Response	"Index contents retrieved successfully"
//	@Failure		500			{object}	Response	"Internal Server Error"
//	@Router			/{partition}/check [get]
func CheckIndexContents(c *fiber.Ctx) error {
	partition := c.Params("partition")
	contents, err := repositories.IndexRepo.CheckIndexContents(partition)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return SuccessResponse(c, http.StatusOK, contents)
}

// PrintTermDictionary prints the term dictionary for a specified field in the index partition
//
//	@Summary		Print term dictionary
//	@Description	Prints the term dictionary for a specified field in the index partition.
//	@Tags			Bleve
//	@Accept			json
//	@Produce		json
//	@Param			partition	path		string		true	"Partition Name"
//	@Param			field		query		string		true	"Field Name"
//	@Success		200			{object}	Response	"Term dictionary retrieved successfully"
//	@Failure		400			{object}	Response	"Field parameter is required"
//	@Failure		500			{object}	Response	"Internal Server Error"
//	@Router			/{partition}/dictionary [get]
func PrintTermDictionary(c *fiber.Ctx) error {
	partition := c.Params("partition")
	field := c.Query("field")

	if field == "" {
		return ErrorResponse(c, http.StatusBadRequest, "Field parameter is required")
	}

	termDict, err := repositories.IndexRepo.PrintTermDictionary(partition, field)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, http.StatusOK, termDict)
}

// DumpIndexContents dumps the contents of the specified index partition with pagination
//
//	@Summary		Dump index contents
//	@Description	Dumps the contents of the specified index partition with pagination.
//	@Tags			Bleve
//	@Accept			json
//	@Produce		json
//	@Param			partition	path		string		true	"Partition Name"
//	@Param			page		query		string		false	"Page number, defaults to 1"
//	@Param			size		query		string		false	"Page size, defaults to 1000"
//	@Success		200			{object}	Response	"Index contents retrieved successfully"
//	@Failure		400			{object}	Response	"Invalid page number or page size"
//	@Failure		404			{object}	Response	"Index not found"
//	@Failure		500			{object}	Response	"Internal Server Error"
//	@Router			/{partition}/dump [get]
func DumpIndexContents(c *fiber.Ctx) error {
	partition := c.Params("partition")
	page := c.Query("page", "1")    // Default to first page
	size := c.Query("size", "1000") // Default to fetching 1000 documents

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	sizeNum, err := strconv.Atoi(size)
	if err != nil || sizeNum <= 0 {
		sizeNum = 1000
	}

	// Dump index contents using repository method with offset and limit
	documents, err := repositories.IndexRepo.DumpIndexContentsWithLimit(partition, pageNum, sizeNum)
	if err != nil {
		// If index not found, return 404
		if errors.Is(err, repositories.ErrIndexNotFound) {
			return ErrorResponse(c, http.StatusNotFound, "index not found")
		}
		// Handle other errors
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Return success response with JSON
	return SuccessResponse(c, http.StatusOK, documents)
}

// ListIndexFields lists all fields in the specified index partition
//
//	@Summary		List index fields
//	@Description	Lists all fields in the specified index partition.
//	@Tags			Bleve
//	@Accept			json
//	@Produce		json
//	@Param			partition	path		string		true	"Partition Name"
//	@Success		200			{object}	Response	"Fields listed successfully"
//	@Failure		404			{object}	Response	"Index not found"
//	@Failure		500			{object}	Response	"Internal Server Error"
//	@Router			/{partition}/fields [get]
func ListIndexFields(c *fiber.Ctx) error {
	partition := c.Params("partition")
	fields, err := repositories.IndexRepo.ListIndexFields(partition)
	if err != nil {
		if errors.Is(err, repositories.ErrIndexNotFound) {
			return ErrorResponse(c, http.StatusNotFound, "index not found")
		}
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, http.StatusOK, fields)
}

// PrintIndexMapping prints the mapping for the specified index partition
//
//	@Summary		Print index mapping
//	@Description	Prints the mapping for the specified index partition.
//	@Tags			Bleve
//	@Accept			json
//	@Produce		json
//	@Param			partition	path		string		true	"Partition Name"
//	@Success		200			{object}	Response	"Mapping printed successfully"
//	@Failure		500			{object}	Response	"Internal Server Error"
//	@Router			/{partition}/mapping [get]
func PrintIndexMapping(c *fiber.Ctx) error {
	partition := c.Params("partition")
	mapping, err := repositories.IndexRepo.PrintIndexMapping(partition)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, http.StatusOK, mapping)
}

// GetPartitions retrieves a list of all index partitions
//
//	@Summary		Get partitions
//	@Description	Retrieves a list of all index partitions.
//	@Tags			Bleve
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response	"Partitions retrieved successfully"
//	@Router			/ [get]
func GetPartitions(c *fiber.Ctx) error {
	partitions := repositories.IndexRepo.GetPartitions()
	return SuccessResponse(c, http.StatusOK, partitions)
}
