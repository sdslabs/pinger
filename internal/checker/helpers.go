// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package checker

import (
	"fmt"
	"net"
	"regexp"
)

// addressRegexPattern is the regex pattern for AddressRegex.
const addressRegexPattern = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`

// MaxMessageSize is the maximum size of the message that can be sent or
// received.
const MaxMessageSize = 65536

// AddressRegex is the regex pattern that is needs to be matched by an
// address.
var AddressRegex *regexp.Regexp

func init() {
	AddressRegex = regexp.MustCompile(addressRegexPattern)
}

// ErrIsTimeout returns true for a timeout error.
func ErrIsTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

type validateComponentFunc = func(string) error

// ValidateComponent takes a map of types and validation of their values
// function and validates the component.
func ValidateComponent(c Component, m map[string]validateComponentFunc) error {
	typ := c.GetType()
	val := c.GetValue()

	validateFunc, ok := m[typ]
	if !ok {
		return fmt.Errorf("invalid component type: %s", typ)
	}

	return validateFunc(val)
}
