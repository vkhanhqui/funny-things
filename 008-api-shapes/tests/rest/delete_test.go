package tests

import (
	"api-shapes/pkg/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAPI_Delete(t *testing.T) {
	u := seedUser(t)
	req, err := http.NewRequest(http.MethodDelete, url+"/users/"+u.ID, nil)
	assert.Nil(t, err)

	status, _, err := client.Request(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, status)
}
