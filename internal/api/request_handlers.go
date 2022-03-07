package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/db/dbconverter"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

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

// /api/v1/{organization_id}/users/{user_id}/
func (h *handler) readUser(w http.ResponseWriter, req *http.Request, caller *types.User) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := vars[orgIDKey]
	userID := vars[userIDKey]

	u, err := h.d.ReadUser(userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if u.OrganizationID != o.ID {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	user := dbconverter.ConvertUser(*u)

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	enc.Encode(user)
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
