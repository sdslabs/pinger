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

// New creates Prober that will skip TLS verification while probing.
func NewHTTPProber() HttpProber {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	return NewWithTLSConfig(tlsConfig)
}

func setOldTransportDefaults(t *http.Transport) *http.Transport {
	if t.DialContext == nil && t.Dial == nil {
		t.DialContext = defaultTransport.DialContext
	}

	if t.TLSHandshakeTimeout == 0 {
		t.TLSHandshakeTimeout = defaultTransport.TLSHandshakeTimeout
	}
	return t
}

func NewWithTLSConfig(config *tls.Config) HttpProber {
	transport := setOldTransportDefaults(
		&http.Transport{
			TLSClientConfig:   config,
			DisableKeepAlives: true,
		})
	return HttpProber{transport}
}

type ProbeResult struct {
	Timeout bool

	StatusCode int
	Body       io.ReadCloser
	Headers    http.Header
}

type HttpProber struct {
	transport *http.Transport
}

func (pr HttpProber) GetProbe(url string, headers map[string]string, payload map[string]string, timeout time.Duration) (*ProbeResult, error) {
	return pr.Probe("GET", url, headers, payload, timeout)
}

func (pr HttpProber) PostProbe(url string, headers map[string]string, payload map[string]string, timeout time.Duration) (*ProbeResult, error) {
	return pr.Probe("POST", url, headers, payload, timeout)
}

// Main entrypoint for doing a HTTP Probe using the package.
// The method specify the type of HTTP request we are trying to make and the other
// parameters are populated accordingly in the request.
func (pr *HttpProber) Probe(method, url string, headers, payload map[string]string, timeout time.Duration) (*ProbeResult, error) {
	client := &http.Client{
		Timeout:   timeout,
		Transport: pr.transport,
	}

	var err error

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Error while marshalling JSON payload.")
	}

	var req *http.Request
	if method == "POST" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payloadJson))

	} else {
		req, err = http.NewRequest(method, url, nil)
		q := req.URL.Query()
		for key, val := range payload {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if err != nil {
		return nil, fmt.Errorf("Error while preparing request: %s", err)
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		// Request errored due to a timeout.
		// Send a curated response in this case.

		return &ProbeResult{Timeout: true}, nil

	} else if err != nil {
		return nil, fmt.Errorf("Error while making request: %s", err)
	}

	return parseResponse(resp), nil
}

// Parse the response obtained from making the reqeust using the prober,
// it takes a few fields of the response and return it in a concise way to
// be digested later.
func parseResponse(resp *http.Response) *ProbeResult {
	return &ProbeResult{
		Timeout: false,

		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       resp.Body,
	}
}
