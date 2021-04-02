package app

import "fmt"

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("Error %v: %v", r.StatusCode, r.Err)
}
