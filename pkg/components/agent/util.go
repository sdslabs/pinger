package agent

// boolean represents a bool variable in string format.
type boolean string

// boolean constants.
const (
	boolTrue  boolean = "true"
	boolFalse boolean = "false"
	boolNil   boolean = "nil"
)

// newBoolean returns a string bool.
func newBoolean(from bool) boolean {
	if from {
		return boolTrue
	}
	return boolFalse
}

// True tells if the boolean string is true.
func (b boolean) True() bool { return b == boolTrue }

// False tells if the boolean string is false.
func (b boolean) False() bool { return b == boolFalse }

// True tells if the boolean string is nil.
func (b boolean) Nil() bool { return b == boolNil }
