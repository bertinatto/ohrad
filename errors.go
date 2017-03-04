package ohrad

import "fmt"

type ErrorTimeout struct {
	where string
}

func (e *ErrorTimeout) Error() string {
	return fmt.Sprintf("Timed out: %s", e.where)
}
