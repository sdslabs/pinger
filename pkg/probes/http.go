package probes

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var defaultTransport = http.DefaultTransport.(*http.Transport)

// NewHTTPProber creates Prober that will skip TLS verification while probing.
func NewHTTPProber() HTTPProber {
	tlsConfig := &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	return NewWithTLSConfig(tlsConfig)
}

// NewWithTLSConfig creates a Prober with provided TLS config.
func NewWithTLSConfig(config *tls.Config) HTTPProber {
	transport := setOldTransportDefaults(
		&http.Transport{
			TLSClientConfig:   config,
			DisableKeepAlives: true,
		})
	return HTTPProber{transport}
}

func setOldTransportDefaults(t *http.Transport) *http.Transport {
	if t.DialContext == nil && t.Dial == nil { //nolint:staticcheck
		t.DialContext = defaultTransport.DialContext
	}

	if t.TLSHandshakeTimeout == 0 {
		t.TLSHandshakeTimeout = defaultTransport.TLSHandshakeTimeout
	}
	return t
}

// Parse the response obtained from making the reqeust using the prober,
// it takes a few fields of the response and return it in a concise way to
// be digested later.
func parseResponse(resp *http.Response, duration time.Duration, startTime time.Time) *HTTPProbeResult {
	return &HTTPProbeResult{
		Timeout: false,

		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       resp.Body,

		StartTime: startTime,
		Duration:  duration,
	}
}

// HTTPProber is the prober for HTTP request checks.
type HTTPProber struct {
	transport *http.Transport
}

// GetProbe executes `Probe` method for a "GET" HTTP request.
func (pr HTTPProber) GetProbe(
	url string,
	headers, payload map[string]string,
	timeout time.Duration) (*HTTPProbeResult, error) {
	return pr.Probe("GET", url, headers, payload, timeout)
}

// PostProbe executes `Probe` method for a "POST" HTTP request.
func (pr HTTPProber) PostProbe(
	url string,
	headers, payload map[string]string,
	timeout time.Duration) (*HTTPProbeResult, error) {
	return pr.Probe("POST", url, headers, payload, timeout)
}

// Probe is the main entrypoint for doing a HTTP Probe using the package.
// The method specify the type of HTTP request we are trying to make and the other
// parameters are populated accordingly in the request.
func (pr *HTTPProber) Probe(
	method, url string,
	headers, payload map[string]string,
	timeout time.Duration) (*HTTPProbeResult, error) {
	client := &http.Client{
		Timeout:   timeout,
		Transport: pr.transport,
	}

	var err error

	if headers == nil {
		headers = map[string]string{}
	}
	if payload == nil {
		payload = map[string]string{}
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling JSON payload")
	}

	var req *http.Request
	if method == "POST" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payloadJSON))
	} else {
		req, err = http.NewRequest(method, url, nil)
		q := req.URL.Query()
		for key, val := range payload {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if err != nil {
		return nil, fmt.Errorf("error while preparing request: %s", err)
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	req.Header.Set("Content-Type", "application/json")
	startTime := time.Now()
	resp, err := client.Do(req)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		// Request errored due to a timeout.
		// Send a curated response in this case.

		return &HTTPProbeResult{Timeout: true}, nil
	} else if err != nil {
		return nil, fmt.Errorf("error while making request: %s", err)
	}

	duration := time.Since(startTime)

	return parseResponse(resp, duration, startTime), nil
}

// HTTPProbeResult is the result of HTTP check probe.
type HTTPProbeResult struct {
	Timeout bool

	StatusCode int
	Body       io.ReadCloser
	Headers    http.Header

	// Time at which the probe execution started.
	StartTime time.Time
	// Duration that the probe lasted for.
	Duration time.Duration
}
