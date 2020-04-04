package probes

import (
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"time"

	"github.com/gorilla/websocket"
)

// WSProber probes a websocket url. This connects with the target and sends messages.
// Next message is sent to the target once the previous response is received. Messages
// received are in the same order as they are sent.
type WSProber struct {
	timeout  time.Duration
	url      *neturl.URL
	headers  map[string]string
	messages []string
}

// GetURL returns the URL of the `WSProber`.
func (pr *WSProber) GetURL() string {
	return pr.url.String()
}

// SetURL sets the URL of the `WSProber`.
func (pr *WSProber) SetURL(rawurl string) error {
	u, err := neturl.Parse(rawurl)
	if err != nil {
		return err
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return fmt.Errorf("url scheme should be ws(s) and not %s", u.Scheme)
	}
	pr.url = u
	return nil
}

func (pr *WSProber) deadline(startTime time.Time) time.Time {
	return startTime.Add(pr.timeout)
}

// NewWSProber returns a websocket prober. Requires a valid "ws" scheme url. Messages and headers can be nil.
func NewWSProber(url string, messages []string, headers map[string]string, timeout time.Duration) (*WSProber, error) {
	prober := &WSProber{}
	if err := prober.SetURL(url); err != nil {
		return nil, err
	}

	if messages == nil {
		prober.messages = []string{}
	} else {
		prober.messages = messages
	}

	if headers == nil {
		prober.headers = map[string]string{}
	} else {
		prober.headers = headers
	}

	prober.timeout = timeout
	return prober, nil
}

// Probe executes the prober. It connects with the target websocket URL, sends and receives messages.
func (pr *WSProber) Probe() (*WSProbeResults, error) {
	startTime := time.Now()
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: pr.timeout,
	}

	conn, resp, err := dialer.Dial(pr.GetURL(), parseHeaders(pr.headers))
	if err != nil {
		if errIsTimeout(err) {
			return &WSProbeResults{
				Messages:  []string{},
				StartTime: startTime,
				Duration:  pr.timeout,
				Timeout:   true,
			}, nil
		}
		return nil, err
	}

	defer conn.Close()      //nolint:errcheck
	defer resp.Body.Close() //nolint:errcheck

	response := &WSProbeResults{
		Messages:   []string{},
		StartTime:  startTime,
		StatusCode: resp.StatusCode,
		Body:       resp.Body,
		Headers:    resp.Header,
	}

	if deadlineErr := conn.SetWriteDeadline(pr.deadline(startTime)); deadlineErr != nil {
		return nil, deadlineErr
	}

	if deadlineErr := conn.SetReadDeadline(pr.deadline(startTime)); deadlineErr != nil {
		return nil, deadlineErr
	}

	for _, msg := range pr.messages {
		if err = conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			if errIsTimeout(err) {
				response.Duration = pr.timeout
				response.Timeout = true
				return response, nil
			}
			return nil, err
		}

		_, recv, err := conn.ReadMessage()
		if err != nil {
			if errIsTimeout(err) {
				response.Duration = pr.timeout
				response.Timeout = true
				return response, nil
			}
			return nil, err
		}
		response.Messages = append(response.Messages, string(recv))
	}

	response.Timeout = false
	response.Duration = time.Since(startTime)
	return response, nil
}

// WSProbeResults contain the results of the probe. This consists of whether the probe was timeout,
// the time and duration, response code, body, headers and the messages received (if any) in response.
type WSProbeResults struct {
	Timeout bool

	StartTime time.Time
	Duration  time.Duration

	Messages []string

	StatusCode int
	Body       io.ReadCloser
	Headers    http.Header
}

func parseHeaders(orig map[string]string) http.Header {
	headers := make(http.Header)
	for key, val := range orig {
		headers.Add(key, val)
	}
	return headers
}
