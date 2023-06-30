package main

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

func Test_selectUsers(t *testing.T) {
	shutdown := initTracer()
	defer shutdown()
	dsn := "root:@tcp(localhost:3306)/otelsql?parseTime=true"
	db, err := otelsql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the MySQL database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name string
		want []*User
	}{
		{
			name: "work",
			want: []*User{
				{
					UserID: 1,
					Name:   "taro yamada",
					Email:  "taro-yamada@example.com",
				},
				{
					UserID: 2,
					Name:   "taro go",
					Email:  "taro-go@example.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectUsers(db)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
