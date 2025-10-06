package api

import "context"

type Client interface {
	Get(ctx context.Context, url string, headers map[string]string) (*Response, error)
	Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error)
	Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*Response, error)
	Patch(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error)
}

type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}
