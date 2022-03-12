package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/db/dbconverter"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func (h *handler) signupHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) createOrganization(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) readOrganization(w http.ResponseWriter, req *http.Request, requester *types.User) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := vars[orgIDKey]

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	org := dbconverter.ConvertOrganization(*o)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &org,
	})
}

func (h *handler) updateOrganization(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteOrganization(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) createUser(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) readUsers(w http.ResponseWriter, req *http.Request, requester *types.User) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := vars[orgIDKey]

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	users := make([]types.User, 0, len(o.Users))
	for _, usr := range o.Users {
		u, err := h.d.ReadUser(usr.Name)
		if err != nil {
			_ = encodeNotFoundResponse(w)
			continue
		}

		user := dbconverter.ConvertUser(*u)
		users = append(users, user)
	}

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &users,
	})
}

// /api/v1/{organization_id}/users/{user_id}/
func (h *handler) readUser(w http.ResponseWriter, req *http.Request, requester *types.User) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := vars[orgIDKey]
	userID := vars[userIDKey]

	u, err := h.d.ReadUser(userID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	if u.OrganizationID != o.ID {
		_ = encodeNotFoundResponse(w)
		return
	}

	user := dbconverter.ConvertUser(*u)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &user,
	})
}

func (h *handler) updateUser(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteUser(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) createUserToken(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteUserToken(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) createMachine(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) readMachine(w http.ResponseWriter, req *http.Request, requester *types.User) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := vars[orgIDKey]
	machineID := vars[machineIDKey]

	m, err := h.d.ReadMachine(machineID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	if m.OrganizationID != o.ID {
		_ = encodeNotFoundResponse(w)
		return
	}

	machine := dbconverter.ConvertMachine(*m)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &machine,
	})
}

func (h *handler) updateMachine(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteMachine(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) createMachineToken(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteMachineToken(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) createMachineSchedule(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) readMachineSchedule(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) updateMachineSchedule(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteMachineSchedule(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) addRecord(w http.ResponseWriter, req *http.Request, machine *types.Machine) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) getRecord(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) getRecords(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteRecord(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) createTask(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) readTask(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) readTasks(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) updateTask(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteTask(w http.ResponseWriter, req *http.Request, requester *types.User) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) loginHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) loginChecker(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) passwordChangeHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}
