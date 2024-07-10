package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/blevesearch/bleve/v2/search/query"
)

// ParseQueries parses a slice of query maps into a slice of bleve.Query
func ParseQueries(queries []interface{}) []query.Query {
	var result []query.Query
	for _, q := range queries {
		queryMap := q.(map[string]interface{})
		queryType := queryMap["type"].(string)
		queryFunc, exists := queryFuncMap[queryType]
		if exists {
			bQuery, _ := queryFunc(queryMap)
			result = append(result, bQuery)
		}
	}
	return result
}

// Helper function to parse a slice of queries
func parseQuerySlice(param interface{}) ([]query.Query, error) {
	if slice, ok := param.([]interface{}); ok {
		return ParseQueries(slice), nil
	}
	return nil, errors.New("invalid query format")
}

// toStringSlice converts a slice of interfaces to a slice of strings
func toStringSlice(slice []interface{}) []string {
	var result []string
	for _, v := range slice {
		result = append(result, v.(string))
	}
	return result
}

// parseCoordinates parses nested slices of interfaces into slices
func parseCoordinates(coords []interface{}) [][][][]float64 {
	var result [][][][]float64
	for _, c := range coords {
		coord1 := c.([]interface{})
		var coord1Result [][][]float64
		for _, c1 := range coord1 {
			coord2 := c1.([]interface{})
			var coord2Result [][]float64
			for _, c2 := range coord2 {
				coord3 := c2.([]interface{})
				var coord3Result []float64
				for _, c3 := range coord3 {
					coord3Result = append(coord3Result, c3.(float64))
				}
				coord2Result = append(coord2Result, coord3Result)
			}
			coord1Result = append(coord1Result, coord2Result)
		}
		result = append(result, coord1Result)
	}
	return result
}

var queryFuncMap map[string]func(map[string]interface{}) (query.Query, error)

