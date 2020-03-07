package form3

import (
	"net/http"
	"net/url"
)

const (
	version          = "v1"
	defaultUserAgent = "form3-client/" + version
	defaultBaseURL   = "http://localhost:8080/"
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
