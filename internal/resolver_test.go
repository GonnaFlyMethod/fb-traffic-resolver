package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	pathToBuildFolder = "./testdata/build"
	prefixOfAPIInPath = "/api"
)

func TestTrafficResolver_Resolve(t *testing.T) {
	apiServerStartControl := make(chan *url.URL)
	apiServerCloseControl := make(chan struct{})

	go startTestAPIServer(apiServerStartControl, apiServerCloseControl)
	addressOfTestAPI := <-apiServerStartControl

	tr := NewTrafficResolver(addressOfTestAPI, prefixOfAPIInPath, pathToBuildFolder)
	resolverHandler := http.HandlerFunc(tr.Resolve)

	testCases := []struct {
		name, requestPath, expectedContentType string
		expectedStatusCode                     int
		expectedResponseBody                   []byte
	}{
		{
			name:        "should get index html",
			requestPath: "/",

			expectedContentType:  "text/html; charset=utf-8",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: loadFile(t, "index.html"),
		},
		{
			name:        "should get index html",
			requestPath: "/non_existent_url",

			expectedContentType:  "text/html; charset=utf-8",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: loadFile(t, "index.html"),
		},
		{
			name:        "should get favicon.ico",
			requestPath: "/favicon.ico",

			expectedContentType:  "image/vnd.microsoft.icon",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: loadFile(t, "favicon.ico"),
		},
		{
			name:        "should get json response from api",
			requestPath: fmt.Sprintf("%s/locations", prefixOfAPIInPath),

			expectedContentType:  "application/json",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: []byte(`[{"location1": "OK"}, {"location2": "OK"}, {"location3": "OK"}]`),
		},
		{
			name:        "should get '404 page not found' response from api",
			requestPath: fmt.Sprintf("%s/non_existent_api_url", prefixOfAPIInPath),

			expectedContentType:  "text/plain; charset=utf-8",
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: []byte("404 page not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tc.requestPath, http.NoBody)
			responseRecorder := httptest.NewRecorder()

			resolverHandler.ServeHTTP(responseRecorder, request)

			actualResponseBody, err := io.ReadAll(responseRecorder.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedContentType, responseRecorder.Header().Get("Content-Type"))
			assert.Equal(t, tc.expectedStatusCode, responseRecorder.Code)
			assert.True(t, bytes.Equal(tc.expectedResponseBody, actualResponseBody))
		})
	}

	apiServerCloseControl <- struct{}{}
	<-apiServerCloseControl
}

func testAPIServerRouter(res http.ResponseWriter, req *http.Request) {
	var responseBody []byte

	switch req.URL.Path {
	case fmt.Sprintf("%s/locations", prefixOfAPIInPath):
		res.Header().Set("Content-Type", "application/json")
		responseBody = []byte(`[{"location1": "OK"}, {"location2": "OK"}, {"location3": "OK"}]`)
	default:
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusNotFound)
		responseBody = []byte("404 page not found")
	}

	if _, err := res.Write(responseBody); err != nil {
		panic("can't write response body as a response in test API server")
	}
}

func startTestAPIServer(apiServerStartControl chan *url.URL, apiServerCloseControl chan struct{}) {
	hf := http.HandlerFunc(testAPIServerRouter)

	s := httptest.NewServer(hf)

	parsedURL, err := url.Parse(s.URL)
	if err != nil {
		s.Close()
		panic("error while parsing address of test api server")
	}

	apiServerStartControl <- parsedURL

	<-apiServerCloseControl
	s.Close()

	fmt.Println("test api server was successfully closed")

	apiServerCloseControl <- struct{}{}
}

func loadFile(t *testing.T, fileName string) []byte {
	filePath := fmt.Sprintf("%s/%s", pathToBuildFolder, fileName)

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	return data
}
