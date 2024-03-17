package client

import (
	"io"
	"net/http"
	"time"
)

func Request(req *http.Request) (int, []byte, error) {
	client := http.DefaultClient
	client.Timeout = 30 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer resp.Body.Close()

	bts, err := io.ReadAll(resp.Body)
	return resp.StatusCode, bts, err
}
