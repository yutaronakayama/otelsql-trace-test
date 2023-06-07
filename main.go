package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/yutaronakayama/otelsql-trace-test/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

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

func main() {
	shutdown := initTracer()
	defer shutdown()

	ctx := context.Background()

	// MySQLデータベースへの接続
	dsn := "root:@tcp(localhost:3306)/otelsql?parseTime=true"
	db, err := otelsql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the MySQL database: %v", err)
	}
	defer db.Close()

	// 取得したデータベース情報をsqlxでラップ
	sqlxDB := sqlx.NewDb(db, "mysql")

	// ユーザテーブル情報を取得
	type UserRepo domain.UserRepository
	users, err := UserRepo.SelectUsers(ctx, sqlxDB)
	if err != nil {
		fmt.Errorf("failed to select users: %v", err)
	}

	for _, user := range users {
		fmt.Printf("User ID: %d, Name: %s, Email: %s\n", user.UserID, user.Name, user.Email)
	}
}
