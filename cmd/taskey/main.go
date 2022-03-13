package main

import (
	"context"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/LassiHeikkila/taskey/internal/api"
	"github.com/LassiHeikkila/taskey/internal/auth"
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

	log.Println("db handler initialized")

	privateKey := os.Getenv("TASKEYJWTKEY")
	privKey, err := hex.DecodeString(privateKey)
	if err != nil {
		log.Println("TASKEYJWTKEY not in hex encoded format!")
		return 1
	}

	a := auth.NewController(privKey)
	if a == nil {
		log.Println("failed to create auth controller!")
		return 1
	}

	log.Println("auth handler initialized")

	h := api.NewHandler(a, c)
	if h == nil {
		log.Println("failed to create API handler!")
		return 1
	}

	if err := h.RegisterOrganizationHandlers(); err != nil {
		log.Println("failed to register organization routes!")
		return 1
	}
	if err := h.RegisterUserHandlers(); err != nil {
		log.Println("failed to register user routes!")
		return 1
	}
	if err := h.RegisterMachineHandlers(); err != nil {
		log.Println("failed to register machine routes!")
		return 1
	}
	if err := h.RegisterScheduleHandlers(); err != nil {
		log.Println("failed to register schedule routes!")
		return 1
	}
	if err := h.RegisterTaskHandlers(); err != nil {
		log.Println("failed to register task routes!")
		return 1
	}
	if err := h.RegisterAuthenticationHandlers(); err != nil {
		log.Println("failed to register authentication routes!")
		return 1
	}

	log.Println("API handler initialized")

	srv := &http.Server{
		Handler:      h,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("HTTP server exited with error:", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down HTTP server")
	// give server 15s to shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("error shutting down HTTP server:", err)
	}

	return 0
}
