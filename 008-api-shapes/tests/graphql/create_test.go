package rest

import (
	"api-shapes/tests"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/graphql"

	"github.com/stretchr/testify/assert"
)

func TestUserAPI_Create(t *testing.T) {
	input := map[string]any{
		"query": `
			mutation Users {
				createUser(input: {name: "created name"}) {
					id
					name
				}
			}`,
		"variables": nil,
	}
	bts, _ := json.Marshal(input)
	res, err := http.Post(tests.URL+"/graphql", "application/json", bytes.NewReader(bts))
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func seedUser(t *testing.T) CreateUser {
	bts, _ := json.Marshal(map[string]any{
		"query": `
			mutation Users {
				createUser(input: {name: "created name"}) {
					id
					name
				}
			}`,
		"variables": nil,
	})
	res, err := http.Post(tests.URL+"/graphql", "application/json", bytes.NewReader(bts))
	assert.Nil(t, err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	var createRes graphql.Response
	err = json.Unmarshal(body, &createRes)
	assert.Nil(t, err)

	var cu Res
	err = json.Unmarshal(createRes.Data, &cu)
	assert.Nil(t, err)

	return cu.CreateUser
}

type Res struct {
	CreateUser CreateUser
}

type CreateUser struct {
	ID   string
	Name string
}
