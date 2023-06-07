package repository

import (
	"context"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nakayamayutaro/otelsql-trace-test/domain"
)

func Test_selectUsers(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *sqlx.DB
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectUsers(tt.args.ctx, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("selectUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
