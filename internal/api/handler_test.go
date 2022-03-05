package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/auth/mock"
	"github.com/LassiHeikkila/taskey/internal/db/mock"
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

	server := httptest.NewServer(h)

	if server == nil {
		t.Fatal("failed to create test server")
	}

	client := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/org123/users/user456/", nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error doing request:", err)
	}
	defer resp.Body.Close()

	response := make(map[string]interface{})
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&response); err != nil {
		t.Fatal("failed to decode response as JSON:", err)
	}
}
