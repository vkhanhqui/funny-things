package soap

import (
	"api-shapes/pkg/client"
	"api-shapes/tests"
	"api-shapes/transport"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAPI_Retrieve(t *testing.T) {
	createRes := seedUser(t)
	req, err := http.NewRequest(http.MethodGet, tests.URL+"/soap/users/"+createRes.ID, nil)
	assert.Nil(t, err)

	status, bts, err := client.Request(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	var res transport.UserRes
	err = xml.Unmarshal(bts, &res)
	assert.Nil(t, err)
	assert.Equal(t, createRes, res)
}
