// Package handlers contains handlers for various app server routes.
package handlers

import (
	"errors"
	"strconv"

	"github.com/sdslabs/status/pkg/utils"
)

var errInvalidIDParam = errors.New("invalid ID parameter")

// paramIsEmail returns true if the parameter is email and false if parameter
// is an integer ID. Returns a non-nil error for other cases. This also returns
// the final parameter as an interface{} that can be later type asserted to get
// either ID or email.
func paramIsEmail(parameter string) (isEmail bool, param interface{}, err error) {
	// Check if it's a valid integer ID
	id, e := strconv.ParseUint(parameter, 10, 0)
	if e == nil {
		return false, uint(id), nil
	}

	// Check if it's a valid email address
	if utils.RegexEmail.MatchString(parameter) {
		return true, parameter, nil
	}

	return false, nil, errInvalidIDParam
}
