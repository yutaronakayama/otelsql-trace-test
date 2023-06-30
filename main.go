package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type User struct {
	UserID int    `db:"id"`
	Name   string `db:"name"`
	Email  string `db:"email"`
}

func initTracer() func() {
	// Jaegerへトレース情報を送るためのエクスポータの作成
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")),
	)
	if err != nil {
		log.Fatal(err)
	}

	// トレースプロバイダの設定
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}
}

func selectUsers(db *sql.DB) ([]*User, error) {
	users, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer users.Close()

	ret := []*User{}
	for users.Next() {
		user := &User{}
		if err := users.Scan(&user.UserID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		ret = append(ret, user)
	}
	return ret, nil
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	// MySQLデータベースへの接続
	dsn := "root:@tcp(localhost:3306)/otelsql?parseTime=true"
	db, err := otelsql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the MySQL database: %v", err)
	}
	defer db.Close()

	// ユーザ情報を取得
	users, err := selectUsers(db)
	if err != nil {
		log.Fatalf("failed to select users: %v", err)
	}
	fmt.Println(users)
}
