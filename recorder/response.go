package recorder

import "fmt"

// Response holds the response message from the recorder state machine.
// It can be one of the Response* types.
type Response interface{}

// ResponseError indicates an error response.
type ResponseError error

// ResponseOK indicates a succesful operation.
type ResponseOK struct{}

// ResponseLevel return the current known recording level.
type ResponseLevel MsgLevel

// NewResponseErrorf returns a new ResponseError
// according to the given format specifier.
func NewResponseErrorf(format string, a ...interface{}) ResponseError {
	return fmt.Errorf(format, a...)
}

// NewResponseError returns a new ResponseError
// based on the given error.
func NewResponseError(err error) ResponseError {
	return err
}
