// Package checker implements the checkers and probers for the various
// protocols.
//
// For any checker it is required that the type implements the Checker
// interface.
//
// Further, a checker can be used to either just validate the check input
// using the validate method:
//
//	if err := checker.Validate(checkConfig); err != nil {
//		// handle error
//	}
//
// NewControllerOpts can be used to create a controller options from the
// check config which can be paired with a controller which executes checker
// at regular intervals of time.
//
// This package also contains some helpers which are common to use among
// checkers, such as, regex for checking if the address is valid or not or
// if the err is timeout.
package checker
