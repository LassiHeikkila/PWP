package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/LassiHeikkila/taskey/internal/auth/mock"
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/internal/db/mock"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

// check that public interface is implemented
var _ Handler = &handler{}

func TestRouteRegistrationUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	err := h.RegisterUserHandlers()
	if err != nil {
		t.Fatal("error returned by handler registration method")
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/org123/users/user456/", nil)
	rm := mux.RouteMatch{}

	matched := h.router.Match(req, &rm)
	if !matched {
		t.Fatal("valid route not matched:", rm.MatchErr)
	}
}

func TestRouteRegistrationOrganization(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	err := h.RegisterOrganizationHandlers()
	if err != nil {
		t.Fatal("error returned by handler registration method")
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/organizations/org123/", nil)
	rm := mux.RouteMatch{}

	matched := h.router.Match(req, &rm)
	if !matched {
		t.Fatal("valid route not matched:", rm.MatchErr)
	}
}

func TestRouteRegistrationMachine(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	err := h.RegisterMachineHandlers()
	if err != nil {
		t.Fatal("error returned by handler registration method")
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/org123/machines/machine456/", nil)
	rm := mux.RouteMatch{}

	matched := h.router.Match(req, &rm)
	if !matched {
		t.Fatal("valid route not matched:", rm.MatchErr)
	}
}

func TestRouteRegistrationTask(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	err := h.RegisterTaskHandlers()
	if err != nil {
		t.Fatal("error returned by handler registration method")
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/org123/tasks/taskXYZ/", nil)
	rm := mux.RouteMatch{}

	matched := h.router.Match(req, &rm)
	if !matched {
		t.Fatal("valid route not matched:", rm.MatchErr)
	}
}

func TestRouteRegistrationSchedule(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	err := h.RegisterScheduleHandlers()
	if err != nil {
		t.Fatal("error returned by handler registration method")
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/org123/machines/machineABC/schedule/", nil)
	rm := mux.RouteMatch{}

	matched := h.router.Match(req, &rm)
	if !matched {
		t.Fatal("valid route not matched:", rm.MatchErr)
	}
}

func TestRouteRegistrationAuth(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	err := h.RegisterAuthenticationHandlers()
	if err != nil {
		t.Fatal("error returned by handler registration method")
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/", nil)
	rm := mux.RouteMatch{}

	matched := h.router.Match(req, &rm)
	if !matched {
		t.Fatal("valid route not matched:", rm.MatchErr)
	}
}

func TestProcessRequestGetOrganization(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	if err := h.RegisterUserHandlers(); err != nil {
		t.Fatal("error registering user handlers:", err)
	}
	if err := h.RegisterOrganizationHandlers(); err != nil {
		t.Fatal("error registering organization handlers:", err)
	}
	if err := h.RegisterMachineHandlers(); err != nil {
		t.Fatal("error registering machine handlers:", err)
	}
	if err := h.RegisterTaskHandlers(); err != nil {
		t.Fatal("error registering task handlers:", err)
	}
	if err := h.RegisterScheduleHandlers(); err != nil {
		t.Fatal("error registering schedule handlers:", err)
	}
	if err := h.RegisterAuthenticationHandlers(); err != nil {
		t.Fatal("error registering authentication handlers:", err)
	}
	if err := h.RegisterRecordHandlers(); err != nil {
		t.Fatal("error registering record handlers:", err)
	}

	server := httptest.NewServer(h)

	if server == nil {
		t.Fatal("failed to create test server")
	}

	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/organizations/org123/", nil)
	req.Header.Set("Authorization", "Bearer my test key")

	a.EXPECT().ValidateUserToken(
		"my test key",
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(tokenString string, user *string, organization *string, role *int) bool {
		if tokenString != "my test key" {
			return false
		}
		if user != nil {
			*user = "user456"
		}
		if organization != nil {
			*organization = "org123"
		}
		if role != nil {
			*role = int(types.RoleUser | types.RoleMaintainer | types.RoleAdministrator | types.RoleRoot)
		}
		return true
	})

	d.EXPECT().ReadOrganization("org123").Return(&db.Organization{
		Model: gorm.Model{
			ID: 123,
		},
		Name: "org123",
	}, nil).Times(2)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error doing request:", err)
	}
	defer resp.Body.Close()

	var response Response
	b, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(b, &response); err != nil {
		t.Fatal("failed to decode response as JSON: \"", err, "\", response was: \"", string(b), "\"")
	}

	if response.Code != 200 {
		t.Fatal("response not 200:", response)
	}
}

func TestProcessRequestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	if err := h.RegisterUserHandlers(); err != nil {
		t.Fatal("error registering user handlers:", err)
	}
	if err := h.RegisterOrganizationHandlers(); err != nil {
		t.Fatal("error registering organization handlers:", err)
	}
	if err := h.RegisterMachineHandlers(); err != nil {
		t.Fatal("error registering machine handlers:", err)
	}
	if err := h.RegisterTaskHandlers(); err != nil {
		t.Fatal("error registering task handlers:", err)
	}
	if err := h.RegisterScheduleHandlers(); err != nil {
		t.Fatal("error registering schedule handlers:", err)
	}
	if err := h.RegisterAuthenticationHandlers(); err != nil {
		t.Fatal("error registering authentication handlers:", err)
	}
	if err := h.RegisterRecordHandlers(); err != nil {
		t.Fatal("error registering record handlers:", err)
	}

	server := httptest.NewServer(h)

	if server == nil {
		t.Fatal("failed to create test server")
	}

	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/org123/users/user456/", nil)
	req.Header.Set("Authorization", "Bearer my test key")

	a.EXPECT().ValidateUserToken(
		"my test key",
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(tokenString string, user *string, organization *string, role *int) bool {
		if tokenString != "my test key" {
			return false
		}
		if user != nil {
			*user = "user456"
		}
		if organization != nil {
			*organization = "org123"
		}
		if role != nil {
			*role = int(types.RoleUser | types.RoleMaintainer | types.RoleAdministrator | types.RoleRoot)
		}
		return true
	})

	d.EXPECT().ReadUser("user456").Return(&db.User{
		Model: gorm.Model{
			ID: 456,
		},
		Name:           "user456",
		Email:          "lassi@example.com",
		OrganizationID: 123,
		Organization: db.Organization{
			Model: gorm.Model{
				ID: 123,
			},
			Name: "org123",
		},
		Role: types.RoleUser | types.RoleMaintainer | types.RoleAdministrator | types.RoleRoot,
	}, nil)
	d.EXPECT().ReadOrganization("org123").Return(&db.Organization{
		Model: gorm.Model{
			ID: 123,
		},
		Name: "org123",
	}, nil).Times(2)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error doing request:", err)
	}
	defer resp.Body.Close()

	var response Response
	b, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(b, &response); err != nil {
		t.Fatal("failed to decode response as JSON: \"", err, "\", response was: \"", string(b), "\"")
	}

	if response.Code != 200 {
		t.Fatal("response not 200:", response)
	}
}

func TestProcessRequestGetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	if err := h.RegisterUserHandlers(); err != nil {
		t.Fatal("error registering user handlers:", err)
	}
	if err := h.RegisterOrganizationHandlers(); err != nil {
		t.Fatal("error registering organization handlers:", err)
	}
	if err := h.RegisterMachineHandlers(); err != nil {
		t.Fatal("error registering machine handlers:", err)
	}
	if err := h.RegisterTaskHandlers(); err != nil {
		t.Fatal("error registering task handlers:", err)
	}
	if err := h.RegisterScheduleHandlers(); err != nil {
		t.Fatal("error registering schedule handlers:", err)
	}
	if err := h.RegisterAuthenticationHandlers(); err != nil {
		t.Fatal("error registering authentication handlers:", err)
	}
	if err := h.RegisterRecordHandlers(); err != nil {
		t.Fatal("error registering record handlers:", err)
	}

	server := httptest.NewServer(h)

	if server == nil {
		t.Fatal("failed to create test server")
	}

	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/org123/users/", nil)
	req.Header.Set("Authorization", "Bearer my test key")

	a.EXPECT().ValidateUserToken(
		"my test key",
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(tokenString string, user *string, organization *string, role *int) bool {
		if tokenString != "my test key" {
			return false
		}
		if user != nil {
			*user = "user456"
		}
		if organization != nil {
			*organization = "org123"
		}
		if role != nil {
			*role = int(types.RoleUser | types.RoleMaintainer | types.RoleAdministrator | types.RoleRoot)
		}
		return true
	})

	user456 := db.User{
		Model: gorm.Model{
			ID: 456,
		},
		Name:           "user456",
		Email:          "lassi@example.com",
		OrganizationID: 123,
		Organization: db.Organization{
			Model: gorm.Model{
				ID: 123,
			},
			Name: "org123",
		},
		Role: types.RoleUser | types.RoleMaintainer | types.RoleAdministrator | types.RoleRoot,
	}
	user567 := db.User{
		Model: gorm.Model{
			ID: 567,
		},
		Name:           "user567",
		Email:          "totallynotlassi@example.com",
		OrganizationID: 123,
		Organization: db.Organization{
			Model: gorm.Model{
				ID: 123,
			},
			Name: "org123",
		},
		Role: types.RoleUser | types.RoleMaintainer,
	}

	d.EXPECT().ReadUser("user456").Return(&user456, nil)
	d.EXPECT().ReadUser("user567").Return(&user567, nil)
	d.EXPECT().ReadOrganization("org123").Return(&db.Organization{
		Model: gorm.Model{
			ID: 123,
		},
		Name: "org123",
		Users: []db.User{
			user456,
			user567,
		},
	}, nil).Times(2)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error doing request:", err)
	}
	defer resp.Body.Close()

	var response Response
	b, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(b, &response); err != nil {
		t.Fatal("failed to decode response as JSON: \"", err, "\", response was: \"", string(b), "\"")
	}

	if response.Code != 200 {
		t.Fatal("response not 200:", response)
	}
}

