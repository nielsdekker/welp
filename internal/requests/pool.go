package requests

import (
	"io"
	"net/http"
	"time"
)

type Pool interface {
	Do(req *http.Request) (*http.Response, error)
	GetPoolSize() int
}

type pool struct {
	semaphore chan struct{}
	client    http.Client
	poolSize  int
}

func NewPool(concurrentRequests int) Pool {
	return &pool{
		semaphore: make(chan struct{}, concurrentRequests),
		client:    http.Client{Timeout: 5 * time.Second},
		poolSize:  concurrentRequests,
	}
}

func (p *pool) GetPoolSize() int {
	return p.poolSize
}

func (p *pool) Do(req *http.Request) (*http.Response, error) {
	// Will succeed when there is room in the channel. So blocking when the
	// request pool is already saturated.
	p.semaphore <- struct{}{}

	response, err := p.client.Do(req)

	if err != nil {
		<-p.semaphore
		return response, err
	}

	// Wrap the response body reader so when it closes the semaphore is
	// released.
	response.Body = &onCloseReader{
		ReadCloser: response.Body,
		onClose:    func() { <-p.semaphore },
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
