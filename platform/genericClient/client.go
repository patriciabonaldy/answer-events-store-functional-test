package genericClient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	// ErrNotFound error
	ErrNotFound = errors.New("resource not found")
	// ErrURLIsEmpty error
	ErrURLIsEmpty = errors.New("request does not have url")
	// ErrBodyIsEmpty error
	ErrBodyIsEmpty = errors.New("request does not have body")
)

// Header represents Header in the request.
type Header struct {
	Key   string
	Value string
}

// Client handler different http methods
type Client interface {
	Delete(ctx context.Context, url string, headers ...Header) error
	Get(ctx context.Context, url string) (resp *http.Response, err error)
	Post(ctx context.Context, url string, data []byte, headers ...Header) (resp *http.Response, err error)
}

// Client defines the communication client.
type client struct {
	httpClient *http.Client
}

// New create a new client
func New() Client {
	return &client{
		httpClient: &http.Client{Transport: &http.Transport{}},
	}
}

func (c *client) Delete(ctx context.Context, url string, headers ...Header) error {
	req, err := c.withHeader(ctx, http.MethodDelete, url, nil, headers)
	if err != nil {
		return err
	}

	_, err = c.do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := c.withHeader(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err = c.do(req)
	if err != nil {
		fmt.Println("paso x aqui")
		return resp, err
	}

	return resp, nil
}

func (c *client) Post(ctx context.Context, url string, data []byte, headers ...Header) (resp *http.Response, err error) {
	fmt.Println("request:", string(data))
	body := bytes.NewReader(data)
	req, err := c.withHeader(ctx, http.MethodPost, url, body, headers)
	if err != nil {
		return nil, err
	}

	if body == nil {
		return nil, ErrBodyIsEmpty
	}

	resp, err = c.do(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *client) withHeader(ctx context.Context, method, url string, body io.Reader, headers []Header) (*http.Request, error) {
	if url == "" {
		return nil, ErrURLIsEmpty
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to make request [%s:%s]: %w", req.Method, req.URL.String(), err)
	}

	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	return req, nil
}

func (c *client) do(req *http.Request) (resp *http.Response, err error) {
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed doing request [%s:%s]: %s", req.Method, req.URL.String(), err.Error())
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, ErrNotFound
	default:
		return resp, fmt.Errorf("failed to do request, %d status code received", resp.StatusCode)
	}
}
