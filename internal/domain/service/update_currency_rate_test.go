package currencyrates

import (
	"context"
	"reflect"
	"testing"

	goredis "github.com/go-redis/redis/v8"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"github.com/vkuzmenkova/currency-rates/internal/db"
	"github.com/vkuzmenkova/currency-rates/internal/domain"
)

func TestCurrenciesService_UpdateRate(t *testing.T) {
	type fields struct {
		Repo         *db.CurrenciesRepo
		Enqueuer     *work.Enqueuer
		Pool         *work.WorkerPool
		CurrencyList *domain.CurrencyList
		KV           *goredis.Client
	}
	type args struct {
		ctx          context.Context
		base         string
		currencyCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CurrenciesService{
				Repo:         tt.fields.Repo,
				Enqueuer:     tt.fields.Enqueuer,
				Pool:         tt.fields.Pool,
				CurrencyList: tt.fields.CurrencyList,
				KV:           tt.fields.KV,
			}
			got, err := s.UpdateRate(tt.args.ctx, tt.args.base, tt.args.currencyCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateRate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
