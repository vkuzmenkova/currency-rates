package currencyrates

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/vkuzmenkova/currency-rates/internal/db"
	"github.com/vkuzmenkova/currency-rates/internal/domain"
	"github.com/vkuzmenkova/currency-rates/mocks"
	"github.com/vkuzmenkova/currency-rates/models"
)

func TestCurrenciesService_GetRateByUUID(t *testing.T) {
	type args struct {
		ctx   context.Context
		uuid  uuid.UUID
		error error
	}
	tests := []struct {
		name    string
		args    args
		want    models.CurrencyRate
		wantErr bool
	}{
		{
			name:    "CurrencyRate",
			args:    args{ctx: context.Background(), uuid: uuid.New(), error: nil},
			want:    models.CurrencyRate{UpdatedAt: time.Time{}.String(), Value: 0},
			wantErr: false,
		},
		{
			name:    "ErrNoRows",
			args:    args{ctx: context.Background(), uuid: uuid.New(), error: pgx.ErrNoRows},
			want:    models.CurrencyRate{},
			wantErr: true,
		},
		{
			name:    "Error",
			args:    args{ctx: context.Background(), uuid: uuid.New(), error: errors.New("err")},
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
			var baseVal, currencyVal uint8

			row.On("Scan", &baseVal, &currencyVal, &value, &updatedAt).Return(tt.args.error)
			conn.On("QueryRow", mock.Anything, mock.Anything, tt.args.uuid).Return(row)
			repo := &db.CurrenciesRepo{Conn: conn}
			s := &CurrenciesService{Repo: repo, CurrencyList: domain.NewCurrencyList()}

			got, err := s.GetRateByUUID(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRateByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRateByUUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
