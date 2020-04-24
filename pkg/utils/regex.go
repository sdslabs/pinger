package utils

import "regexp"

// Compiled regex patterns.
var (
	RegexEmail *regexp.Regexp
)

func init() {
	// compile patterns during init of the function since "MustCompile" panics
	// so if there is any invalid regex it will panic before the anything even
	// starts it's execution

	RegexEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
}
