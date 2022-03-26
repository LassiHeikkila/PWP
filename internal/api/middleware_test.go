package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/LassiHeikkila/taskey/internal/auth/mock"
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/internal/db/mock"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func TestAuthenticatedUserMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)

	expectedUserToken := &db.UserToken{
		User: db.User{
			Name:           "Lassi",
			Email:          "lassi@example.com",
			OrganizationID: 123,
			Role:           types.RoleUser | types.RoleMaintainer | types.RoleAdministrator | types.RoleRoot,
		},
	}
	d.EXPECT().ReadUserToken(db.StringToUUID(`cf6525ce-9fbb-4cd1-a1f1-d96f4220b3d2`)).Return(expectedUserToken, nil)

	called := false

	next := func(w http.ResponseWriter, req *http.Request) {
		called = true

		_, _ = w.Write([]byte(`ok`))
	}

	// check that when token validation returns a user, next is called

	mw := NewAuthUserMiddleware(next, a, d, types.RoleUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/example", nil)
	req.Header.Set("Authorization", "Key cf6525ce-9fbb-4cd1-a1f1-d96f4220b3d2")
	mw.ServeHTTP(w, req)

	if !called {
		t.Fatal("next not called")
	}

	// check that when token validation returns an error, next is not called

	called = false
	d.EXPECT().ReadUserToken(db.StringToUUID(`91f2a925-a2c5-4a59-aa13-adeb7950faa8`)).Return(nil, errors.New("invalid token"))

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/example", nil)
	req2.Header.Set("Authorization", "Key 91f2a925-a2c5-4a59-aa13-adeb7950faa8")
	mw.ServeHTTP(w2, req2)

	if called {
		t.Fatal("next called!")
	}
}

func TestAuthenticatedMachineMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)

	expectedMachine := &types.Machine{
		Name:        "TestMachine",
		Description: "Machine for testing",
		OS:          "Linux",
		Arch:        "x86_64",
	}
	expectedMachineToken := &db.MachineToken{
		Machine: db.Machine{
			Name:        "TestMachine",
			Description: "Machine for testing",
			OS:          "Linux",
			Arch:        "x86_64",
		},
	}
	d.EXPECT().ReadMachineToken(db.StringToUUID(`519aa433-418e-4fc2-bd72-5d196a62fc85`)).Return(expectedMachineToken, nil)

	called := false
	var calledWithMachine *types.Machine

	next := AuthenticatedMachineHandler(func(w http.ResponseWriter, req *http.Request, machine *types.Machine) {
		called = true
		calledWithMachine = machine

		_, _ = w.Write([]byte(`ok`))
	})

	// check that when token validation returns a machine, next is called

	mw := NewAuthMachineMiddleware(next, a, d)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/example", nil)
	req.Header.Set("Authorization", "Key 519aa433-418e-4fc2-bd72-5d196a62fc85")
	mw.ServeHTTP(w, req)

	if !called {
		t.Fatal("next not called")
	}
	require.Equal(t, expectedMachine, calledWithMachine)

	// check that when token validation fails, next is not called

	called = false
	d.EXPECT().ReadMachineToken(db.StringToUUID(`c87fcad1-2f2a-40f2-8da0-38712863003a`)).Return(nil, errors.New("invalid token"))

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/example", nil)
	req2.Header.Set("Authorization", "Key c87fcad1-2f2a-40f2-8da0-38712863003a")
	mw.ServeHTTP(w2, req2)

	if called {
		t.Fatal("next called")
	}
}
