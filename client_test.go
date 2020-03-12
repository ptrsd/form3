package form3

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type assertion struct {
	expected interface{}
	actual   interface{}
	name     string
}

type assertions []assertion

func Test_NewDefaultClient(t *testing.T) {
	client := NewDefaultClient(nil)

	notNils := assertions{
		{actual: client, name: "Client"},
		{actual: client.httpClient, name: "Client.httpClient"},
	}
	assertNotNil(t, notNils)

	equals := assertions{
		{actual: client.BaseURL.String(), expected: defaultBaseURL, name: "Client.BaseURL"},
		{actual: client.UserAgent, expected: defaultUserAgent, name: "Client.UserAgent"},
	}
	assertEquals(t, equals)
}

func Test_whenRequestWithoutBodyThenContentTypeIsNotSet(t *testing.T) {
	client := testClient("")

	request, err := client.newRequest(http.MethodGet, &url.URL{Path: "/"}, nil)
	if err != nil {
		t.Errorf("error while creating new request, %s", err.Error())
	}

	assertEquals(t, assertions{
		{actual: request.Header.Get("Content-Type"), expected: "", name: "Client.ContentType"},
	})
}

func Test_whenRequestWithBodyThenContentTypeIsSet(t *testing.T) {
	client := testClient("")

	request, err := client.newRequest(http.MethodPost, &url.URL{Path: "/"}, "")
	if err != nil {
		t.Errorf("error while creating new request, %s", err.Error())
	}

	assertEquals(t, assertions{
		{actual: request.Header.Get("Content-Type"), expected: contentType, name: "Client.ContentType"},
	})
}

func Test_whenRequestWithInvalidBodyThenReturnError(t *testing.T) {
	client := testClient("")

	_, err := client.newRequest(http.MethodGet, &url.URL{Path: "/"}, make(chan int))
	assertNotNil(t, assertions{
		{actual: err, name: "Client.InvalidBodyError"},
	})

	assertEquals(t, assertions{
		{actual: err.Error(), expected: "json: unsupported type: chan int", name: "Client.InvalidBodyErrorMessage"},
	})
}

func Test_whenCallingExistingServiceThenReturnBody(t *testing.T) {
	server := startServer()
	defer server.Close()

	client := testClient(server.URL)

	req, err := client.newRequest(http.MethodGet, &url.URL{Path: "/"}, nil)
	if err != nil {
		t.Errorf("error while creating new request, %s", err.Error())
	}

	resBody := &struct {
		Status string `json:"status"`
	}{}

	err = client.do(req, resBody)
	if err != nil {
		t.Errorf("error while calling service, %s", err.Error())
	}

	assertEquals(t, assertions{
		{actual: resBody.Status, expected: "success", name: "Client.ResponseBody"},
	})
}

func Test_whenCallingNotExistingServiceThenReturnError(t *testing.T) {
	client := testClient("http://not-existing-service")

	req, err := client.newRequest(http.MethodGet, &url.URL{Path: "/"}, nil)
	if err != nil {
		t.Errorf("error while creating new request, %s", err.Error())
	}

	resBody := &struct {
		Status string `json:"status"`
	}{}

	err = client.do(req, resBody)
	assertNotNil(t, assertions{
		{actual: err, name: "Client.CallingNotExistingService"},
	})
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

func assertNotNil(t *testing.T, notNils assertions) {
	for _, assertion := range notNils {
		if assertion.actual == nil {
			t.Errorf("%s: Expected not to be nil", assertion.name)
		}
	}
}

func assertEquals(t *testing.T, equals assertions) {
	for _, assertion := range equals {
		if !reflect.DeepEqual(assertion.expected, assertion.actual) {
			t.Errorf("%s:\nExpected: %#v\n  Actual: %#v", assertion.name, assertion.expected, assertion.actual)
		}
	}
}

func assertNotEmpty(t *testing.T, notEmpty []assertion) {
	for _, assertion := range notEmpty {
		if assertion.actual == reflect.Zero(reflect.TypeOf(assertion.actual)).Interface() {
			t.Errorf("%s: Expected not to be empty", assertion.name)
		}
	}
}
