package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type User struct {
	UserID int    `db:"user_id"`
	Name   string `db:"name"`
	Email  string `db:"email"`
}

func initTracer() func() {
	// Set the Jaeger exporter to send traces to the Jaeger instance.
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Configure the trace provider with the Jaeger exporter and a sampler that always samples.
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)

	// Set the propagators for the global OpenTelemetry configuration.
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
	var users []User
	err = sqlxDB.SelectContext(ctx, &users, "SELECT user_id, name, email FROM users")
	if err != nil {
		log.Fatalf("failed to query users table: %v", err)
	}

	for _, user := range users {
		fmt.Printf("User ID: %d, Name: %s, Email: %s\n", user.UserID, user.Name, user.Email)
	}
}
