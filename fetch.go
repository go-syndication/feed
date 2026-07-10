package feed

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Fetch GETs url using client (or [http.DefaultClient] if nil) and parses the
// body into a normalized [Feed].
func Fetch(ctx context.Context, client *http.Client, url string) (*Feed, error) {
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("feed: fetch %s: unexpected status %s", url, resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}
