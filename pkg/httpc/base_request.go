package httpc

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	DefaultDialTimeout    = 5 * time.Second
	DefaultRequestTimeout = 10 * time.Second
)

func request(req *http.Request) (*http.Response, error) {
	transport := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, DefaultDialTimeout)
		},
	}
	cli := http.Client{
		Transport: &transport,
		Timeout:   DefaultRequestTimeout,
	}
	return cli.Do(req)
}

func requestBody(req *http.Request) (int, []byte, error) {
	resp, err := request(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	var body []byte
	code := resp.StatusCode
	body, err = io.ReadAll(resp.Body)
	return code, body, err
}
