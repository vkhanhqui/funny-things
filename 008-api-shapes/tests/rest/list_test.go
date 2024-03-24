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
	req, err := http.NewRequest(http.MethodGet, tests.URL+"/v2/users", nil)
	assert.Nil(t, err)

	status, bts, err := client.Request(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	var res []transport.UserRes
	err = json.Unmarshal(bts, &res)
	assert.Nil(t, err)
	assert.Equal(t, createRes, res[len(res)-1])
}
