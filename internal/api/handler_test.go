package api

import (
	"net/http"
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
