package currencyrates

import "fmt"

type BaseAndCodeAreEqualError struct {
	Message string
}

func (e BaseAndCodeAreEqualError) Error() string {
	return "Base and currency codes must be different"
}

type UnavailableCurrencyError struct {
	CurrencyList string
}

func (e UnavailableCurrencyError) Error() string {
	return fmt.Sprintf("Currency unavailable. List of available currencies: %s", e.CurrencyList)
}

type InvalidUUIDError struct {
}

func (e InvalidUUIDError) Error() string {
	return "Invalid UUID"
}

type NoUUIDFoundError struct {
	Message string
}

func (e NoUUIDFoundError) Error() string {
	return "Updating is in progress or UUID does not exist."
}

type NoValueFoundError struct {
	Currency string
}

func (e NoValueFoundError) Error() string {
	return fmt.Sprintf("No info about: %s. Update currency rate first.", e.Currency)
}
