package currencyrates

import "fmt"

type NoUUIDFoundError struct {
	Message string
}

func (e NoUUIDFoundError) Error() string {
	return fmt.Sprintf("UUID not found: %s", e.Message)
}

type NoValueFoundError struct {
	Currency string
}

func (e NoValueFoundError) Error() string {
	return fmt.Sprintf("No info about: %s. Update currency rate first.", e.Currency)
}
