// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sdslabs/status/internal/checker"
)

const (
	// contentTypeHeader is the key for request header content type.
	contentTypeHeader = "Content-Type"

	// jsonContentType is the content type value for JSON.
	jsonContentType = "application/json"
)

// Prober sends an HTTP request to the request url and accepts the response
// from the server.
type Prober struct {
	method   string
	url      string
	headers  map[string]string
	payloads map[string]interface{}
	timeout  time.Duration
}

// NewProber creates a new HTTP Prober.
func NewProber(
	method, targetURL string,
	headers map[string]string,
	payloads map[string]interface{},
	timeout time.Duration,
) (*Prober, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout should be > 0")
	}

	if headers == nil {
		headers = map[string]string{}
	}
	if payloads == nil {
		payloads = map[string]interface{}{}
	}

	return &Prober{
		method:   method,
		url:      targetURL,
		headers:  headers,
		payloads: payloads,
		timeout:  timeout,
	}, nil
}

// Probe sends an HTTP request to the url.
func (p *Prober) Probe(ctx context.Context) (*ProbeResult, error) {
	startTime := time.Now()

	baseCtx := ctx
	if ctx == nil {
		baseCtx = context.Background()
	}

	probeCtx, cancel := context.WithTimeout(baseCtx, p.timeout)
	defer cancel()

	timeoutResult := &ProbeResult{
		Timeout:   true,
		StartTime: startTime,
		Duration:  p.timeout,
	}

	req, err := p.getRequest(probeCtx)
	if err != nil {
		if checker.ErrIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}

	cli := &http.Client{
		Timeout: p.timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // nolint:gosec
			},
		},
	}

	resp, err := cli.Do(req)
	if err != nil {
		if checker.ErrIsTimeout(err) {
			return timeoutResult, nil
		}

		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ProbeResult{
		Timeout:    false,
		StartTime:  startTime,
		Duration:   time.Since(startTime),
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       buf.String(),
	}, nil
}

// getRequest creates a new HTTP request from the context passed to it.
func (p *Prober) getRequest(ctx context.Context) (*http.Request, error) {
	var req *http.Request

	if p.method == http.MethodGet || p.method == "" {
		var err error

		req, err = http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
		if err != nil {
			return nil, fmt.Errorf("cannot create request: %w", err)
		}

		// Add payloads as GET queries to the request.
		q := req.URL.Query()
		for key, val := range p.payloads {
			strVal := fmt.Sprint(val)
			q.Add(key, strVal)
		}

		req.URL.RawQuery = q.Encode()
	} else {
		var body io.Reader

		contentType, ok := p.headers[contentTypeHeader]
		if ok && contentType == jsonContentType {
			pjson, err := json.Marshal(p.payloads)
			if err != nil {
				return nil, fmt.Errorf("cannot create JSON payload: %w", err)
			}

			body = bytes.NewBuffer(pjson)
		} else { // assume it's form data and let the header be
			form := url.Values{}

			for key, val := range p.payloads {
				strVal := fmt.Sprint(val)
				form.Add(key, strVal)
			}

			body = strings.NewReader(form.Encode())
		}

		var err error

		req, err = http.NewRequestWithContext(ctx, p.method, p.url, body)
		if err != nil {
			return nil, fmt.Errorf("cannot create request: %w", err)
		}
	}

	for key, val := range p.headers {
		req.Header.Add(key, val)
	}

	return req, nil
}

// ProbeResult is the result of HTTP probe.
type ProbeResult struct {
	Timeout    bool
	StatusCode int
	Body       string
	Headers    http.Header
	StartTime  time.Time
	Duration   time.Duration
}
