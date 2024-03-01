package currencyrates

import "fmt"

type NoUUIDFoundError struct {
	Message string
}

func (e NoUUIDFoundError) Error() string {
	return fmt.Sprintf("Updating is in progress or UUID is not exist.")
}

type NoValueFoundError struct {
	Currency string
}

func (e NoValueFoundError) Error() string {
	return fmt.Sprintf("No info about: %s. Update currency rate first.", e.Currency)
}
