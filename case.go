package webtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// TestCase is an http test case.
type TestCase struct {
	Name           string
	Method         string
	URL            string
	BodyString     string
	BodyReader     io.Reader
	ResponseEntity interface{}
	Validate       func(interface{}) error
}

// Run runs a slice of test cases.
func Run(t Fataler, handler http.Handler, cases []TestCase) {
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		var req *http.Request
		var err error
		if tc.BodyReader != nil {
			req, err = http.NewRequest(tc.Method, tc.URL, tc.BodyReader)
		} else {
			req, err = http.NewRequest(tc.Method, tc.URL, bytes.NewReader([]byte(tc.BodyString)))
		}

		if err != nil {
			t.Fatal(failStr(tc, "unable to create HTTP request", err))
		}

		handler.ServeHTTP(rec, req)

		ct := rec.Header().Get("Content-Type")
		if strings.Contains(ct, "json") {
			if err := json.Unmarshal(rec.Body.Bytes(), tc.ResponseEntity); err != nil {
				t.Fatal(failStr(tc, "unable to decode response as JSON", err))
			}
			if err := tc.Validate(tc.ResponseEntity); err != nil {
				t.Fatal(failStr(tc, "validation failed", err))
			}
		} else if strings.Contains(ct, "xml") {
			if err := json.Unmarshal(rec.Body.Bytes(), tc.ResponseEntity); err != nil {
				t.Fatal(failStr(tc, "unable to decode response as JSON", err))
			}
			if err := tc.Validate(tc.ResponseEntity); err != nil {
				t.Fatal(failStr(tc, "validation failed", err))
			}
		} else if tc.Validate != nil {
			t.Fatal(failStr(tc, "unable to validate response", errors.New("missing Content-Type header")))
		}
	}
}

func failStr(tc TestCase, msg string, err error) string {
	return fmt.Sprintf("failed test case %q: %s: %s", tc.Name, msg, err)
}
