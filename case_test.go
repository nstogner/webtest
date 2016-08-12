package webtest

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type mockT struct {
	realT *testing.T
}

func (m mockT) Fatalf(form string, stuff ...interface{}) {
	m.ShouldActuallyFail(fmt.Sprintf(form, stuff))
}

func (m mockT) Fatal(msg ...interface{}) {
	m.ShouldActuallyFail(fmt.Sprint(msg))
}

func (m mockT) ShouldActuallyFail(errMsg string) {
	if !strings.Contains(errMsg, "THIS SHOULD FAIL") {
		m.realT.Fatal("actual test should not have failed here:", errMsg)
	}
}

func TestRun(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users-json":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"id":"abc"}`)
		case "/users-xml":
			w.Header().Set("Content-Type", "application/xml")
			fmt.Fprint(w, `<User><Id>abc</Id></User>`)
		case "/users-with-no-content-type":
			fmt.Fprint(w, `{"id":"abc"}`)
		case "/users-with-empty-resp":
			fmt.Fprint(w, `{"id":"abc"}`)
		default:
			t.Fatal("This point should never have been reached")
		}
	})

	type user struct {
		XMLName xml.Name `xml:"User"`
		ID      string   `json:"id" xml:"Id"`
	}

	cases := []TestCase{
		{
			Name:           "retrieve users",
			Method:         "GET",
			URL:            "/users-json",
			ResponseEntity: &user{},
			Validate: func(e interface{}) error {
				u := e.(*user)
				if u.ID != "abc" {
					return errors.New("expected field user.id == 'abc'")
				}
				return nil
			},
		},
		{
			Name:           "retrieve users",
			Method:         "GET",
			URL:            "/users-xml",
			ResponseEntity: &user{},
			Validate: func(e interface{}) error {
				u := e.(*user)
				if u.ID != "abc" {
					return errors.New("expected field user.id == 'abc'")
				}
				return nil
			},
		},
		{
			Name:   "retrieve users - assume no response",
			Method: "GET",
			URL:    "/users-with-no-content-type",
		},
		{
			Name:           "retrieve users - THIS SHOULD FAIL",
			Method:         "GET",
			URL:            "/users-with-no-content-type",
			ResponseEntity: &user{},
			Validate: func(e interface{}) error {
				return nil
			},
		},
	}

	mt := mockT{t}

	Run(mt, handler, cases)
}
