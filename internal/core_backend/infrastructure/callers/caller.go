package callers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"backend-service/internal/core_backend/common/helper"
	"backend-service/internal/core_backend/common/logger"
)

type Caller struct{}

func NewCaller() *Caller {
	return &Caller{}
}

func (c *Caller) CallGetMethod(host, path string, queryParams map[string]string) (io.ReadCloser, error) {
	requestURL := fmt.Sprintf("%s/%s?%s", host, path, helper.MapToString(queryParams))
	res, err := http.Get(requestURL)
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to call GET request to %s", host))
	}

	return res.Body, err
}

func CreateRequest(url, method string, data interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func SendRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
