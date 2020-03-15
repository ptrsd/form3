package form3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	version          = "v1"
	defaultUserAgent = "form3-client/" + version
	defaultBaseURL   = "http://localhost:8080"
	contentType      = "application/vnd.api+json"

	defaultPageSize = "100"
)

// Client provides interface for accessing Form3 server.
type Client struct {
	httpClient *http.Client

	// BaseURL is a base url for Form3 server. Default value: http://localhost:8080
	BaseURL        *url.URL
	UserAgent      string
	AccountService *AccountService
}

type ErrorMessage struct {
	ErrorMessage string `json:"error_message"`
}

//ListOptions defines page number and size of a page for a list operation.
type ListOptions struct {
	Page     int
	PageSize int
}

// NewDefaultClient creates a new instance of Form3 client. If httpClient is not provided then client will use http.DefaultClient.
func NewDefaultClient(httpClient *http.Client) (client *Client) {
	return NewClient(httpClient, defaultBaseURL)
}

// NewClient creates a new instance of Form3 client. If httpClient is not provided then client will use http.DefaultClient.
// BaseURL must be the valid URL value. If validation of BaseURL fails then client will return an error.
func NewClient(httpClient *http.Client, baseURL string) (client *Client) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	userBaseURL, _ := url.Parse(baseURL)
	client = &Client{httpClient: httpClient, BaseURL: userBaseURL, UserAgent: defaultUserAgent}

	client.AccountService = &AccountService{client}

	return client
}

func (c *Client) newRequest(method string, url *url.URL, body interface{}) (req *http.Request, err error) {
	reqURL := c.BaseURL.ResolveReference(url)
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

func (c *Client) do(req *http.Request, respType interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return err
	}

	if respType != nil {
		err = json.NewDecoder(resp.Body).Decode(respType)
	}

	return err
}

func checkError(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 204:
		return nil
	case 400, 401, 403, 404, 405, 406, 409, 429, 500, 502, 503, 504:
		errMsg := &ErrorMessage{}

		if err := json.NewDecoder(resp.Body).Decode(errMsg); err != nil {
			return err
		}

		if errMsg.ErrorMessage == "" {
			return fmt.Errorf(resp.Status)
		}
		return fmt.Errorf(errMsg.ErrorMessage)
	default:
		return fmt.Errorf("unknown status code %d", resp.StatusCode)
	}
}
