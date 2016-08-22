package webtest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

// Timeout is the timeout on a given request.
var Timeout = 1 * time.Minute

// Fataler is defined to match a testing.B or testing.T.
type Fataler interface {
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

// Response wraps a *http.Response.
type Response struct {
	*http.Response
	Fail func(error)
}

// DoReq runs a *http.Request.
func DoReq(t Fataler, req *http.Request) *Response {
	resp, err := (&http.Client{Timeout: Timeout}).Do(req)
	if err != nil {
		t.Fatalf("unexpected error issuing HTTP request: %s", err)
	}
	return &Response{resp, func(err error) {
		t.Fatalf("%s %s ... failed: %s", req.Method, req.URL.Path, err)
	}}
}

// Do constructs a *http.Request and runs it.
func Do(t Fataler, ts *httptest.Server, method, path string, body []byte) *Response {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", ts.URL, path), bytes.NewReader(body))
	if err != nil {
		t.Fatalf("unable to generate HTTP request: %s", err)
	}
	return DoReq(t, req)
}

// E is short for expectation.
type E struct {
	Status int
	JSON   interface{}
	XML    interface{}
}

// Expect enforces expectations (E). It also closes the Request.Body.
func (r *Response) Expect(exp E) *Response {
	defer r.Body.Close()

	// Catch bad test-cases.
	if exp.JSON != nil && exp.XML != nil {
		r.Fail(errors.New("BAD TEST CASE: UNABLE TO UNMARSHAL RESPONSE TO XML AND JSON"))
	}

	// Run all expectations.

	if exp.Status > 0 {
		if got := r.StatusCode; got != exp.Status {
			b, _ := ioutil.ReadAll(r.Body)
			r.Fail(fmt.Errorf("expected status %v, got: %v with body:\n%s", exp.Status, got, b))
		}
	}

	if exp.JSON != nil {
		if err := json.NewDecoder(r.Body).Decode(exp.JSON); err != nil {
			r.Fail(fmt.Errorf("unable to decode response as JSON: %s", err))
		}
	}
	if exp.XML != nil {
		if err := xml.NewDecoder(r.Body).Decode(exp.XML); err != nil {
			r.Fail(fmt.Errorf("unable to decode response as JSON: %s", err))
		}
	}

	return r
}

// Context runs a function for further validation and will spit out some request
// context-specific error message if that function returns an error.
func (r *Response) Context(f func() error) {
	if err := f(); err != nil {
		r.Fail(err)
	}
}
