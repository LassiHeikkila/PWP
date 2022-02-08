package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"

	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func compareTimeWithMicroSecondEpsilon(a, b time.Time) bool {
	delta := a.Sub(b)
	if delta < 0 {
		delta = -delta
	}
	return delta < time.Microsecond
}

var timeCmp = cmp.Comparer(compareTimeWithMicroSecondEpsilon)

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

	user := User{Name: "Lassi"}
	loginInfo := LoginInfo{
		Username: "lassi",
		Password: "changeme",
		User:     user,
	}
	machine := Machine{
		Name:        "Lassi's Raspberry Pi",
		Description: "the one on the desk",
		OS:          "Linux",
		Arch:        "armv6l",
	}
	org := Organization{
		Users:    []User{user},
		Machines: []Machine{machine},
	}
	sj := pgtype.JSON{}
	sj.Set([]byte(`{"some":"json object"}`))
	schedule := Schedule{
		Machine: machine,
		Content: sj,
	}
	tj := pgtype.JSON{}
	tj.Set([]byte(`{"key":"value"}`))
	task := Task{
		Name:        "some task",
		Description: "blink some light for some time",
		Content:     tj,
	}
	machineToken := MachineToken{
		Value: pgtype.UUID{
			Bytes:  [16]byte{0xfa, 0x1a, 0xfe, 0x1},
			Status: pgtype.Present,
		},
		Expiration: time.Now().Add(time.Hour),
		// id is set after adding machine to db, it is added by gorm to the machine object
	}
	userToken := UserToken{
		Value: pgtype.UUID{
			Bytes:  [16]byte{0xf0, 0x07, 0xba, 0x11},
			Status: pgtype.Present,
		},
		Expiration: time.Now().Add(time.Hour),
		// id is set after adding user to db, it is added by gorm to the user object
	}

	t.Run("test user creation", func(t *testing.T) {
		res := db.Create(&user)
		if err := res.Error; err != nil {
			t.Fatal("error creating User:", err)
		}
	})
	userToken.UserID = user.ID

	t.Run("test login account creation", func(t *testing.T) {
		res := db.Create(&loginInfo)
		if err := res.Error; err != nil {
			t.Fatal("error creating LoginInfo:", err)
		}
	})

	t.Run("test machine creation", func(t *testing.T) {
		res := db.Create(&machine)
		if err := res.Error; err != nil {
			t.Fatal("error creating Machine:", err)
		}
	})
	machineToken.MachineID = machine.ID

	t.Run("test organization creation", func(t *testing.T) {
		res := db.Create(&org)
		if err := res.Error; err != nil {
			t.Fatal("error creating Organization:", err)
		}
	})

	t.Run("test schedule creation", func(t *testing.T) {
		res := db.Create(&schedule)
		if err := res.Error; err != nil {
			t.Fatal("error creating Schedule:", err)
		}
	})

	t.Run("test task creation", func(t *testing.T) {
		res := db.Create(&task)
		if err := res.Error; err != nil {
			t.Fatal("error creating Task:", err)
		}
	})

	t.Run("test user token creation", func(t *testing.T) {
		res := db.Create(&userToken)
		if err := res.Error; err != nil {
			t.Fatal("error creating UserToken:", err)
		}
	})

	t.Run("test machine token creation", func(t *testing.T) {
		res := db.Create(&machineToken)
		if err := res.Error; err != nil {
			t.Fatal("error creating MachineToken:", err)
		}
	})

	t.Run("test user read", func(t *testing.T) {
		u := User{}
		// just user First, we know there's only one
		res := db.First(&u)
		if err := res.Error; err != nil {
			t.Fatal("non-nil error returned:", err)
		}
		if !cmp.Equal(user, u, timeCmp) {
			t.Fatal(cmp.Diff(user, u, timeCmp))
		}
	})
}
