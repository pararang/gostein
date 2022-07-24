package gostein

import "fmt"

type ErrNot2XX struct {
	StatusCode int
	Message    string
}

func (e ErrNot2XX) Error() string {
	return fmt.Sprintf("http status code %d: %s", e.StatusCode, e.Message)
}

type ErrDecodeJSON struct {
	Err error
}

func (e ErrDecodeJSON) Error() string {
	return fmt.Sprintf("decode json error: %s", e.Err)
}

type ErrParamCantBeNil struct {
	ParamName string
}

func (e ErrParamCantBeNil) Error() string {
	return fmt.Sprintf("param %s should have value", e.ParamName)
}
