package main

import (
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

func Test_selectUsers(t *testing.T) {
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
			name: "問題なくUser情報が取得できること",
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
			got, _ := selectUsers(db)
			fmt.Println(got)
			// fmt.Printf("got: %v", got)
			// fmt.Printf("want: %v", tt.want)
			assert.ElementsMatch(t, tt.want, got)
			// assert.Equal(t, v.UserID, tt.want.UserID)
			// assert.Equal(t, v.Name, tt.want.Name)
			// assert.Equal(t, v.Email, tt.want.Email)

		})
	}
}
