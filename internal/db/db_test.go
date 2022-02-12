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

func compareTimeWithMilliSecondEpsilon(a, b time.Time) bool {
	x := a.Round(time.Millisecond)
	y := b.Round(time.Millisecond)
	return x.Equal(y)
}

var timeCmp = cmp.Comparer(compareTimeWithMilliSecondEpsilon)

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

	c := NewController(db)
	if c == nil {
		t.Fatal("failed to create db controller!")
	}

	user := User{
		Name: "Lassi",
		Role: int8(RoleRoot | RoleAdministrator),
	}
	loginInfo := LoginInfo{
		Username: "lassi",
		Password: "changeme",
	}
	machine := Machine{
		Name:        "Lassi's Raspberry Pi",
		Description: "the one on the desk",
		OS:          "Linux",
		Arch:        "armv6l",
	}
	org := Organization{
		Name: "example-org",
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
		err := c.CreateUser(&user)
		if err != nil {
			t.Fatal("error creating User:", err)
		}
	})
	userToken.UserID = user.ID
	loginInfo.UserID = user.ID
	org.Users = append(org.Users, user)

	t.Run("test login account creation", func(t *testing.T) {
		err := c.CreateLoginInfo(&loginInfo)
		if err != nil {
			t.Fatal("error creating LoginInfo:", err)
		}
	})

	t.Run("test machine creation", func(t *testing.T) {
		err := c.CreateMachine(&machine)
		if err != nil {
			t.Fatal("error creating Machine:", err)
		}
	})
	machineToken.MachineID = machine.ID
	schedule.MachineID = machine.ID
	org.Machines = append(org.Machines, machine)

	t.Run("test organization creation", func(t *testing.T) {
		err := c.CreateOrganization(&org)
		if err != nil {
			t.Fatal("error creating Organization:", err)
		}
	})

	user.OrganizationID = org.ID
	machine.OrganizationID = org.ID
	userToken.User = user
	machineToken.Machine = machine
	loginInfo.UserID = user.ID
	loginInfo.User = user

	record := Record{
		MachineID: machine.ID,
		TaskID:    task.ID,
		Timestamp: time.Now(),
		Status:    0,
		Output:    "success",
	}

	t.Run("test schedule creation", func(t *testing.T) {
		err := c.CreateSchedule(&schedule)
		if err != nil {
			t.Fatal("error creating Schedule:", err)
		}
	})

	t.Run("test task creation", func(t *testing.T) {
		err := c.CreateTask(&task)
		if err != nil {
			t.Fatal("error creating Task:", err)
		}
	})

	t.Run("test user token creation", func(t *testing.T) {
		err := c.CreateUserToken(&userToken)
		if err != nil {
			t.Fatal("error creating UserToken:", err)
		}
	})

	t.Run("test machine token creation", func(t *testing.T) {
		err := c.CreateMachineToken(&machineToken)
		if err != nil {
			t.Fatal("error creating MachineToken:", err)
		}
	})

	t.Run("test record creation", func(t *testing.T) {
		err := c.CreateRecord(&record)
		if err != nil {
			t.Fatal("error creating Record:", err)
		}
	})

	t.Run("test user read", func(t *testing.T) {
		name := user.Name
		u, err := c.ReadUser(name)
		if err != nil {
			t.Fatal("error reading User:", err)
		}
		if u == nil {
			t.Fatal("nil User returned")
		}
		if !cmp.Equal(user, *u, timeCmp) {
			t.Fatal(cmp.Diff(user, *u, timeCmp))
		}
	})

	t.Run("test machine read", func(t *testing.T) {
		name := machine.Name
		m, err := c.ReadMachine(name)
		if err != nil {
			t.Fatal("error reading Machine:", err)
		}
		if m == nil {
			t.Fatal("nil Machine returned")
		}
		if !cmp.Equal(machine, *m, timeCmp) {
			t.Fatal(cmp.Diff(machine, *m, timeCmp))
		}
	})

	t.Run("test task read", func(t *testing.T) {
		name := task.Name
		tsk, err := c.ReadTask(name)
		if err != nil {
			t.Fatal("error reading Task:", err)
		}
		if tsk == nil {
			t.Fatal("nil Task returned")
		}
		if !cmp.Equal(task, *tsk, timeCmp) {
			t.Fatal(cmp.Diff(task, *tsk, timeCmp))
		}
	})

	t.Run("test organization read", func(t *testing.T) {
		name := org.Name
		o, err := c.ReadOrganization(name)
		if err != nil {
			t.Fatal("error reading Organization:", err)
		}
		if o == nil {
			t.Fatal("nil Organization returned")
		}
		if !cmp.Equal(org, *o, timeCmp) {
			t.Fatal(cmp.Diff(org, *o, timeCmp))
		}
	})

	t.Run("test schedule read", func(t *testing.T) {
		machineName := machine.Name
		s, err := c.ReadSchedule(machineName)
		if err != nil {
			t.Fatal("error reading Schedule:", err)
		}
		if s == nil {
			t.Fatal("nil Schedule returned")
		}
		if !cmp.Equal(schedule, *s, timeCmp) {
			t.Fatal(cmp.Diff(schedule, *s, timeCmp))
		}
	})

	t.Run("test user token read", func(t *testing.T) {
		tok := userToken.Value
		ut, err := c.ReadUserToken(tok)
		if err != nil {
			t.Fatal("error reading UserToken:", err)
		}
		if ut == nil {
			t.Fatal("nil UserToken returned")
		}
		if !cmp.Equal(userToken, *ut, timeCmp) {
			t.Fatal(cmp.Diff(userToken, *ut, timeCmp))
		}
	})

	t.Run("test machine token read", func(t *testing.T) {
		tok := machineToken.Value
		mt, err := c.ReadMachineToken(tok)
		if err != nil {
			t.Fatal("error reading MachineToken:", err)
		}
		if mt == nil {
			t.Fatal("nil MachineToken returned")
		}
		if !cmp.Equal(machineToken, *mt, timeCmp) {
			t.Fatal(cmp.Diff(machineToken, *mt, timeCmp))
		}
	})

	t.Run("test login info read", func(t *testing.T) {
		u := loginInfo.Username
		l, err := c.ReadLoginInfo(u)
		if err != nil {
			t.Fatal("error reading LoginInfo:", err)
		}
		if l == nil {
			t.Fatal("nil LoginInfo returned")
		}
		if !cmp.Equal(loginInfo, *l, timeCmp) {
			t.Fatal(cmp.Diff(loginInfo, *l, timeCmp))
		}
	})

	t.Run("test record read", func(t *testing.T) {
		r, err := c.ReadRecords(machine.Name)
		if err != nil {
			t.Fatal("error reading machine Records:", err)
		}
		if len(r) == 0 {
			t.Fatal("no records found")
		}
	})

	t.Run("update user", func(t *testing.T) {
		user.Name = "Lassi2"
		err := c.UpdateUser(&user)
		if err != nil {
			t.Fatal("error updating User:", err)
		}
	})

	t.Run("update machine", func(t *testing.T) {
		machine.Name = "the-other-raspbi"
		err := c.UpdateMachine(&machine)
		if err != nil {
			t.Fatal("error updating Machine:", err)
		}
	})

	t.Run("update organization", func(t *testing.T) {
		org.Name = "example-org2"
		err := c.UpdateOrganization(&org)
		if err != nil {
			t.Fatal("error updating Organization:", err)
		}
	})

	t.Run("update schedule", func(t *testing.T) {
		sj := pgtype.JSON{}
		sj.Set([]byte(`{"description":"schedule","schedule":{"task":"echo 1 > /dev/null"}}`))
		schedule.Content = sj
		err := c.UpdateSchedule(&schedule)
		if err != nil {
			t.Fatal("error updating Schedule:", err)
		}
	})

	t.Run("update task", func(t *testing.T) {
		task.Name = "shutoff"
		err := c.UpdateTask(&task)
		if err != nil {
			t.Fatal("error updating Task:", err)
		}
	})

	t.Run("update user token", func(t *testing.T) {
		userToken.Expiration = time.Now().AddDate(1, 0, 0)
		err := c.UpdateUserToken(&userToken)
		if err != nil {
			t.Fatal("error updating UserToken:", err)
		}
	})

	t.Run("update machine token", func(t *testing.T) {
		machineToken.Expiration = time.Now().AddDate(3, 0, 0)
		err := c.UpdateMachineToken(&machineToken)
		if err != nil {
			t.Fatal("error updating MachineToken:", err)
		}
	})

	t.Run("update login info", func(t *testing.T) {
		loginInfo.Password = "p4ssw0rd"
		err := c.UpdateLoginInfo(&loginInfo)
		if err != nil {
			t.Fatal("error updating LoginInfo:", err)
		}
	})

	t.Run("update record", func(t *testing.T) {
		record.Output = "redacted"
		err := c.UpdateRecord(&record)
		if err != nil {
			t.Fatal("error updating record:", err)
		}
	})

	t.Run("delete user", func(t *testing.T) {
		err := c.DeleteUser(user.Name)
		if err != nil {
			t.Fatal("error deleting User:", err)
		}
	})

	// delete schedule before machine, otherwise won't be able to find the schedule with machine name
	t.Run("delete schedule", func(t *testing.T) {
		// this will fail if updating machine also failed
		err := c.DeleteSchedule(machine.Name)
		if err != nil {
			t.Fatal("error deleting Schedule:", err)
		}
	})

	// delete records before machne, otherwise won't be able to find the records with machine name
	t.Run("delete machine records", func(t *testing.T) {
		err := c.DeleteRecords(machine.Name)
		if err != nil {
			t.Fatal("error deleting machine Records:", err)
		}
	})

	t.Run("delete machine", func(t *testing.T) {
		err := c.DeleteMachine(machine.Name)
		if err != nil {
			t.Fatal("error deleting Machine:", err)
		}
	})

	t.Run("delete organization", func(t *testing.T) {
		err := c.DeleteOrganization(org.Name)
		if err != nil {
			t.Fatal("error deleting Organization:", err)
		}
	})

	t.Run("delete task", func(t *testing.T) {
		err := c.DeleteTask(task.Name)
		if err != nil {
			t.Fatal("error deleting Task:", err)
		}
	})

	t.Run("delete user token", func(t *testing.T) {
		err := c.DeleteUserToken(userToken.Value)
		if err != nil {
			t.Fatal("error deleting UserToken:", err)
		}
	})

	t.Run("delete machine token", func(t *testing.T) {
		err := c.DeleteMachineToken(machineToken.Value)
		if err != nil {
			t.Fatal("error deleting MachineToken:", err)
		}
	})

	t.Run("delete login info", func(t *testing.T) {
		err := c.DeleteLoginInfo(loginInfo.Username)
		if err != nil {
			t.Fatal("error deleting LoginInfo:", err)
		}
	})
}
