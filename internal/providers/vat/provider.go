package vat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/avast/retry-go"
	"github.com/vkuzmenkova/currency-rates/internal/providers"
	"github.com/vkuzmenkova/currency-rates/models"
)

const (
	HOST           = "https://api.vatcomply.com"
	RequestTimeout = 2 * time.Second
	RETRY          = 5
	RetryDelay     = 2 * time.Second
)

type VATProvider struct {
	Host       string
	Client     http.Client
	Retry      uint
	RetryDelay time.Duration
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0  --name=HTTP
type HTTP interface {
	Get(url string) (resp *http.Response, err error)
}

func NewVATProvider() providers.RatesProvider {
	return &VATProvider{
		Host: HOST,
		Client: http.Client{
			Timeout: RequestTimeout,
		},
		Retry:      RETRY,
		RetryDelay: RetryDelay,
	}
}

func (p *VATProvider) GetRate(ctx context.Context, base string, currencyCode string) (*models.CurrencyRate, error) {
	var respBody []byte

	err := retry.Do(
		func() error {
			resp, err := p.Client.Get(fmt.Sprintf("%s/rates?base=%s", p.Host, base))
			if err != nil {
				return fmt.Errorf("client.Get: %w", err)
			}

			respBody, err = io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("io.ReadAll: %w", err)
			}

			err = resp.Body.Close()
			if err != nil {
				return fmt.Errorf("body.Close: %w", err)
			}
			return nil
		},
		retry.Context(ctx),
		retry.Delay(p.RetryDelay),
		retry.Attempts(p.Retry),
		retry.OnRetry(func(n uint, err error) {
			log.Err(err).Msgf("retry %d: %s\n", n, err.Error())
		}),
	)
	if err != nil {
		return &models.CurrencyRate{}, fmt.Errorf("retry.Do: %w", err)
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
		Currency:  currencyCode,
		Value:     c.Rates[currencyCode],
	}
	//time.Sleep(1 * time.Second) // Imitates long response

	return &rate, nil
}
