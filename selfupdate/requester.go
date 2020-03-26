package selfupdate

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Requester interface allows developers to customize the method in which
// requests are made to retrieve the version and binary
type Requester interface {
	Fetch(ctx context.Context, url string) (io.ReadCloser, error)
}

// HTTPRequester is the normal requester that is used and does an HTTP
// to the url location requested to retrieve the specified data.
type HTTPRequester struct {
}

// Fetch will return an HTTP request to the specified url and return
// the body of the result. An error will occur for a non 200 status code.
func (httpRequester *HTTPRequester) Fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad http status from %s: %v", url, resp.Status)
	}

	return resp.Body, nil
}

// mockRequester used for some mock testing to ensure the requester contract
// works as specified.
type mockRequester struct {
	currentIndex int
	fetches      []func(context.Context, string) (io.ReadCloser, error)
}

func (mr *mockRequester) handleRequest(requestHandler func(context.Context, string) (io.ReadCloser, error)) {
	if mr.fetches == nil {
		mr.fetches = []func(context.Context, string) (io.ReadCloser, error){}
	}
	mr.fetches = append(mr.fetches, requestHandler)
}

func (mr *mockRequester) Fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	if len(mr.fetches) <= mr.currentIndex {
		return nil, fmt.Errorf("No for currentIndex %d to mock", mr.currentIndex)
	}
	current := mr.fetches[mr.currentIndex]
	mr.currentIndex++

	return current(ctx, url)
}
