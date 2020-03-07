package form3

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	version          = "v1"
	defaultUserAgent = "form3-client/" + version
	defaultBaseURL   = "http://localhost:8080/"
	contentType      = "application/vnd.api+json"
)

type Client struct {
	httpClient *http.Client

	BaseURL   *url.URL
	UserAgent string
}

func NewDefaultClient(httpClient *http.Client) (client *Client) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	client = &Client{httpClient: httpClient, BaseURL: baseURL, UserAgent: defaultUserAgent}

	return client
}

func (c *Client) newRequest(method, path string, body interface{}) (req *http.Request, err error) {
	reqURL := c.BaseURL.ResolveReference(&url.URL{Path: path})
	buf := bytes.Buffer{}

	if body != nil {
		encoder := json.NewEncoder(&buf)
		if err = encoder.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err = http.NewRequest(method, reqURL.String(), &buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", defaultUserAgent)

	return req, nil
}

func (c *Client) do(req *http.Request, respType interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(respType)
	return resp, err
}
