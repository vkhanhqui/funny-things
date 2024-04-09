package rest

import (
	"api-shapes/pkg/client"
	"api-shapes/tests"
	"api-shapes/transport"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAPI_Create(t *testing.T) {
	bts, _ := json.Marshal(transport.UserReq{Name: "mock"})
	req, err := http.NewRequest(http.MethodPost, tests.URL+"/rest/users", bytes.NewReader(bts))
	assert.Nil(t, err)

	status, _, err := client.Request(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)
}

func seedUser(t *testing.T) transport.UserRes {
	bts, _ := json.Marshal(transport.UserReq{Name: "mock"})
	req, err := http.NewRequest(http.MethodPost, tests.URL+"/rest/users", bytes.NewReader(bts))
	if err != nil {
		t.Fatal(err)
	}

	_, bts, err = client.Request(req)
	if err != nil {
		t.Fatal(err)
	}

	var createRes transport.UserRes
	err = json.Unmarshal(bts, &createRes)
	if err != nil {
		t.Fatal(err)
	}
	return createRes
}
