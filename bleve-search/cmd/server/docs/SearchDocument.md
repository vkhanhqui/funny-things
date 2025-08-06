### Search Document

## **Bleve Query Types and Usage**:
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

---
## **Query parameters**:

### Query
- **type**: Specifies the type of query.
- **term**: Specifies the term to search for.

### Fields
- **fields**: Specifies the fields to include in the search results.

### Pagination
- **size**: Specifies the maximum number of results to return.
- **from**: Specifies the offset from which to start returning results.

### Sorting
- **sort**:
  - **field**: Specifies the field to sort by.
  - **order**: Specifies the sorting order (`asc` for ascending, `desc` for descending).

### Facets
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
        }
    }
]
```
