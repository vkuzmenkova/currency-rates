package controller

import "fmt"

type BaseAndCodeAreEqual struct {
	Message string
}

func (e BaseAndCodeAreEqual) Error() string {
	return fmt.Sprintf("Base and currency code are the same: %s", e.Message)
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
	return "Ivalid UUID"
}
