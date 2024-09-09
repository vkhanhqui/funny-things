# Bleve Search

## Purpose

This project provides a REST API for managing search indexes using Bleve, a modern text indexing library for Go. The API allows users to create and delete indexes, index documents, search for documents, and perform various other operations related to index management. The primary goal is to offer a robust and efficient search solution that can be easily integrated into other applications.

```bash
├── cmd
│   └── restapi
├── go.mod
├── go.sum
├── integration
│   ├── integration_tests
│   ├── intergration_test.go
│   └── poc_test.go
├── internal
│   ├── blevefunc
│   ├── controllers
│   ├── repositories
│   ├── routes
│   └── utils
├── pkg
│   └── log
└── README.md
```

## API Endpoints

* [Create Index](#1-create-index)
* [Delete Index](#2-delete-index)
* [Index Document](#3-index-document)
* [Search Document](#4-search-document)
* [Get Document Count](#5-get-document-count)
* [Bulk Load Documents](#6-bulk-load-documents)
* [Check Index Contents](#7-check-index-contents)
* [Print Term Dictionary](#8-print-term-dictionary)
* [Dump Index Contents](#9-dump-index-contents)
* [List Index Fields](#10-list-index-fields)
* [Print Index Mapping](#11-print-index-mapping)

---
### 1. Create Index

- **Description:** Creates a new index for a specified partition.
- **Method:** `POST`
- **Endpoint:** `/api/v1/{partition}`
- **Request Body:** (example)
```json
{
    "analysis": {
        "token_filters": {
            "to_lower": {
                "type": "to_lower"
            }
        }
    },
    "default_mapping": {
        "dynamic": true,
        "enabled": true,
        "properties": {
            "defaultField": {
                "dynamic": true,
                "enabled": true,
                "fields": [
                    {
                        "name": "defaultField",
                        "type": "text"
                    }
                ]
            }
        }
    },
    "types": {
        "product": {
            "properties": {
                "name": {
                    "fields": [
                        {
                            "name": "name",
                            "type": "text"
                        }
                    ]
                }
            }
        }
    },
    "docvalues_dynamic": false,
    "index_dynamic": true,
    "store_dynamic": true,
    "type_field": "_type",
    "default_type": "_default",
    "default_analyzer": "standard",
    "default_datetime_parser": "dateTimeOptional",
    "default_field": "defaultField"
}
  ```
### Explanation of Fields:
- **`analysis`**: Optional section for defining custom analyzers.
  - **`char_filters`**: Defines character filters to preprocess the text  before tokenization. 
  - **`tokenizers`**: Defines how the text is split into tokens (words or   terms).
  - **`token_maps`**: Defines maps of tokens, such as stop words.
  - **`token_filters`**: Defines filters to apply to tokens after tokenization.
    - List token filters availables:
        - apostrophe
        - camelCase
        - dict_compound
        - edge_ngram
        - elision
        - hierarchy
        - keyword_marker
        - length
        - ngram
        - normalize_unicode
        - possessive_en
        - reverse
        - shingle
        - stemmer_en_plural
        - stemmer_en_snowball
        - stemmer_porter
        - stop_en
        - stop_tokens
        - to_lower
        - truncate_token
        - unique
  - **`analyzers`**: Defines custom analyzers that specify a tokenizer and  a sequence of token filters.
  - **`date_time_parsers`**: Defines custom date-time parsers for handling  different date-time formats.

- **`default_mapping`**: Specifies default mappings for fields not explicitly mapped under specific types.
  - **`enabled`**: Indicates whether the document mapping is enabled. If set to false, the document type will not be indexed.
  - **`dynamic`**: Determines whether the document type allows dynamic addition of new fields that are not explicitly defined in the mapping.
  - **`properties`**: A map of field names to their respective DocumentMapping objects, allowing for nested fields within the document. Each field can have its own mapping configuration.
  - **`fields`**: A list of field mappings that specify how individual fields within the document are indexed and analyzed.
  - **`defaultAnalyzer`**: Specifies the default analyzer to use for fields in this document type, unless overridden by individual field settings.
  - **`structTagKey`**: Overrides the default struct tag key ("json") used when looking for field names in struct tags. This is useful when you want to use a different tag key for field names.

- **`types`**: Defines mappings for specific document types. Each type can have its own set of properties (`product`, `user`, etc.).
  - **`defaultAnalyzer`**: Specifies the default analyzer to use for fields in this   document type, unless overridden by individual field settings.
  - **`structTagKey`**: Overrides the default struct tag key ("json") used when looking   for field names in struct tags. This is useful when you want to use a different tag key   for field names.
  - **`properties`**: A map of field names to their respective `DocumentMapping` objects,   allowing for nested fields within the document. Each field can have its own mapping   configuration.
  - **`enabled`**: Indicates whether the document mapping is enabled. If set to `false`, the document type will not be indexed.
  - **`dynamic`**: Determines whether the document type allows dynamic addition of new  fields that are not explicitly defined in the mapping.
  - **`fields`**: A list of field mappings that specify how individual fields within the  document are indexed and analyzed.
    - **`name`**: The name of the field.
    - **`type`**: The data type of the field (e.g., text, numeric).
    - **`analyzer`**: Specifies the name of the analyzer to use for this field. If empty, it will use the first non-empty DefaultAnalyzer found in the DocumentMapping tree or the IndexMapping.DefaultAnalyzer.
    - **`store`**: Indicates whether to store field values in the index. Stored values can be retrieved from search results using SearchRequest.Fields.
    - **`index`**: Indicates whether the field should be indexed.
    - **`include_term_vectors`**: If true, term occurrences are recorded for this field, including term positions and offsets in the source document field. Required for phrase queries or terms highlighting in source documents.
    - **`include_in_all`**: Indicates whether the field should be included in the _all field.
    - **`date_format`**: Specifies the date format for date fields.
    - **`docvalues`**: If true, enables the index to support uninverting for this field. Useful for faceting and sorting queries.
    - **`skip_freq_norm`**: If true, avoids indexing frequency and norm values of the tokens for this field. Useful for saving processing when default score-based relevancy isn't needed.
    - **`dims`**: Specifies the dimensionality of the vector.
    - **`similarity`**: The similarity algorithm used for scoring vector fields. See index.DefaultSimilarityMetric and index.SupportedSimilarityMetrics.
    - **`vector_index_optimized_for`**: Optimization setting for vector fields.

- **`docvalues_dynamic`**: Controls dynamic mapping behavior for doc values (`true` or `false`).
- **`index_dynamic`**: Controls dynamic mapping behavior for indexing fields (`true` or `false`).
- **`store_dynamic`**: Controls dynamic mapping behavior for storing fields (`true` or `false`).
- **`type_field`**: Defines the field in documents that determines their type (`_type` by default).
- **`default_type`**: Specifies the default type for documents that don't specify a type (`_default` by default).
- **`default_analyzer`**: Specifies the default analyzer to use for text fields (`standard` analyzer by default).
- **`default_datetime_parser`**: Specifies the default date-time parser for date fields (`dateTimeOptional` by default).
- **`default_field`**: Specifies the default field to use if not explicitly defined (`defaultField` by default).

### 2. Delete Index

- **Description:** Deletes an existing index for a specified partition.
- **Method:** `DELETE`
- **Endpoint:** `/api/v1/{partition}`
- **Request Body:** None
  

### 3. Index Document

- **Description:** Indexes a document in the specified partition.
- **Method:** `POST`
- **Endpoint:** `/api/v1/{partition}/document`
- **Request Body:** JSON object representing the document

### 4. Search Document

- **Description:** Searches for documents in the specified partition based on a query string.
- **Method:** `GET`
- **Endpoint:** `/api/v1/{partition}/search?q={query}`
- **Request Body:** None
- **Bleve Query Types and Usage**:
  #### 1. `query_string`

  - **Description:** Executes a query string query.
  - **Parameters:** 
    - `query`: The query string to execute.
  - **Example:**
  ```json
    [{
      "type": "query_string",
      "query": "search term"
    }]
  ```
    This executes a query string search for "search term".

  #### 2. `bool_field`

  - **Description:** Matches documents based on a boolean field value.
  - **Parameters:** 
    - `value`: Boolean value (`true` or `false`).
  - **Example:**
  ```json
    [{
      "type": "bool_field",
      "value": true
    }]
  ```
    Matches documents where the boolean field is `true`.

  #### 3. `boolean`

  - **Description:** Executes a boolean combination of queries (must,   should, must_not).
  - **Parameters:** 
    - `must`: Array of queries that must match.
    - `should`: Array of queries that should match.
    - `must_not`: Array of queries that must not match.
  - **Example:**
  ```json
    [{
      "type": "boolean",
      "must": [
        [{
          "type": "term",
          "term": "keyword"
        }]
      ],
      "should": [
        [{
          "type": "match_phrase",
          "match_phrase": "exact phrase"
        }]
      ],
      "must_not": []
    }]
  ```
    Executes a boolean query with a must clause for term "keyword" and a should clause for exact phrase "exact phrase".

  #### 4. `conjunction`

  - **Description:** Executes a conjunction (AND) of queries.
  - **Parameters:** 
    - `conjuncts`: Array of queries to combine with AND.
  - **Example:**
  ```json
    [{
      "type": "conjunction",
      "conjuncts": [
        [{
          "type": "term",
          "term": "keyword1"
        }],
        [{
          "type": "term",
          "term": "keyword2"
        }]
      ]
    }]
  ```
    Executes a conjunction query combining "keyword1" AND "keyword2".

  #### 5. `date_range`

  - **Description:** Matches documents within a specified date range.
  - **Parameters:** 
    - `start`: Start date in RFC3339 format.
    - `end`: End date in RFC3339 format.
  - **Example:**
  ```json
    [{
      "type": "date_range",
      "start": "2023-01-01T00:00:00Z",
      "end": "2023-12-31T23:59:59Z"
    }]
  ```
    Matches documents within the year 2023.

  #### 6. `disjunction`

  - **Description:** Executes a disjunction (OR) of queries.
  - **Parameters:** 
    - `disjuncts`: Array of queries to combine with OR.
  - **Example:**
  ```json
    [{
      "type": "disjunction",
      "disjuncts": [
        [{
          "type": "term",
          "term": "keyword1"
        }],
        [{
          "type": "term",
          "term": "keyword2"
        }]
      ]
    }]
  ```
    Executes a disjunction query combining "keyword1" OR "keyword2".

  #### 7. `doc_id`

  - **Description:** Matches documents by specific document IDs.
  - **Parameters:** 
    - `ids`: Array of document IDs to match.
  - **Example:**
  ```json
    [{
      "type": "doc_id",
      "ids": ["doc1", "doc2"]
    }]
  ```
    Matches documents with IDs "doc1" or "doc2".

  #### 8. `fuzzy`

  - **Description:** Executes a fuzzy query for approximate matching.
  - **Parameters:** 
    - `term`: Term to match approximately.
  - **Example:**
  ```json
    [{
      "type": "fuzzy",
      "term": "approximate"
    }]
  ```
    Matches documents with terms similar to "approximate".

  #### 9. `match`

  - **Description:** Executes a full-text match query.
  - **Parameters:** 
    - `match`: Text to match.
  - **Example:**
  ```json
    [{
      "type": "match",
      "match": "full text query"
    }]
  ```
    Matches documents containing "full text query".

  #### 10. `match_all`

  - **Description:** Matches all documents.
  - **Example:**
  ```json
    [{
      "type": "match_all"
    }]
  ```
    Matches all documents in the index.

  #### 11. `match_none`

  - **Description:** Matches no documents.
  - **Example:**
  ```json
    [{
      "type": "match_none"
    }]
  ```
    Does not match any documents.

  #### 12. `match_phrase`

  - **Description:** Executes a phrase query for exact matching of phrases.
  - **Parameters:** 
    - `match_phrase`: Exact phrase to match.
  - **Example:**
  ```json
    [{
      "type": "match_phrase",
      "match_phrase": "exact phrase"
    }]
  ```
    Matches documents with the exact phrase "exact phrase".

  #### 13. `numeric_range`

  - **Description:** Matches documents within a specified numeric range.
  - **Parameters:** 
    - `min`: Minimum value (optional).
    - `max`: Maximum value (optional).
  - **Example:**
  ```json
    [{
      "type": "numeric_range",
      "min": 100,
      "max": 200
    }]
  ```
    Matches numeric values between 100 and 200.

  #### 14. `phrase`

  - **Description:** Executes a phrase query for matching exact phrases.
  - **Parameters:** 
    - `terms`: Array of terms in the exact phrase.
    - `field`: Field in which to search for the phrase.
  - **Example:**
  ```json
    [{
      "type": "phrase",
      "terms": ["exact", "phrase"],
      "field": "content"
    }]
  ```
    Matches documents with the exact phrase "exact phrase" in the "content" field.

  #### 15. `prefix`

  - **Description:** Matches documents with terms starting with a   specified prefix.
  - **Parameters:** 
    - `prefix`: Prefix to match.
  - **Example:**
  ```json
    [{
      "type": "prefix",
      "prefix": "pre"
    }]
  ```
    Matches documents with terms starting with "pre".

  #### 16. `regexp`

  - **Description:** Executes a regular expression query.
  - **Parameters:** 
    - `regexp`: Regular expression pattern to match.
  - **Example:**
  ```json
    [{
      "type": "regexp",
      "regexp": "^prefix.*"
    }]
  ```
    Matches documents with terms matching the regular expression "^prefix.*".

  #### 17. `term`

  - **Description:** Executes a term query for exact matching.
  - **Parameters:** 
    - `term`: Term to match.
  - **Example:**
  ```json
    [{
      "type": "term",
      "term": "exact_term"
    }]
  ```
    Matches documents with the exact term "exact_term".

  #### 18. `wildcard`

  - **Description:** Executes a wildcard query for matching terms using   wildcards.
  - **Parameters:** 
    - `wildcard`: Wildcard pattern to match.
  - **Example:**
  ```json
    [{
      "type": "wildcard",
      "wildcard": "term*"
    }]
  ```
    Matches documents with terms matching the wildcard pattern "term*".

  #### 19. `geo_bounding_box`

  - **Description:** Matches documents within a specified geographic  bounding box.
  - **Parameters:** 
    - `top_left_lon`: Longitude of the top-left corner.
    - `top_left_lat`: Latitude of the top-left corner.
    - `bottom_right_lon`: Longitude of the bottom-right corner.
    - `bottom_right_lat`: Latitude of the bottom-right corner.
  - **Example:**
  ```json
    [{
      "type": "geo_bounding_box",
      "top_left_lon": -74.04728500751165,
      "top_left_lat": 40.88221709544296,
      "bottom_right_lon": -73.90757083892822,
      "bottom_right_lat": 40.68228753441558
    }]
  ```
    Matches documents within the specified geographic bounding box.

  #### 20. `geo_distance`

  - **Description:** Matches documents within a specified distance from a   geographic point.
  - **Parameters:** 
    - `lon`: Longitude of the reference point.
    - `lat`: Latitude of the reference point.
    - `distance`: Distance from the reference point (e.g., "10km").
  - **Example:**
  ```json
    [{
      "type": "geo_distance",
      "lon": -73.9851303100586,
      "lat": 40.748817610874455,
      "distance": "10km"
    }]
  ```
    Matches documents within 10 kilometers from the specified geographic point.

  #### 21. `ip_range`

  - **Description:** Matches documents within a specified IP address range.
  - **Parameters:** 
    - `cidr`: CIDR notation IP address range (e.g., "192.168.1.0/24").
  - **Example:**
  ```json
   [{
      "type": "ip_range",
      "cidr": "192.168.1.0/24"
    }]
  ```
    Matches documents within the specified IP address range.

  #### 22. `geo_shape`

  - **Description:** Matches documents with a specified geographic shape.
  - **Parameters:** 
    - `coordinates`: Array of coordinates defining the shape.
    - `type`: Type of the shape (e.g., "polygon", "circle").
    - `relation`: Spatial relation to the shape (e.g., "within", "intersects").
  - **Example:**
  ```json
    [{
      "type": "geo_shape",
      "coordinates": [[[-73.99756, 40.73083], [-73.99756, 40.741404], [-73.988135, 40.741404], [-73.988135, 40.73083], [-73.99756, 40.73083]]],
      "type": "polygon",
      "relation": "within"
    }]
  ```
    Matches documents within the specified polygon shape.

- **Query parameters**:

  - Query
    - **type**: Specifies the type of query.
    - **term**: Specifies the term to search for.

  - Fields
    - **fields**: Specifies the fields to include in the search results.

  - Pagination
    - **size**: Specifies the maximum number of results to return.
    - **from**: Specifies the offset from which to start returning results.

  - Sorting
    - **sort**:
      - **field**: Specifies the field to sort by.
      - **order**: Specifies the sorting order (`asc` for ascending, `desc` for descending).

  - Facets
  Facets provide aggregated data based on specified fields or ranges.
    - **titles**:
      - **size**: Specifies the maximum number of facet terms to return.
      - **field**: Specifies the field for which to generate facet terms.
      - **numericRanges**: Optional. Specifies numeric ranges for aggregating facet terms.
      - **dateTimeRanges**: Optional. Specifies date ranges for aggregating facet terms.
  Example:
  ```json
  [
      {
          "query": {
              "type": "term",
              "term": "title1"
          },
          "fields": [
              "name",
              "age",
              "title",
              "content"
          ],
          "size": 2,
          "from": 0,
          "sort": {
              "field": "id",
              "order": "asc"
          },
          "facets": {
              "age_ranges": {
                  "field": "age",
                  "numeric_ranges": [
                      {
                          "name": "young",
                          "max": 29
                      },
                      {
                          "name": "middle_aged",
                          "min": 30,
                          "max": 39
                      },
                      {
                          "name": "senior",
                          "min": 40
                      }
                  ]
              },
              "titles": {
                  "size": 10,
                  "field": "title"
              }
          }
      }
  ]
  ```


### 5. Get Document Count

- **Description:** Retrieves the count of documents in the specified partition.
- **Method:** `GET`
- **Endpoint:** `/api/v1/{partition}/count`
- **Request Body:** None
  

### 6. Bulk Load Documents

- **Description:** Bulk loads documents into the specified partition from a provided data source.
- **Method:** `POST`
- **Endpoint:** `/api/v1/{partition}/bulk`
- **Request Body:** (example)
  ```json
    {"id": "1", "name": "John Doe", "age": 30, "title":"some title1", "content":"some content1"}
    {"id": "2", "name": "Jane Doe", "age": 25, "title":"some title2", "content":"some content2"}
    {"id": "3", "name": "Jake Smith", "age": 40, "title":"some title3", "content":"some content3"}
  ```
  

### 7. Check Index Contents

- **Description:** Retrieves detailed contents of the specified index.
- **Method:** `GET`
- **Endpoint:** `/api/v1/{partition}/contents`
- **Request Body:** None

### 8. Print Term Dictionary

- **Description:** Retrieves the term dictionary for a specified field in the index.
- **Method:** `GET`
- **Endpoint:** `/api/v1/{partition}/terms?field={field}`
- **Request Body:** None

### 9. Dump Index Contents

- **Description:** Dumps all documents in the specified index.
- **Method:** `GET`
- **Endpoint:** `/api/v1/{partition}/dump`
- **Request Body:** None
  

### 10. List Index Fields

- **Description:** Lists all fields in the specified index.
- **Method:** `GET`
- **Endpoint:** `/api/v1/{partition}/fields`
- **Request Body:** None

### 11. Print Index Mapping

- **Description:** Retrieves the mapping of the specified index.
- **Method:** `GET`
- **Endpoint:** `/api/v1/{partition}/mapping`
- **Request Body:** None

### 12. Get Partitions
- **Description:** Retrieves a list of all available index partitions.
- **Method:** `GET`
- **Endpoint:** `/api/v1/partitions`
- **Request Body:** None

---