func init() {
	queryFuncMap = map[string]func(map[string]interface{}) (query.Query, error){
		"query_string": func(params map[string]interface{}) (query.Query, error) {
			queryString, ok := params["query"].(string)
			if !ok {
				return nil, errors.New("query string is required")
			}
			return bleve.NewQueryStringQuery(queryString), nil
		},
		"bool_field": func(params map[string]interface{}) (query.Query, error) {
			val, ok := params["value"].(bool)
			if !ok {
				return nil, errors.New("boolean value is required")
			}
			return bleve.NewBoolFieldQuery(val), nil
		},
		"boolean": func(params map[string]interface{}) (query.Query, error) {
			var err error

			// Check for 'must' queries
			var mustQueries []query.Query
			if mustParam, exists := params["must"]; exists {
				mustQueries, err = parseQuerySlice(mustParam)
				if err != nil {
					return nil, err
				}
			}

			// Check for 'should' queries
			var shouldQueries []query.Query
			if shouldParam, exists := params["should"]; exists {
				shouldQueries, err = parseQuerySlice(shouldParam)
				if err != nil {
					return nil, err
				}
			}

			// Check for 'must_not' queries
			var mustNotQueries []query.Query
			if mustNotParam, exists := params["must_not"]; exists {
				mustNotQueries, err = parseQuerySlice(mustNotParam)
				if err != nil {
					return nil, err
				}
			}

			// Create a new BooleanQuery and add the parsed queries
			bQuery := bleve.NewBooleanQuery()
			bQuery.AddMust(mustQueries...)
			bQuery.AddShould(shouldQueries...)
			bQuery.AddMustNot(mustNotQueries...)

			return bQuery, nil
		},
		"conjunction": func(params map[string]interface{}) (query.Query, error) {
			conjuncts := ParseQueries(params["conjuncts"].([]interface{}))
			return bleve.NewConjunctionQuery(conjuncts...), nil
		},
		"date_range": func(params map[string]interface{}) (query.Query, error) {
			start, _ := time.Parse(time.RFC3339, params["start"].(string))
			end, _ := time.Parse(time.RFC3339, params["end"].(string))
			return bleve.NewDateRangeQuery(start, end), nil
		},
		"disjunction": func(params map[string]interface{}) (query.Query, error) {
			disjuncts := ParseQueries(params["disjuncts"].([]interface{}))
			return bleve.NewDisjunctionQuery(disjuncts...), nil
		},
		"doc_id": func(params map[string]interface{}) (query.Query, error) {
			ids := toStringSlice(params["ids"].([]interface{}))
			return bleve.NewDocIDQuery(ids), nil
		},
		"fuzzy": func(params map[string]interface{}) (query.Query, error) {
			term, ok := params["term"].(string)
			if !ok {
				return nil, errors.New("term string is required")
			}
			return bleve.NewFuzzyQuery(term), nil
		},
		"match": func(params map[string]interface{}) (query.Query, error) {
			match, ok := params["match"].(string)
			if !ok {
				return nil, errors.New("match string is required")
			}
			return bleve.NewMatchQuery(match), nil
		},
		"match_all": func(params map[string]interface{}) (query.Query, error) {
			return bleve.NewMatchAllQuery(), nil
		},
		"match_none": func(params map[string]interface{}) (query.Query, error) {
			return bleve.NewMatchNoneQuery(), nil
		},
		"match_phrase": func(params map[string]interface{}) (query.Query, error) {
			matchPhrase, ok := params["match_phrase"].(string)
			if !ok {
				return nil, errors.New("match phrase string is required")
			}
			return bleve.NewMatchPhraseQuery(matchPhrase), nil
		},
		"numeric_range": func(params map[string]interface{}) (query.Query, error) {
			min, okMin := params["min"].(float64)
			max, okMax := params["max"].(float64)
			var minPtr, maxPtr *float64
			if okMin {
				minPtr = &min
			}
			if okMax {
				maxPtr = &max
			}
			return bleve.NewNumericRangeQuery(minPtr, maxPtr), nil
		},
		"phrase": func(params map[string]interface{}) (query.Query, error) {
			terms := toStringSlice(params["terms"].([]interface{}))
			field, ok := params["field"].(string)
			if !ok {
				return nil, errors.New("field string is required")
			}
			return bleve.NewPhraseQuery(terms, field), nil
		},
		"prefix": func(params map[string]interface{}) (query.Query, error) {
			prefix, ok := params["prefix"].(string)
			if !ok {
				return nil, errors.New("prefix string is required")
			}
			return bleve.NewPrefixQuery(prefix), nil
		},
		"regexp": func(params map[string]interface{}) (query.Query, error) {
			regexp, ok := params["regexp"].(string)
			if !ok {
				return nil, errors.New("regexp string is required")
			}
			return bleve.NewRegexpQuery(regexp), nil
		},
		"term": func(params map[string]interface{}) (query.Query, error) {
			term, ok := params["term"].(string)
			if !ok {
				return nil, errors.New("term string is required")
			}
			return bleve.NewTermQuery(term), nil
		},
		"wildcard": func(params map[string]interface{}) (query.Query, error) {
			wildcard, ok := params["wildcard"].(string)
			if !ok {
				return nil, errors.New("wildcard string is required")
			}
			return bleve.NewWildcardQuery(wildcard), nil
		},
		"geo_bounding_box": func(params map[string]interface{}) (query.Query, error) {
			topLeftLon, _ := params["top_left_lon"].(float64)
			topLeftLat, _ := params["top_left_lat"].(float64)
			bottomRightLon, _ := params["bottom_right_lon"].(float64)
			bottomRightLat, _ := params["bottom_right_lat"].(float64)
			return bleve.NewGeoBoundingBoxQuery(topLeftLon, topLeftLat, bottomRightLon, bottomRightLat), nil
		},
		"geo_distance": func(params map[string]interface{}) (query.Query, error) {
			lon, _ := params["lon"].(float64)
			lat, _ := params["lat"].(float64)
			distance, _ := params["distance"].(string)
			return bleve.NewGeoDistanceQuery(lon, lat, distance), nil
		},
		"ip_range": func(params map[string]interface{}) (query.Query, error) {
			cidr, _ := params["cidr"].(string)
			return bleve.NewIPRangeQuery(cidr), nil
		},
		"geo_shape": func(params map[string]interface{}) (query.Query, error) {
			coordinates := parseCoordinates(params["coordinates"].([]interface{}))
			typ, _ := params["type"].(string)
			relation, _ := params["relation"].(string)
			return bleve.NewGeoShapeQuery(coordinates, typ, relation)
		},
	}
}

func parseSortOrder(sortPart map[string]interface{}) (search.SortOrder, error) {
	sortField, ok := sortPart["field"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid sort field")
	}

	order, ok := sortPart["order"].(string)
	if !ok {
		order = "asc" // Default to ascending order if not specified
	}
	return search.SortOrder{
		&search.SortField{
			Field: sortField,
			Desc:  order == "desc",
		},
	}, nil
}

