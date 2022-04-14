package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/auth"
	"github.com/LassiHeikkila/taskey/internal/db"
)

type Handler interface {
	http.Handler

	RegisterOrganizationHandlers() error
	RegisterUserHandlers() error
	RegisterMachineHandlers() error
	RegisterScheduleHandlers() error
	RegisterTaskHandlers() error
	RegisterAuthenticationHandlers() error
}

type handler struct {
	router *mux.Router

	a auth.Controller
	d db.Controller
}

func NewHandler(a auth.Controller, d db.Controller) *handler {
	m := mux.NewRouter()

	return &handler{
		router: m,
		a:      a,
		d:      d,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.router == nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	if r.Method == http.MethodOptions {
		return
	}
	h.router.ServeHTTP(w, r)
}

func (h *handler) RegisterOrganizationHandlers() error {
	h.setOrgRoutesV1()
	return nil
}

func (h *handler) RegisterUserHandlers() error {
	h.setUserRoutesV1()
	return nil
}

func (h *handler) RegisterMachineHandlers() error {
	h.setMachineRoutesV1()
	return nil
}

func (h *handler) RegisterScheduleHandlers() error {
	h.setScheduleRoutesV1()
	return nil
}

func (h *handler) RegisterTaskHandlers() error {
	h.setTaskRoutesV1()
	return nil
}

func (h *handler) RegisterAuthenticationHandlers() error {
	h.setAuthRoutesV1()
	return nil
}

func (h *handler) RegisterRecordHandlers() error {
	h.setRecordRoutesV1()
	return nil
}

func (h *handler) RegisterSignUpHandlers() error {
	h.setSignUpRoutesV1()
	return nil
}

func (h *handler) RegisterExtraRoute(path string, handlerFunc func(http.ResponseWriter, *http.Request)) error {
	h.router.HandleFunc(path, handlerFunc)
	return nil
}
