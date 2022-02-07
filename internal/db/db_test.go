package db

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestDBIntegration(t *testing.T) {
	now := time.Now().UTC()

	if testing.Short() {
		t.Skip()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const (
		pw = "test1234"
	)

	var terminateContainerFunc func()
	var containerIP string
	var containerPort int

	t.Run("setup container", func(t *testing.T) {
		req := testcontainers.ContainerRequest{
			Image:        "postgres",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForListeningPort("5432/tcp"),
			Env: map[string]string{
				"POSTGRES_PASSWORD": pw,
			},
		}
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			t.Fatal("failed to start postgres container:", err)
		}
		terminateContainerFunc = func() {
			if err := container.Terminate(ctx); err != nil {
				t.Log("error terminating container:", err)
			}
		}

		ip, err := container.Host(ctx)
		if err != nil {
			t.Fatal("failed to get container ip:", err)
		}
		containerIP = ip

		port, err := container.MappedPort(ctx, "5432")
		if err != nil {
			t.Fatal("failed to get mapped port:", err)
		}
		containerPort = port.Int()
	})
	if containerIP == "" || containerPort == 0 {
		t.Fatal("failed to setup container, cannot perform rest of tests...")
	}
	defer terminateContainerFunc()

	t.Log("setting up container took:", time.Since(now))
	time.Sleep(time.Second)

	var db *gorm.DB
	t.Run("open and initialize database", func(t *testing.T) {
		db = OpenDB(
			WithHost(containerIP),
			WithPort(containerPort),
			WithUsername("postgres"),
			WithPassword(pw),
			WithSSLMode("disable"),
		)
		if db == nil {
			t.Fatal("got nil db, expecting a non-null pointer")
		}

		if err := InitializeDB(db); err != nil {
			t.Fatal("error initalizing db:", err)
		}
	})

}
