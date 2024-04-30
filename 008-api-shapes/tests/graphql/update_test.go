package rest

import (
	"api-shapes/tests"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAPI_Update(t *testing.T) {
	u := seedUser(t)

	input := map[string]any{
		"query": fmt.Sprintf(`
			mutation Users {
				updateUser(input: {id: "%s", name: "updated name"}) {
					id
					name
				}
			}`, u.ID),
		"variables": nil,
	}
	bts, _ := json.Marshal(input)
	res, err := http.Post(tests.URL+"/graphql", "application/json", bytes.NewReader(bts))
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
