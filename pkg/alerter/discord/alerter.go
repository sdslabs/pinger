// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sdslabs/pinger/pkg/alerter"
	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"

	"github.com/sirupsen/logrus"
)

// serviceName is the name of the service used to send the alert.
const serviceName = "discord"

func init() {
	alerter.Register(serviceName, func() alerter.Alerter { return new(Alerter) })
}

// reqBody is the JSON request body format for slack webhook request.
type reqBody struct {
	Text string `json:"content"`
}

// Alerter sends an alert for test status.
type Alerter struct {
	log *logrus.Logger
}

// Provision initializes required fields for a's execution.
func (a *Alerter) Provision(ctx *appcontext.Context, _ alerter.Provider) error {
	a.log = ctx.Logger()
	return nil
}

// Alert sends the notification on discord.
func (a *Alerter) Alert(ctx context.Context, metrics []checker.Metric, amap map[string]alerter.Alert) error {
	for i := range metrics {
		metric := metrics[i]
		alt, ok := amap[metric.GetCheckID()]
		if !ok {
			continue
		}

		if err := a.alert(ctx, metric, alt); err != nil {
			a.log.Errorf("check %s: %v", metric.GetCheckID(), err)
			continue
		}
	}

	return nil
}

// alert sends an individual notification.
func (a *Alerter) alert(ctx context.Context, metric checker.Metric, alt alerter.Alert) error {
	var msg string
	if metric.IsSuccessful() {
		msg = fmt.Sprintf("%s is back up", metric.GetCheckName())
	} else {
		msg = fmt.Sprintf("%s is down", metric.GetCheckName())
		if metric.IsTimeout() {
			msg = fmt.Sprintf("%s: timeout", metric.GetCheckName())
		}
	}

	body, err := json.Marshal(reqBody{Text: msg})
	if err != nil {
		return fmt.Errorf("unexpected error while marshaling: %v", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		alt.GetTarget(),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}

	defer resp.Body.Close() // nolint:errcheck

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response: %v", err)
	}
	if buf.String() != "ok" {
		return fmt.Errorf("not-ok response returned from slack")
	}

	return nil
}

// Interface guard.
var _ alerter.Alerter = (*Alerter)(nil)
