package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestCreateShortURL(t *testing.T) {
	ts := httptest.NewServer(routes())
	defer ts.Close()

	var testTable = []struct {
		url     string
		want    string
		status  int
		method  string
		longURL string
	}{
		{"/", `^http:\/\/localhost:8080\/[A-Z]+$`, http.StatusCreated, http.MethodPost, "https://ya.ru"},
		{"/", "", http.StatusBadRequest, http.MethodPost, "https://ya.ru"},
		{"/", "", http.StatusMethodNotAllowed, http.MethodGet, "https://ya.ru"},
	}
	for _, v := range testTable {
		resp, short := testRequest(t, ts, v.method, v.url)
		defer resp.Body.Close()
		assert.Equal(t, v.status, resp.StatusCode)
		re := regexp.MustCompile(v.want)
		if !re.MatchString(short) {
			t.Errorf("handler returned wrong response: got %v expected %v",
				short, v.want)
		}

	}
}

func TestGetURL(t *testing.T) {

	testID := "SVHZQO"
	_, err := URLList[testID]
	if !err {
		URLList[testID] = "https://ya.ru/"
	}

	type expected struct {
		code     int
		location string
	}
	type request struct {
		id     string
		method string
	}
	tests := []struct {
		name     string
		expected expected
		request  request
	}{
		{
			name: "positive test",
			expected: expected{
				code:     307,
				location: `https://ya.ru/`,
			},
			request: request{
				id:     testID,
				method: http.MethodGet,
			},
		},
		{
			name: "negative test1",
			expected: expected{
				code:     405,
				location: ``,
			},
			request: request{
				id:     testID,
				method: http.MethodPost,
			},
		},
		{
			name: "negative test2",
			expected: expected{
				code:     404,
				location: ``,
			},
			request: request{
				id:     "h123456",
				method: http.MethodGet,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(test.request.method, "/", nil)
			req.SetPathValue("id", test.request.id)
			rr := httptest.NewRecorder()
			getURL(rr, req)

			if status := rr.Code; status != test.expected.code {
				t.Errorf("handler returned wrong status code: got %v expected %v",
					status, test.expected.code)
			}

			if location := rr.Header().Get("Location"); location != test.expected.location {
				t.Errorf("handler returned wrong response: got %v expected %v",
					location, test.expected.location)
			}
		})
	}
}
