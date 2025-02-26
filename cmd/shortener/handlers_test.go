package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestCreateShortURL(t *testing.T) {
	type expected struct {
		code     int
		response string
	}
	type request struct {
		url    string
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
				code:     201,
				response: `^http:\/\/localhost:8080\/[A-Z]+$`,
			},
			request: request{
				url:    "https://ya.ru/",
				method: http.MethodPost,
			},
		},
		{
			name: "negative test1",
			expected: expected{
				code:     405,
				response: ``,
				//contentType: "text/plain",
			},
			request: request{
				url:    "",
				method: http.MethodGet,
			},
		},
		{
			name: "negative test2",
			expected: expected{
				code:     400,
				response: ``,
			},
			request: request{
				url:    "https://ya.ru/",
				method: http.MethodPost,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.request.method, "/", bytes.NewBufferString(test.request.url))
			rr := httptest.NewRecorder()
			createShortURL(rr, req)

			if status := rr.Code; status != test.expected.code {
				t.Errorf("handler returned wrong status code: got %v expected %v",
					status, test.expected.code)
			}

			re := regexp.MustCompile(test.expected.response)
			if body := rr.Body.String(); !re.MatchString(body) {
				t.Errorf("handler returned wrong response: got %v expected %v",
					body, test.expected.response)
			}
		})
	}
}

func TestGetURL(t *testing.T) {

	testId := "SVHZQO"
	_, err := URLList[testId]
	if !err {
		URLList[testId] = "https://ya.ru/"
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
				id:     testId,
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
				id:     testId,
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
