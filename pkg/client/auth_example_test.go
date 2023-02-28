package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleTokenAuth() {
	ts := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer ts.Close()

	cl := New(Options{
		BaseURL: ts.URL,
		Auth:    &TokenAuth{"mytoken1234"},
	})

	if err := cl.Ping(context.Background()); err != nil {
		fmt.Printf("Pinging server failed: %v\n", err)
	} else {
		fmt.Println("Success!")
	}

	// Output: Success!
}
