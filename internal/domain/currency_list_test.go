package domain

import (
	"testing"
)

func TestCurrencyList_GetCodeByValue(t *testing.T) {
	type fields struct {
		Base                string
		AvailableCurrencies map[string]uint8
	}
	type args struct {
		num uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Get value from the list",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{"USD": 1, "EUR": 2, "MXN": 3}},
			args:   args{num: 3},
			want:   "MXN",
		},
		{
			name:   "Get value out of the list",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{"USD": 1, "EUR": 2, "MXN": 3}},
			args:   args{num: 0},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &CurrencyList{
				Base:                tt.fields.Base,
				AvailableCurrencies: tt.fields.AvailableCurrencies,
			}
			if got := cl.GetCodeByValue(tt.args.num); got != tt.want {
				t.Errorf("GetCodeByValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrencyList_GetCurrencyList(t *testing.T) {
	type fields struct {
		Base                string
		AvailableCurrencies map[string]uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Get list",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{"USD": 1}},
			want:   "USD,usd",
		},
		{
			name:   "Get empty list",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{}},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &CurrencyList{
				Base:                tt.fields.Base,
				AvailableCurrencies: tt.fields.AvailableCurrencies,
			}
			if got := cl.GetCurrencyList(); got != tt.want {
				t.Errorf("GetCurrencyList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrencyList_GetCurrencyListUpper(t *testing.T) {
	type fields struct {
		Base                string
		AvailableCurrencies map[string]uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Get list",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{"USD": 1}},
			want:   "USD",
		},
		{
			name:   "Get empty list",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{}},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &CurrencyList{
				Base:                tt.fields.Base,
				AvailableCurrencies: tt.fields.AvailableCurrencies,
			}
			if got := cl.GetCurrencyListUpper(); got != tt.want {
				t.Errorf("GetCurrencyListUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrencyList_GetValueByCode(t *testing.T) {
	type fields struct {
		Base                string
		AvailableCurrencies map[string]uint8
	}
	type args struct {
		code string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint8
	}{
		{
			name:   "Get value",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{"USD": 1, "EUR": 2}},
			args:   args{code: "EUR"},
			want:   2,
		},
		{
			name:   "Zero value",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{}},
			args:   args{code: "MXN"},
			want:   uint8(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &CurrencyList{
				Base:                tt.fields.Base,
				AvailableCurrencies: tt.fields.AvailableCurrencies,
			}
			if got := cl.GetValueByCode(tt.args.code); got != tt.want {
				t.Errorf("GetValueByCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrencyList_IsCurrencyAvailable(t *testing.T) {
	type fields struct {
		Base                string
		AvailableCurrencies map[string]uint8
	}
	type args struct {
		code string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Available",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{"USD": 1, "EUR": 2}},
			args:   args{code: "EUR"},
			want:   true,
		},
		{
			name:   "Not available",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{"USD": 1, "EUR": 2}},
			args:   args{code: "MXN"},
			want:   false,
		},
		{
			name:   "Empty list",
			fields: fields{Base: "USD", AvailableCurrencies: map[string]uint8{}},
			args:   args{code: "MXN"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &CurrencyList{
				Base:                tt.fields.Base,
				AvailableCurrencies: tt.fields.AvailableCurrencies,
			}
			if got := cl.IsCurrencyAvailable(tt.args.code); got != tt.want {
				t.Errorf("IsCurrencyAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
