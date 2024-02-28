package vat

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vkuzmenkova/currency-rates/internal/providers"
	"github.com/vkuzmenkova/currency-rates/models"
)

const (
	HOST = "https://api.vatcomply.com"
)

type VATProvider struct {
	Host string
}

func NewVATProvider() providers.RatesProvider {
	return &VATProvider{
		Host: HOST,
	}
}

func (p *VATProvider) GetRate(base string, currency_code string) (*models.CurrencyRate, error) {
	resp, err := http.Get(fmt.Sprintf("%s/rates?base=%s", p.Host, base))
	if err != nil {
		return &models.CurrencyRate{}, fmt.Errorf("GET %s/rates?base=%s %d: %w", p.Host, base, resp.StatusCode, err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &models.CurrencyRate{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	type currency struct {
		Date  string
		Base  string
		Rates map[string]float64
	}

	var c currency
	err = json.Unmarshal(respBody, &c)
	if err != nil {
		return &models.CurrencyRate{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	rate := models.CurrencyRate{
		UpdatedAt: time.Now().String(),
		Base:      base,
		Currency:  currency_code,
		Value:     c.Rates[currency_code],
	}
	log.Println(rate)
	time.Sleep(1 * time.Minute)

	return &rate, nil
}
