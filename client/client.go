package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/pkg/errors"
)

type httpClient interface {
	Do(context.Context, *http.Request) (*http.Response, []byte, error)
}

type simpleHttpClient struct {
	client  *http.Client
	timeout time.Duration
	verbose bool
	logger  *log.Logger // stdout when nil, verbose should be true
}

func (c *simpleHttpClient) Do(ctx context.Context, r *http.Request) (*http.Response, []byte, error) {
	if c.verbose {
		if c.logger == nil {
			debug.Info("[URL]", r.URL.String())
			debug.Info("[METHOD]", r.Method)
			for n, h := range r.Header {
				debug.Info(fmt.Sprintf("[REQUEST] %s: %s", n, h))
			}
		} else {
			c.logger.Println("[URL]", r.URL.String())
			c.logger.Println("[METHOD]", r.Method)
			for n, h := range r.Header {
				c.logger.Println(fmt.Sprintf("[REQUEST] %s: %s", n, h))
			}
		}
	}

	var lctx context.Context
	var lcancel context.CancelFunc
	req := r
	if c.timeout > 0 {
		lctx, lcancel = context.WithTimeout(ctx, c.timeout)
		req = r.WithContext(lctx)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return resp, nil, errors.Wrap(err, "do failed")
	}

	defer resp.Body.Close()
	if lcancel != nil {
		defer lcancel()
	}

	if c.verbose {
		if c.logger == nil {
			for n, h := range resp.Header {
				debug.Info(fmt.Sprintf("[RESPONSE] %s: %s", n, h))
			}

			debug.Info("[STATUS]", resp.Status)
		} else {
			for n, h := range resp.Header {
				c.logger.Println(fmt.Sprintf("[RESPONSE] %s: %s", n, h))
			}

			c.logger.Println("[STATUS]", resp.Status)
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, body, errors.Wrap(err, "readall failed")
	}

	return resp, body, nil
}

// NewSimpleHttpClient creates an instance of simpleHttpClient.
//
// If client is nil, http.DefaultClient is used. If logger is nil, standard log is used when
// verbose is set to true.
func NewSimpleHttpClient(timeout time.Duration, verbose bool, logger *log.Logger) *simpleHttpClient {
	return &simpleHttpClient{
		client:  http.DefaultClient,
		timeout: timeout,
		verbose: verbose,
		logger:  logger,
	}
}
