package api

import (
	"github.com/htquangg/a-wasm/config"

	"github.com/go-resty/resty/v2"
)

var RedactedHeaders = []string{" X-Request-Id"}

type Client struct {
	HTTPClient *resty.Client
}

type ClientOptions struct {
	Debug bool
}

func NewClient(p *ClientOptions) *Client {
	httpClient := resty.New()
	httpClient.
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")

	if p.Debug {
		httpClient.EnableTrace()

		httpClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			logRequest(req)
			return nil
		})

		httpClient.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			logResponse(resp)
			return nil
		})
	}

	httpClient.SetBaseURL(config.AWASM_URL)

	return &Client{
		HTTPClient: httpClient,
	}
}
