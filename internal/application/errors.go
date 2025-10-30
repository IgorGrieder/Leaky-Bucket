package application

type NoTokensError struct {
}

func (e *NoTokensError) Error() string {
	return "error: no tokens available"
}