func TestProcessRequestGetMachine(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	if err := h.RegisterUserHandlers(); err != nil {
		t.Fatal("error registering user handlers:", err)
	}
	if err := h.RegisterOrganizationHandlers(); err != nil {
		t.Fatal("error registering organization handlers:", err)
	}
	if err := h.RegisterMachineHandlers(); err != nil {
		t.Fatal("error registering machine handlers:", err)
	}
	if err := h.RegisterTaskHandlers(); err != nil {
		t.Fatal("error registering task handlers:", err)
	}
	if err := h.RegisterScheduleHandlers(); err != nil {
		t.Fatal("error registering schedule handlers:", err)
	}
	if err := h.RegisterAuthenticationHandlers(); err != nil {
		t.Fatal("error registering authentication handlers:", err)
	}
	if err := h.RegisterRecordHandlers(); err != nil {
		t.Fatal("error registering record handlers:", err)
	}

	server := httptest.NewServer(h)

	if server == nil {
		t.Fatal("failed to create test server")
	}

	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/org123/machines/machineXYZ/", nil)
	req.Header.Set("Authorization", "Bearer my test key")

	a.EXPECT().ValidateUserToken(
		"my test key",
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(tokenString string, user *string, organization *string, role *int) bool {
		if tokenString != "my test key" {
			return false
		}
		if user != nil {
			*user = "user456"
		}
		if organization != nil {
			*organization = "org123"
		}
		if role != nil {
			*role = int(types.RoleUser | types.RoleMaintainer | types.RoleAdministrator | types.RoleRoot)
		}
		return true
	})

	d.EXPECT().ReadMachine("machineXYZ").Return(&db.Machine{
		Model: gorm.Model{
			ID: 678,
		},
		Name:           "machineXYZ",
		Description:    "test machine",
		OS:             "linux",
		Arch:           "amd64",
		OrganizationID: 123,
	}, nil)
	d.EXPECT().ReadOrganization("org123").Return(&db.Organization{
		Model: gorm.Model{
			ID: 123,
		},
		Name: "org123",
	}, nil).Times(2)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error doing request:", err)
	}
	defer resp.Body.Close()

	var response Response
	b, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(b, &response); err != nil {
		t.Fatal("failed to decode response as JSON: \"", err, "\", response was: \"", string(b), "\"")
	}

	if response.Code != 200 {
		t.Fatal("response not 200:", response)
	}
}

func TestProcessRequestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	if err := h.RegisterUserHandlers(); err != nil {
		t.Fatal("error registering user handlers:", err)
	}
	if err := h.RegisterOrganizationHandlers(); err != nil {
		t.Fatal("error registering organization handlers:", err)
	}
	if err := h.RegisterMachineHandlers(); err != nil {
		t.Fatal("error registering machine handlers:", err)
	}
	if err := h.RegisterTaskHandlers(); err != nil {
		t.Fatal("error registering task handlers:", err)
	}
	if err := h.RegisterScheduleHandlers(); err != nil {
		t.Fatal("error registering schedule handlers:", err)
	}
	if err := h.RegisterAuthenticationHandlers(); err != nil {
		t.Fatal("error registering authentication handlers:", err)
	}
	if err := h.RegisterRecordHandlers(); err != nil {
		t.Fatal("error registering record handlers:", err)
	}

	server := httptest.NewServer(h)

	if server == nil {
		t.Fatal("failed to create test server")
	}

	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/org123/tasks/task1234/", nil)
	req.Header.Set("Authorization", "Bearer my test key")

	a.EXPECT().ValidateUserToken(
		"my test key",
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(tokenString string, user *string, organization *string, role *int) bool {
		if tokenString != "my test key" {
			return false
		}
		if user != nil {
			*user = "user456"
		}
		if organization != nil {
			*organization = "org123"
		}
		if role != nil {
			*role = int(types.RoleUser | types.RoleMaintainer | types.RoleAdministrator | types.RoleRoot)
		}
		return true
	})

	d.EXPECT().ReadTask("task1234").Return(&db.Task{
		Model: gorm.Model{
			ID: 1234,
		},
		Name:           "task1234",
		Description:    "test task",
		Content:        db.StringToJSON(`{"definition": "content schema is not defined yet so just use placeholder"}`),
		OrganizationID: 123,
	}, nil)
	d.EXPECT().ReadOrganization("org123").Return(&db.Organization{
		Model: gorm.Model{
			ID: 123,
		},
		Name: "org123",
	}, nil).Times(2)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error doing request:", err)
	}
	defer resp.Body.Close()

	var response Response
	b, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(b, &response); err != nil {
		t.Fatal("failed to decode response as JSON: \"", err, "\", response was: \"", string(b), "\"")
	}

	if response.Code != 200 {
		t.Fatal("response not 200:", response)
	}
}

