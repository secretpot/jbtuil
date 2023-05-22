package errs

import (
	"fmt"
)

type cerror struct {
	code    int
	message string
}

func (e cerror) Error() string {
	return fmt.Sprintf("code: %d, msg: %v", e.code, e.message)
}

func New(code int, message string) error {
	return cerror{
		code,
		message,
	}
}
func Err(message string) error {
	return cerror{
		-1,
		message,
	}
}

func ERQ(err interface{}) { // Exit Request
	if err != nil {
		if _, ok := err.(error); ok {
			panic(fmt.Sprintf("%v", err))
		}
	}
}
