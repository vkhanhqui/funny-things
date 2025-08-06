### Create Index

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