// extractQueries extracts and parses the query part from the queries array
func extractQueries(searchRequest *bleve.SearchRequest, queries []map[string]interface{}) error {
	var bleveQueries []query.Query

	for _, queryData := range queries {
		queryPart, ok := queryData["query"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid query format")
		}

		parsedQueries := ParseQueries([]interface{}{queryPart})
		if len(parsedQueries) > 0 {
			bleveQueries = append(bleveQueries, parsedQueries[0])
		}
	}

	searchRequest.Query = bleve.NewConjunctionQuery(bleveQueries...)
	return nil
}

// extractFields extracts and sets the fields part from the queries array
func extractFields(searchRequest *bleve.SearchRequest, queries []map[string]interface{}) error {
	var fields []string

	for _, queryData := range queries {
		if fieldsPart, ok := queryData["fields"].([]interface{}); ok {
			for _, field := range fieldsPart {
				if fieldStr, ok := field.(string); ok {
					fields = append(fields, fieldStr)
				}
			}
		}
	}

	if len(fields) > 0 {
		searchRequest.Fields = fields
	}
	return nil
}

// extractPagination extracts and sets the size and from parts from the queries array
func extractPagination(searchRequest *bleve.SearchRequest, queries []map[string]interface{}) error {
	for _, queryData := range queries {
		if sizePart, ok := queryData["size"].(float64); ok {
			searchRequest.Size = int(sizePart)
		}

		if fromPart, ok := queryData["from"].(float64); ok {
			searchRequest.From = int(fromPart)
		}
	}
	return nil
}

// extractMiscellaneous extracts and sets includeLocations, score, searchAfter, searchBefore, and explain parts
func extractMiscellaneous(searchRequest *bleve.SearchRequest, queries []map[string]interface{}) error {
	for _, queryData := range queries {
		if includeLocationsPart, ok := queryData["includeLocations"].(bool); ok {
			searchRequest.IncludeLocations = includeLocationsPart
		}

		if scorePart, ok := queryData["score"].(string); ok {
			searchRequest.Score = scorePart
		}

		if searchAfterPart, ok := queryData["search_after"].([]interface{}); ok {
			var searchAfter []string
			for _, item := range searchAfterPart {
				if itemStr, ok := item.(string); ok {
					searchAfter = append(searchAfter, itemStr)
				}
			}
			searchRequest.SearchAfter = searchAfter
		}

		if searchBeforePart, ok := queryData["search_before"].([]interface{}); ok {
			var searchBefore []string
			for _, item := range searchBeforePart {
				if itemStr, ok := item.(string); ok {
					searchBefore = append(searchBefore, itemStr)
				}
			}
			searchRequest.SearchBefore = searchBefore
		}

		if explainPart, ok := queryData["explain"].(bool); ok {
			searchRequest.Explain = explainPart
		}
	}
	return nil
}

// extractFacets extracts and sets facets from the queries array
func extractFacets(searchRequest *bleve.SearchRequest, queries []map[string]interface{}) error {
	facets := make(bleve.FacetsRequest)

	for _, queryData := range queries {
		if facetsPart, ok := queryData["facets"].(map[string]interface{}); ok {
			for facetName, facetData := range facetsPart {
				facetJSON, err := json.Marshal(facetData)
				if err != nil {
					return err
				}

				var facetRequest bleve.FacetRequest
				if err := json.Unmarshal(facetJSON, &facetRequest); err != nil {
					return err
				}
				facets[facetName] = &facetRequest
			}
		}
	}

	if len(facets) > 0 {
		searchRequest.Facets = facets
	}
	return nil
}

// extractSortOrders extracts and sets sort orders from the queries array
func extractSortOrders(searchRequest *bleve.SearchRequest, queries []map[string]interface{}) error {
	var sortOrders []search.SortOrder

	for _, queryData := range queries {
		if sortPart, ok := queryData["sort"].(map[string]interface{}); ok {
			sortOrder, err := parseSortOrder(sortPart)
			if err != nil {
				return err
			}
			sortOrders = append(sortOrders, sortOrder)
		}
	}

	if len(sortOrders) > 0 {
		searchRequest.Sort = sortOrders[0] // Assuming only support single sort order for simplicity
	}
	return nil
}
