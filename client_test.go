package form3

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewDefaultClient(t *testing.T) {
	client := NewDefaultClient(nil)

	require.NotNil(t, client)
	require.NotNil(t, client.httpClient)

	require.Equal(t, "http://localhost:8080/", client.BaseURL.String())
	require.Equal(t, "form3-client/v1", client.UserAgent)
}

func Test_whenRequestWithoutBodyThenContentTypeNotSet(t *testing.T) {
	client := testClient("")

	request, err := client.newRequest(http.MethodGet, "/", nil)
	require.Nil(t, err)
	require.Equal(t, request.Header.Get("Content-Type"), "")
}

func Test_whenRequestWithBodyThenContentTypeIsSet(t *testing.T) {
	client := testClient("")

	request, err := client.newRequest(http.MethodGet, "/", "")
	require.Nil(t, err)
	require.Equal(t, request.Header.Get("Content-Type"), contentType)
}

func Test_whenRequestWithInvalidBodyThenReturnError(t *testing.T) {
	client := testClient("")

	req, err := client.newRequest(http.MethodGet, "/", make(chan int))
	require.NotNil(t, err, err.Error())
	require.Nil(t, req)
}

func Test_whenCallingExistingServiceThenFillBody(t *testing.T) {
	server := startServer()
	client := testClient(server.URL)
	defer server.Close()

	req, err := client.newRequest(http.MethodGet, "/", nil)
	require.Nil(t, err)
	require.NotNil(t, req)

	resBody := &struct {
		Status string `json:"status"`
	}{}

	_, err = client.do(req, resBody)
	require.Nil(t, err)
	require.Equal(t, resBody.Status, "success")
}

func Test_whenCallingNotExistingServiceThenReturnError(t *testing.T) {
	client := testClient("http://not-existing-service")

	req, err := client.newRequest(http.MethodGet, "/", nil)
	require.Nil(t, err)
	require.NotNil(t, req)

	resBody := &struct {
		Status string `json:"status"`
	}{}

	_, err = client.do(req, resBody)
	require.NotNil(t, err)
}

func testClient(addr string) *Client {
	testURL, _ := url.Parse(addr)
	client := &Client{
		httpClient: http.DefaultClient,
		BaseURL:    testURL,
		UserAgent:  "",
	}

	return client
}

func startServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":"success"}`)
	}))
}
