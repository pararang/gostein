package gostein

import "fmt"

type ErrNot2XX struct {
	StatusCode int
}

func (e ErrNot2XX) Error() string {
	return fmt.Sprintf("http status code %d", e.StatusCode)
}

type ErrDecodeJSON struct {
	Err error
}

func (e ErrDecodeJSON) Error() string {
	return fmt.Sprintf("decode json error: %s", e.Err)
}