func TestProcessRequestGetTasks(t *testing.T) {
	ctrl := gomock.NewController(t)

	a := mock_auth.NewMockController(ctrl)
	d := mock_db.NewMockController(ctrl)
	h := NewHandler(a, d)
	if h == nil {
		t.Fatal("nil handler created")
	}

	if err := h.RegisterUserHandlers(); err != nil {
		t.Fatal("error registering user handlers:", err)
	}
	if err := h.RegisterOrganizationHandlers(); err != nil {
		t.Fatal("error registering organization handlers:", err)
	}
	if err := h.RegisterMachineHandlers(); err != nil {
		t.Fatal("error registering machine handlers:", err)
	}
	if err := h.RegisterTaskHandlers(); err != nil {
		t.Fatal("error registering task handlers:", err)
	}
	if err := h.RegisterScheduleHandlers(); err != nil {
		t.Fatal("error registering schedule handlers:", err)
	}
	if err := h.RegisterAuthenticationHandlers(); err != nil {
		t.Fatal("error registering authentication handlers:", err)
	}
	if err := h.RegisterRecordHandlers(); err != nil {
		t.Fatal("error registering record handlers:", err)
	}

	server := httptest.NewServer(h)

	if server == nil {
		t.Fatal("failed to create test server")
	}

	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/org123/tasks/", nil)
	req.Header.Set("Authorization", "Bearer my test key")

	a.EXPECT().ValidateUserToken(
		"my test key",
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(tokenString string, user *string, organization *string, role *int) bool {
		if tokenString != "my test key" {
			return false
		}
		if user != nil {
			*user = "user456"
		}
		if organization != nil {
			*organization = "org123"
		}
		if role != nil {
			*role = int(types.RoleUser | types.RoleMaintainer | types.RoleAdministrator | types.RoleRoot)
		}
		return true
	})

	task1234 := db.Task{
		Model: gorm.Model{
			ID: 1234,
		},
		Name:           "task1234",
		Description:    "test task",
		Content:        db.StringToJSON(`{"definition": "content schema is not defined yet so just use placeholder"}`),
		OrganizationID: 123,
	}
	task5678 := db.Task{
		Model: gorm.Model{
			ID: 5678,
		},
		Name:           "task5678",
		Description:    "test task",
		Content:        db.StringToJSON(`{"definition": "content schema is not defined yet so just use placeholder"}`),
		OrganizationID: 123,
	}

	d.EXPECT().ReadTask("task1234").Return(&task1234, nil)
	d.EXPECT().ReadTask("task5678").Return(&task5678, nil)
	d.EXPECT().ReadOrganization("org123").Return(&db.Organization{
		Model: gorm.Model{
			ID: 123,
		},
		Name:  "org123",
		Tasks: []db.Task{task1234, task5678},
	}, nil).Times(2)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error doing request:", err)
	}
	defer resp.Body.Close()

	var response Response
	b, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(b, &response); err != nil {
		t.Fatal("failed to decode response as JSON: \"", err, "\", response was: \"", string(b), "\"")
	}

	if response.Code != 200 {
		t.Fatal("response not 200:", response)
	}
}
