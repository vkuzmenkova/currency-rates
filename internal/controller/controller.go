package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	currencyrates "github.com/vkuzmenkova/currency-rates/internal/domain/service"
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

// UpdateRate godoc
// @Summary      Initiates updating rate
// @Description  Initiates updating rate
// @Tags         rate
// @Produce      json
// @Param	     code    path     string  USD  "currency base, f.e. EUR"  Format(string)
// @Param	     base    query     string  USD  "currency base, default=USD"  Format(string)
// @Success      200  {object} models.CurrencyUpdateUUID
// @Failure      400  {object} UnavailableCurrencyError
// @Failure      400  {object} BaseAndCodeAreEqual
// @Failure      500  {object} BaseAndCodeAreEqual
// @Router       /rates/{code}/update [get]
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
		http.Error(resp, BaseAndCodeAreEqualError{Message: fmt.Sprintf("%s:%s", base, code)}.Error(), http.StatusBadRequest)
		return
	}

	uuidUpdate, err := c.service.UpdateRate(ctx, base, code)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(
		models.CurrencyUpdateUUID{
			Base:     base,
			Currency: code,
			UUID:     uuidUpdate,
		},
	)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetLastRate godoc
// @Summary      Gets rate from the last update
// @Description  Gets rate from the last update
// @Tags         rate
// @Produce      json
// @Param	     code    path     string  USD  "currency base, f.e. EUR"  Format(string)
// @Param	     base    query     string  USD  "currency base, default=USD"  Format(string)
// @Success      200  {object} models.CurrencyRate
// @Failure      400  {object} UnavailableCurrencyError
// @Failure      400  {object} BaseAndCodeAreEqual
// @Failure      500  {object} BaseAndCodeAreEqual
// @Router       /rates/{code} [get]
func (c *Controller) GetLastRate(resp http.ResponseWriter, req *http.Request) {
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
		http.Error(resp, BaseAndCodeAreEqualError{Message: fmt.Sprintf("%s:%s", base, code)}.Error(), http.StatusBadRequest)
		return
	}

	var cr models.CurrencyRate

	cr, err = c.service.GetLastRate(ctx, base, code)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(cr)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetRateByUUID godoc
// @Summary      Gets rate by UUID
// @Description  Gets rate by UUID
// @Tags         rate
// @Produce      json
// @Param	     uuid    query     string  _  "uuid of update"  Format(string)
// @Success      200  {object} models.CurrencyRate
// @Failure      400  {object} InvalidUUIDError
// @Failure      400  {object} currencyrates.NoUUIDFoundError
// @Failure      500  {object} string
// @Router       /rates [get]
func (c *Controller) GetRateByUUID(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	UUID, err := ExtractUUID(req)
	if err != nil {
		http.Error(resp, InvalidUUIDError{}.Error(), http.StatusBadRequest)
		return
	}

	var cr models.CurrencyRate

	// Ищем uuid в базе
	cr, err = c.service.GetRateByUUID(ctx, UUID)
	// Если не найден uuid
	var noUUIDErr *currencyrates.NoUUIDFoundError
	if errors.As(err, &noUUIDErr) {
		http.Error(resp, noUUIDErr.Error(), http.StatusBadRequest)
		return
	}
	// Иные ошибки
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(cr)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = resp.Write(result)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}
