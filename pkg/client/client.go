package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL *url.URL
}

func New(baseURL string) (*Client, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL: parsedBaseURL,
	}, nil
}

func MakeRequest[Response any](c *Client, method, path string, payload any) (int, Response, error) {
	var zero Response
	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(payload); err != nil {
		return 0, zero, err
	}

	request, err := http.NewRequest(method, c.baseURL.JoinPath(path).String(), &body)
	if err != nil {
		return 0, zero, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, zero, err
	}

	var responseBody Response
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return 0, zero, err
	}

	return response.StatusCode, responseBody, nil
}
