package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	currencyrates "github.com/vkuzmenkova/currency-rates/internal/domain/currency_rates"
	"github.com/vkuzmenkova/currency-rates/models"
)

type Controller struct {
	service *currencyrates.CurrenciesService
}

func NewController(ctx context.Context) (*Controller, error) {
	service, err := currencyrates.NewCurrenciesService(ctx)
	if err != nil {
		return &Controller{}, fmt.Errorf("NewCurrenciesService: %w", err)
	}

	return &Controller{
		service: service,
	}, nil
}

func ExtractBase(req *http.Request, base string) string {
	reqBase := strings.ToUpper(req.URL.Query().Get("base"))
	if reqBase == "" {
		reqBase = base
	}

	return strings.ToUpper(reqBase)
}

func ExtractUUID(req *http.Request) (uuid.UUID, error) {
	uuidString := req.URL.Query().Get("uuid")
	if uuidString == "" {
		return uuid.Nil, nil
	}
	uuidValue, err := uuid.Parse(uuidString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s invaild", uuidValue)
	}

	return uuidValue, nil
}

func ExtractCode(req *http.Request) (string, error) {
	vars := mux.Vars(req)
	var code string

	if value, ok := vars["code"]; !ok {
		return "", fmt.Errorf("unable to get currency code from vars")
	} else {
		code = strings.ToUpper(value)
	}

	return code, nil

}

// Ручка обновления  ставки
func (c *Controller) UpdateRate(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	base := ExtractBase(req, c.service.CurrencyList.Base)
	code, err := ExtractCode(req)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if !c.service.CurrencyList.IsCurrencyAvailable(base) || !c.service.CurrencyList.IsCurrencyAvailable(code) {
		http.Error(resp, UnavailableCurrencyError{CurrencyList: c.service.CurrencyList.GetCurrencyListUpper()}.Error(), http.StatusBadRequest)
		return
	}
	if base == code {
		http.Error(resp, BaseAndCodeAreEqual{Message: fmt.Sprintf("%s:%s", base, code)}.Error(), http.StatusBadRequest)
		return
	}

	uuid, err := c.service.UpdateRateV2(ctx, base, code)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(
		models.CurrencyUpdateUUID{
			Base:     base,
			Currency: code,
			UUID:     uuid,
		},
	)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Write(result)
}

// Ручка получения ставки
func (c *Controller) GetRate(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	base := ExtractBase(req, c.service.CurrencyList.Base)
	code, err := ExtractCode(req)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if !c.service.CurrencyList.IsCurrencyAvailable(base) || !c.service.CurrencyList.IsCurrencyAvailable(code) {
		http.Error(resp, UnavailableCurrencyError{CurrencyList: c.service.CurrencyList.GetCurrencyListUpper()}.Error(), http.StatusBadRequest)
		return
	}
	if c.service.CurrencyList.IsCurrencyCodeEqualsBase(base, code) {
		http.Error(resp, BaseAndCodeAreEqual{Message: fmt.Sprintf("%s:%s", base, code)}.Error(), http.StatusBadRequest)
		return
	}

	UUID, err := ExtractUUID(req)
	if err != nil {
		http.Error(resp, InvalidUUIDError{}.Error(), http.StatusBadRequest)
		return
	}

	var cr models.CurrencyRate

	if UUID == uuid.Nil {
		// Если uuid не указан
		cr, err = c.service.GetLastRate(ctx, base, code)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Если запрос еще обновляется
		_, err := c.service.KV.Get(ctx, fmt.Sprintf("%s_%s", code, base)).Result()
		if err == nil {
			http.Error(resp, fmt.Sprintf("Update request with uuid=%s is in process", UUID.String()), http.StatusAccepted)
			return
		}
		// Ищем uuid в базе
		cr, err = c.service.GetRateByUUID(ctx, UUID)
		// Если не найден uuid
		if noUUIDErr, ok := err.(*currencyrates.NoUUIDFoundError); ok {
			http.Error(resp, noUUIDErr.Error(), http.StatusBadRequest)
			return
		}
		// Иные ошибки
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	result, err := json.Marshal(cr)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Write(result)
}
