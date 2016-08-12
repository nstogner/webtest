package webtest

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"foo":"bar"}`)
	})

	type foo struct {
		Foo string `json:"foo"`
	}
	cases := []TestCase{
		{
			Name:           "retrieve users",
			Method:         "GET",
			URL:            "/users",
			ResponseEntity: &foo{},
			Validate: func(e interface{}) error {
				f := e.(*foo)
				if f.Foo != "bar" {
					return errors.New("expected field foo == 'bar'")
				}
				return nil
			},
		},
	}

	Run(t, handler, cases)
}
