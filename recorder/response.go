package recorder

import "fmt"

type Response interface{}

type ResponseError error

type ResponseOK struct{}

func NewResponseErrorf(format string, a ...interface{}) Response {
	return fmt.Errorf(format, a...)
}

func NewResponseError(err error) Response {
	return err
}
