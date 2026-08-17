package requests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Pool interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

type pool struct {
	semaphore    chan struct{}
	client       http.Client
	poolSize     int
	openRequests int
}

func NewPool(concurrentRequests int) Pool {
	return &pool{
		semaphore: make(chan struct{}, concurrentRequests),
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		poolSize: concurrentRequests,
	}
}

func (p *pool) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Will succeed when there is room in the channel. So blocking when the
	// request pool is already saturated.
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("Context closed")
	case p.semaphore <- struct{}{}:
	}

	if req.Header == nil {
		req.Header = http.Header{}
	}
	req.Header.Add("User-Agent", "welp/v1")
	response, err := p.client.Do(req)

	if err != nil {
		<-p.semaphore
		return response, err
	}

	// Wrap the response body reader so when it closes the semaphore is
	// released.
	response.Body = &onCloseReader{
		ReadCloser: response.Body,
		onClose: func() {
			<-p.semaphore
		},
	}

	return response, nil
}

type onCloseReader struct {
	io.ReadCloser
	onClose func()
}

func (o *onCloseReader) Close() error {
	defer o.onClose()
	return o.ReadCloser.Close()
}
