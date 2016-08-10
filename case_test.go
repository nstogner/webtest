package webtest

import (
	"log"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
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
