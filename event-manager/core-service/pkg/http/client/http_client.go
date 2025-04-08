package client

import (
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	cli     *http.Client
	baseUrl string
}

func NewClient(baseUrl string) Client {
	return Client{
		cli: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:      10,
				IdleConnTimeout:   30 * time.Second,
				DisableKeepAlives: false,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
			},
		},
		baseUrl: baseUrl,
	}
}

func (c *Client) Invoke(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseUrl+endpoint, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "new request")
	}
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "do")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "read all")
	}

	return bodyBytes, nil
}
