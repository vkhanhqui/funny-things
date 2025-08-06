# Bleve Search - gRPC API

## Purpose

This project provides a gRPC API for managing search indexes using Bleve, a modern text indexing library for Go. The API allows users to create and delete indexes, index documents, search for documents, and perform various other operations related to index management. The primary goal is to offer a robust and efficient search solution that can be easily integrated into other applications.

## Client Usage

Clients can use the generated gRPC client code to interact with the server. See the gRPC documentation for your preferred language for details on how to generate client code and make gRPC calls.

```bash
├── cmd
│   └── grpc
├── go.mod
├── go.sum
├── integration
│   ├── integration_tests
│   ├── intergration_test.go
│   └── poc_test.go
├── internal
│   ├── blevefunc
│   ├── controllers
│   ├── grpc
│   └── utils
├── pkg
│   └── log
├── proto
│   ├── index_grpc.pb.go
│   ├── index.pb.go
│   └── index.proto
└── README.md
```
## Running the gRPC Server

1.  **Environment Variables:**
    -   `INDEX_DIRECTORY`: The directory where index data will be stored. Default: `./index`
    -   `GRPC_PORT`: The gRPC server port. Default: `:50051`

2.  **Run:**
    ```bash
    go run main.go
    ```

## gRPC Service Definition

The gRPC service definition is located in the `proto/index.proto` file.

### Service: `IndexService`

#### RPC Methods:

1.  **`CreateIndex`**
    -   **Request:** `CreateIndexRequest`
        -   `partition`: `string` - The partition for the index.
        -   `indexMapping`: `google.protobuf.Struct` - The index mapping configuration.
    -   **Response:** `CreateIndexResponse`
        -   `message`: `string` - Success or error message.
2.  **`DeleteIndex`**
    -   **Request:** `DeleteIndexRequest`
        -   `partition`: `string` - The partition of the index to delete.
    -   **Response:** `DeleteIndexResponse`
        -   `message`: `string` - Success or error message.
3.  **`IndexDocument`**
    -   **Request:** `IndexDocumentRequest`
        -   `partition`: `string` - The partition of the index.
        -   `docId`: `string` - The document ID.
        -   `document`: `google.protobuf.Struct` - The document to index.
    -   **Response:** `IndexDocumentResponse`
        -   `message`: `string` - Success or error message.
4.  **`Search`**
    -   **Request:** `SearchRequest`
        -   `partition`: `string` - The partition of the index to search.
        -   `query`: `string` - The search query string (in JSON format representing a Bleve search request).
    -   **Response:** `SearchResult`
        -   `hits`: `repeated SearchHit` - The list of search hits.
        -   `total`: `int32` - The total number of hits.
        -   `max_score`: `double` - The maximum score among the hits.
        -   `took`: `int32` - The time taken for the search (in milliseconds).
5.  **`GetDocumentCount`**
    -   **Request:** `GetDocumentCountRequest`
        -   `partition`: `string` - The partition of the index.
    -   **Response:** `GetDocumentCountResponse`
        -   `count`: `uint64` - The number of documents in the index.
6.  **`BulkLoadDocuments`**
    -   **Request:** `BulkLoadDocumentsRequest`
        -   `partition`: `string` - The partition of the index.
        -   `documents`: `bytes` - The newline-delimited JSON data containing the documents to load.
    -   **Response:** `BulkLoadDocumentsResponse`
        -   `message`: `string` - Success or error message.
7.  **`CheckIndexContents`**
    -   **Request:** `CheckIndexContentsRequest`
        -   `partition`: `string` - The partition of the index.
    -   **Response:** `CheckIndexContentsResponse`
        -   `fields`: `map<string, FieldInfo>` - A map of field names to `FieldInfo`, which contains information about the field (e.g., term count).
8.  **`PrintTermDictionary`**
    -   **Request:** `PrintTermDictionaryRequest`
        -   `partition`: `string` - The partition of the index.
        -   `field`: `string` - The field for which to retrieve the term dictionary.
    -   **Response:** `PrintTermDictionaryResponse`
        -   `termCounts`: `map<string, uint64>` - A map of terms to their frequencies in the specified field.
9.  **`DumpIndexContentsWithLimit`**
    -   **Request:** `DumpIndexContentsWithLimitRequest`
        -   `partition`: `string` - The partition of the index.
        -   `page`: `int32` - The page number (1-based).
        -   `size`: `int32` - The number of documents per page.
    -   **Response:** `DumpIndexContentsWithLimitResponse`
        -   `documents`: `repeated google.protobuf.Struct` - The list of documents (as `google.protobuf.Struct`).
10. **`ListIndexFields`**
    -   **Request:** `ListIndexFieldsRequest`
        -   `partition`: `string` - The partition of the index.
    -   **Response:** `ListIndexFieldsResponse`
        -   `fields`: `repeated string` - The list of fields in the index.
11. **`PrintIndexMapping`**
    -   **Request:** `PrintIndexMappingRequest`
        -   `partition`: `string` - The partition of the index.
    -   **Response:** `PrintIndexMappingResponse`
        -   `mapping`: `google.protobuf.Struct` - The index mapping as a `google.protobuf.Struct`.
12. **`GetPartitions`**
    -   **Request:** `GetPartitionsRequest` (empty message).
    -   **Response:** `GetPartitionsResponse`
        -   `partitions`: `repeated string` - The list of available index partitions.

