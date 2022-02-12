package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/LassiHeikkila/taskey/internal/db"
)

const (
	postgresHost = "127.0.0.1"
	postgresPort = 5432
	postgresUser = "postgres"
	postgresPw   = "test1234"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, os.Interrupt)
		<-sc
		log.Println("caught interrupt signal")
		cancel()
	}()
	s := run(ctx)
	os.Exit(s)
}

func run(ctx context.Context) int {
	d := db.OpenDB(
		db.WithHost(postgresHost),
		db.WithPort(postgresPort),
		db.WithUsername(postgresUser),
		db.WithPassword(postgresPw),
		db.WithSSLMode("disable"),
	)
	if d == nil {
		log.Println("failed to initialize database!")
		return 1
	}

	if err := db.InitializeDB(d); err != nil {
		log.Println("error initalizing db:", err)
		return 1
	}

	c := db.NewController(d)
	if c == nil {
		log.Println("failed to create db controller!")
		return 1
	}

	log.Println("db initialized")

	<-ctx.Done()
	return 0
}
