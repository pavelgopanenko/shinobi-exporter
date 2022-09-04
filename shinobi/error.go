package shinobi

import (
	"errors"
	"fmt"
)

var ErrInconsistentResponseFormat = errors.New("inconsistent response format")

type ErrorUnexpecterAPIResponseStatus int

func (e ErrorUnexpecterAPIResponseStatus) Error() string {
	return fmt.Sprintf("unexpected api response status %d", int(e))
}
