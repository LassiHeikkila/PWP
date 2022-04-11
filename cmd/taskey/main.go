package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/handlers"

	"github.com/LassiHeikkila/taskey/internal/api"
	"github.com/LassiHeikkila/taskey/internal/auth"
	"github.com/LassiHeikkila/taskey/internal/db"
)

const (
	dbHostEnvKey     = "TASKEYDBHOST"
	dbPortEnvKey     = "TASKEYDBPORT"
	dbUserEnvKey     = "TASKEYDBUSER"
	dbPasswordEnvKey = "TASKEYDBPASSWORD"
	dbDbEnvKey       = "TASKEYDBDB"
	jwtKeyEnvKey     = "TASKEYJWTKEY"
)

var (
	dbHost     = getEnvOrDefault(dbHostEnvKey, defaultDbHost)
	dbPort     = getEnvOrDefaultInt(dbPortEnvKey, defaultDbPort)
	dbDb       = getEnvOrDefault(dbDbEnvKey, defaultDbDb)
	dbUser     = os.Getenv(dbUserEnvKey)
	dbPassword = os.Getenv(dbPasswordEnvKey)

	privateKey = os.Getenv(jwtKeyEnvKey)

	httpPort = defaultHttpPort
)

func main() {
	flag.IntVar(&httpPort, "port", defaultHttpPort, "Port to use for HTTP interface")

	flag.Parse()

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
	if err := validateArgs(); err != nil {
		log.Println("error validating arguments:", err)
		return 1
	}
	d := db.OpenDB(
		db.WithHost(dbHost),
		db.WithPort(dbPort),
		db.WithUsername(dbUser),
		db.WithPassword(dbPassword),
		db.WithSSLMode("disable"),
		db.WithDBName(dbDb),
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
	if err := h.RegisterSignUpHandlers(); err != nil {
		log.Println("failed to register signup routes!")
		return 1
	}
	if err := h.RegisterExtraRoute("/api/v1/", ServeSwaggerUI); err != nil {
		log.Println("failed to register swagger UI route!")
		return 2
	}
	if err := h.RegisterExtraRoute("/api/v1/openapi.yml", ServeOpenAPI); err != nil {
		log.Println("failed to register openapi route!")
		return 2
	}

	log.Println("API handler initialized")

	srv := &http.Server{
		Handler:      handlers.CombinedLoggingHandler(log.Writer(), api.ExecutionTimeHandler(h)),
		Addr:         fmt.Sprintf(":%d", httpPort),
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

func validateArgs() error {
	if dbHost == "" {
		return errors.New("empty db host")
	}
	if dbPort == 0 {
		return errors.New("empty db port")
	}
	if dbUser == "" {
		return errors.New("empty db user")
	}
	if dbPassword == "" {
		return errors.New("empty db password")
	}
	if privateKey == "" {
		return errors.New("empty private key for token generator")
	}
	return nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvOrDefaultInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, _ := strconv.Atoi(v)
	return i
}
