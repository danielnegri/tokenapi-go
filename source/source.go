package source

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/danielnegri/adheretech/errors"
	"github.com/danielnegri/adheretech/ledger"
	"github.com/danielnegri/adheretech/log"
	"github.com/danielnegri/adheretech/version"
	"github.com/go-resty/resty/v2"
)

const (
	DefaultRetry   = 5
	DefaultTimeout = 20 * time.Second
	DefaultURL     = "https://us-east4-at-devops-inhouse.cloudfunctions.net/be-interview-env-datasource-0e709f7f"
)

var userAgent = fmt.Sprintf("AdhereTech/Ledger Go HTTP Client %s", version.Version)

type Source interface {
	Check(ctx context.Context) error
	Generate(ctx context.Context, n int) ([]ledger.Token, error)
}

var _ Source = (*client)(nil)

type client struct {
	httpClient *resty.Client
	trace      bool
}

type Config struct {
	Retry   int
	Timeout time.Duration
	URL     *url.URL

	Trace bool
}

func New(cfg *Config) *client {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.Retry == 0 {
		cfg.Retry = DefaultRetry
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}

	if cfg.URL == nil {
		url, err := url.Parse(DefaultURL)
		if err != nil {
			panic(err)
		}

		cfg.URL = url
	}

	httpClient := resty.New().
		SetHostURL(cfg.URL.String()).
		SetLogger(log.Logger()).
		SetRetryCount(cfg.Retry).
		SetTimeout(cfg.Timeout)

	return &client{
		httpClient: httpClient,
		trace:      cfg.Trace,
	}
}

func (c *client) Check(ctx context.Context) error {
	log.Debugf("Checking Token source at %s", c.httpClient.HostURL)

	const op errors.Op = "source/client.Check"
	var apiErr interface{}
	resp, err := c.newRequest(ctx).SetError(&apiErr).SetQueryParam("size", "0").Post("/")
	if err != nil {
		log.Errorf("failed to request health check: %v", err)
		return errors.E(op, errors.Internal, err)
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("health check failed: status=%d, body=%s, error=%v", resp.StatusCode(), string(resp.Body()), apiErr)
		return errors.E(op, errors.Internal)
	}

	return nil
}

func (c *client) Generate(ctx context.Context, n int) ([]ledger.Token, error) {
	log.Debugf("Generating %d tokens", n)

	const op errors.Op = "source/client.Generate"
	if n <= 0 {
		return nil, errors.E(op, errors.Invalid, "number of tokens must be greater than zero")
	}

	var apiErr interface{}
	resp, err := c.newRequest(ctx).
		SetError(&apiErr).
		SetQueryParam("size", strconv.Itoa(n)).Post("/")
	if err != nil {
		log.Errorf("failed to request new tokens: %v", err)
		return nil, errors.E(op, errors.Internal, err)
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("request new tokens failed: status=%d, body=%s, error=%v", resp.StatusCode(), string(resp.Body()), apiErr)
		return nil, errors.E(op, errors.Internal)
	}

	tokens := make([]ledger.Token, 0)
	for _, line := range strings.Split(string(resp.Body()), "\n") {
		if line != "" {
			token := ledger.Token(line)
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

func (c *client) newRequest(ctx context.Context) *resty.Request {
	req := c.httpClient.R()
	if c.trace {
		req = req.EnableTrace()
	}

	return req
}
