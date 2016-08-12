package webtest

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
	})

	cases := []TestCase{
		{
			Name:   "retrieve users",
			Method: "GET",
			URL:    "/users",
			Validate: func(e interface{}) error {
				return nil
			},
		},
	}

	Run(t, handler, cases)
}
