// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package checker

import "time"

// Check is the interface which every check that needs to be processed here
// should implement.
type Check interface {
	GetID() string   // Returns the ID.
	GetName() string // Returns the name.

	GetInterval() time.Duration // Returns the interval after which check is run.
	GetTimeout() time.Duration  // Returns the timeout.

	GetInput() Component      // Returns the input.
	GetOutput() Component     // Returns the output.
	GetTarget() Component     // Returns the target.
	GetPayloads() []Component // Returns the payloads.
}

// Component is the Type Value component for check components like Input,
// Output, Target etc.
type Component interface {
	GetType() string  // Returns the type.
	GetValue() string // Returns the value.
}

// Metric is anything that tells if the check is successful or not. It also
// tells if the check timed out, the start time and duration of check.
type Metric interface {
	GetCheckID() string
	GetCheckName() string

	IsSuccessful() bool
	IsTimeout() bool

	GetStartTime() time.Time
	GetDuration() time.Duration
}
