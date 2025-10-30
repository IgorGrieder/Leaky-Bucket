package application

import "fmt"

type NoTokensError struct {
}

func (e *NoTokensError) Error() string {
	return fmt.Sprintf("error: no tokens available")
}
