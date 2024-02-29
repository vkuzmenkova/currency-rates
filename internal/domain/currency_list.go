package domain

import (
	"strings"
)

type CurrencyList struct {
	Base                string
	AvailableCurrencies map[string]uint8
}

func (cl *CurrencyList) GetCurrencyList() string {
	var c []string
	for code := range cl.AvailableCurrencies {
		c = append(c, code, strings.ToLower(code))
	}

	return strings.Join(c, ",")
}

func (cl *CurrencyList) GetCurrencyListUpper() string {
	var c []string
	for code := range cl.AvailableCurrencies {
		c = append(c, code)
	}

	return strings.Join(c, ",")
}

func (cl *CurrencyList) GetCodeByValue(num uint8) string {
	for code, value := range cl.AvailableCurrencies {
		if value == num {
			return code
		}
	}

	return ""
}

func (cl *CurrencyList) IsCurrencyCodeEqualsBase(base string, code string) bool {
	return base == code
}

func (cl *CurrencyList) IsCurrencyAvailable(code string) bool {
	_, ok := cl.AvailableCurrencies[strings.ToUpper(code)]
	return ok
}

func (cl *CurrencyList) GetValueByCode(code string) uint8 {
	value := cl.AvailableCurrencies[strings.ToUpper(code)]
	return value
}

func NewCurrencyList() *CurrencyList {
	base := "USD"
	currencies := map[string]uint8{
		"USD": 1,
		"EUR": 2,
		"MXN": 3,
	}

	return &CurrencyList{
		Base:                base,
		AvailableCurrencies: currencies,
	}
}
