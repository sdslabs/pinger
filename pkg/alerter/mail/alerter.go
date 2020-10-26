// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package mail

import (
	"context"
	"fmt"
	"time"

	gomail "gopkg.in/mail.v2"

	"github.com/sdslabs/pinger/pkg/alerter"
	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"

	"github.com/sirupsen/logrus"
)

// serviceName is the name of the service used to send the alert.
const serviceName = "mail"

// defaultTimeout is the time after which a mail being sent is considered failed.
const defaultTimeout = time.Minute

func init() {
	alerter.Register(serviceName, func() alerter.Alerter { return new(Alerter) })
}

// senderDetails stores the config for sending E-mail.
type senderDetails struct {
	Host   string
	Port   uint16
	User   string
	Secret string
}

// Alerter sends an alert for test status.
type Alerter struct {
	log    *logrus.Logger
	sender senderDetails
}

// Provision initializes required fields for a's execution.
func (a *Alerter) Provision(ctx *appcontext.Context, prov alerter.Provider) error {
	a.log = ctx.Logger()
	a.sender = senderDetails{prov.GetHost(), prov.GetPort(), prov.GetUser(), prov.GetSecret()}
	return nil
}

// Alert sends the notification on mail.
func (a *Alerter) Alert(ctx context.Context, metrics []checker.Metric, amap map[string]alerter.Alert) error {
	for i := range metrics {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

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
	var (
		thisCtx = ctx
		cancel  func()
	)
	if _, ok := thisCtx.Deadline(); !ok {
		thisCtx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}

	var msg string
	if metric.IsSuccessful() {
		msg = fmt.Sprintf("%s is back up", metric.GetCheckName())
	} else {
		msg = fmt.Sprintf("%s is down", metric.GetCheckName())
		if metric.IsTimeout() {
			msg = fmt.Sprintf("%s: timeout", metric.GetCheckName())
		}
	}

	to := alt.GetTarget()

	m := gomail.NewMessage()
	m.SetHeader("From", a.sender.User)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Pinger Alert: "+msg)
	m.SetBody("text/plain", msg)

	d := gomail.NewDialer(a.sender.Host, int(a.sender.Port), a.sender.User, a.sender.Secret)

	errChan := make(chan error)
	go func(dialer *gomail.Dialer, message *gomail.Message, stream chan<- error) {
		if err := d.DialAndSend(m); err != nil {
			stream <- fmt.Errorf("could not send request: %v", err)
		}
	}(d, m, errChan)

	select {
	case <-thisCtx.Done():
		return thisCtx.Err()
	case err := <-errChan:
		return err
	}
}

// Interface guard.
var _ alerter.Alerter = (*Alerter)(nil)
