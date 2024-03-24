package tests

import (
	"api-shapes/pkg/client"
	"api-shapes/transport"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const url = "http://localhost:8080"

func TestUserAPI_Create(t *testing.T) {
	bts, _ := json.Marshal(transport.UserReq{Name: "mock"})
	req, err := http.NewRequest(http.MethodPost, url+"/users", bytes.NewReader(bts))
	assert.Nil(t, err)

	status, _, err := client.Request(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)
}

func seedUser(t *testing.T) transport.UserRes {
	bts, _ := json.Marshal(transport.UserReq{Name: "mock"})
	req, err := http.NewRequest(http.MethodPost, url+"/users", bytes.NewReader(bts))
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
