package feed

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchSuccessNilClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, jsonSample)
	}))
	defer srv.Close()

	// nil client exercises the http.DefaultClient branch.
	f, err := Fetch(context.Background(), nil, srv.URL)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if f.Title != "JSON Example" {
		t.Errorf("Title = %q", f.Title)
	}
}

func TestFetchSuccessCustomClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, rssSample)
	}))
	defer srv.Close()

	f, err := Fetch(context.Background(), srv.Client(), srv.URL)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if f.Title != "Example Channel" {
		t.Errorf("Title = %q", f.Title)
	}
}

func TestFetchBadRequest(t *testing.T) {
	// A control character in the URL makes NewRequestWithContext fail.
	_, err := Fetch(context.Background(), http.DefaultClient, "http://\x7f/bad")
	if err == nil {
		t.Fatal("expected request-construction error")
	}
}

type errRoundTripper struct{}

func (errRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("transport boom")
}

func TestFetchTransportError(t *testing.T) {
	client := &http.Client{Transport: errRoundTripper{}}
	_, err := Fetch(context.Background(), client, "http://example.invalid/feed")
	if err == nil {
		t.Fatal("expected transport error")
	}
}

func TestFetchNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := Fetch(context.Background(), srv.Client(), srv.URL)
	if err == nil || !strings.Contains(err.Error(), "unexpected status") {
		t.Fatalf("want status error, got %v", err)
	}
}

type errBody struct{}

func (errBody) Read([]byte) (int, error) { return 0, errors.New("read boom") }
func (errBody) Close() error             { return nil }

type bodyErrRoundTripper struct{}

func (bodyErrRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       errBody{},
		Header:     make(http.Header),
	}, nil
}

func TestFetchBodyReadError(t *testing.T) {
	client := &http.Client{Transport: bodyErrRoundTripper{}}
	_, err := Fetch(context.Background(), client, "http://example.invalid/feed")
	if err == nil || !strings.Contains(err.Error(), "read boom") {
		t.Fatalf("want read error, got %v", err)
	}
}
