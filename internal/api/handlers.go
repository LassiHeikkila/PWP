package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/auth"
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

type AuthenticatedUserHandler func(http.ResponseWriter, *http.Request, *types.User)
type AuthenticatedMachineHandler func(http.ResponseWriter, *http.Request, *types.Machine)

type handler struct {
	router *mux.Router

	a auth.Controller
	d db.Controller
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.router == nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	h.router.ServeHTTP(w, r)
}

func (h *handler) signupHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) createOrganization(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) readOrganization(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) updateOrganization(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) deleteOrganization(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) createUser(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) readUser(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) updateUser(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) deleteUser(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) createUserToken(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) deleteUserToken(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) createMachine(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) readMachine(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) updateMachine(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) deleteMachine(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) createMachineToken(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) deleteMachineToken(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) createMachineSchedule(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) readMachineSchedule(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) updateMachineSchedule(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) deleteMachineSchedule(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) addRecord(w http.ResponseWriter, req *http.Request, machine *types.Machine) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) getRecord(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) getRecords(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) deleteRecord(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) createTask(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) readTask(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) readTasks(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) updateTask(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) deleteTask(w http.ResponseWriter, req *http.Request, user *types.User) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) loginHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) loginChecker(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (h *handler) passwordChangeHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}
