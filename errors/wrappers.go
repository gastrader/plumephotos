package errors

import "errors"

//give us access to existing functions in std library
var(
	As = errors.As
	Is = errors.Is
)

