package currencyrates

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/vkuzmenkova/currency-rates/internal/db"
	"github.com/vkuzmenkova/currency-rates/internal/domain"
	"github.com/vkuzmenkova/currency-rates/mocks"
	"github.com/vkuzmenkova/currency-rates/models"
)

func TestCurrenciesService_GetLastRate(t *testing.T) {
	type args struct {
		ctx          context.Context
		base         string
		currencyCode string
		error        error
	}
	tests := []struct {
		name    string
		args    args
		want    models.CurrencyRate
		wantErr bool
	}{
		{
			name:    "CurrencyRate",
			args:    args{ctx: context.Background(), base: "USD", currencyCode: "EUR", error: nil},
			want:    models.CurrencyRate{UpdatedAt: time.Time{}.String(), Base: "USD", Currency: "EUR", Value: 0},
			wantErr: false,
		},
		{
			name:    "ErrNoRows",
			args:    args{ctx: context.Background(), base: "USD", currencyCode: "EUR", error: pgx.ErrNoRows},
			want:    models.CurrencyRate{},
			wantErr: true,
		},
		{
			name:    "Error",
			args:    args{ctx: context.Background(), base: "USD", currencyCode: "EUR", error: errors.New("err")},
			want:    models.CurrencyRate{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			conn := mocks.NewDB(t)
			row := mocks.NewRow(t)
			var updatedAt time.Time
			var value float64

			row.On("Scan", &value, &updatedAt).Return(tt.args.error)
			conn.On("QueryRow", mock.Anything, mock.Anything, uint8(1), uint8(2)).Return(row)

			repo := &db.CurrenciesRepo{Conn: conn}
			s := &CurrenciesService{Repo: repo, CurrencyList: domain.NewCurrencyList()}

			got, err := s.GetLastRate(tt.args.ctx, tt.args.base, tt.args.currencyCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastRate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
