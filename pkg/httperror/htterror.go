package httperror

import "net/http"

type Error struct {
	err  string
	code int
}

func (e Error) Error() string {
	return e.err
}

func (e Error) Code() int {
	if e.code == 0 {
		return http.StatusInternalServerError
	}
	return e.code
}

// New is the way to create a new httperror
func New(code int, msg string) error {
	return Error{
		code: code,
		err:  msg,
	}
}
