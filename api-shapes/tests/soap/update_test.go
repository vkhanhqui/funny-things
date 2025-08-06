package soap

import (
	"api-shapes/pkg/client"
	"api-shapes/tests"
	"api-shapes/transport"
	"bytes"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAPI_Update(t *testing.T) {
	u := seedUser(t)
	bts, _ := xml.Marshal(transport.UserReq{Name: "new"})
	req, err := http.NewRequest(http.MethodPut, tests.URL+"/soap/users/"+u.ID, bytes.NewReader(bts))
	assert.Nil(t, err)

	status, bts, err := client.Request(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	var res transport.UserRes
	err = xml.Unmarshal(bts, &res)

	assert.Nil(t, err)
	assert.Equal(t, u.ID, res.ID)
	assert.Equal(t, "new", res.Name)
}
