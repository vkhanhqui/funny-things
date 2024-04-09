package rest

import (
	"api-shapes/pkg/client"
	"api-shapes/tests"
	"api-shapes/transport"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAPI_List(t *testing.T) {
	createRes := seedUser(t)
	req, err := http.NewRequest(http.MethodGet, tests.URL+"/rest/users", nil)
	assert.Nil(t, err)

	status, bts, err := client.Request(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	var res transport.ListRes
	err = json.Unmarshal(bts, &res)
	assert.Nil(t, err)
	assert.Equal(t, createRes, res.Users[len(res.Users)-1])
}
