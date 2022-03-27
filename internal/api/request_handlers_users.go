package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/internal/db/dbconverter"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func (h *handler) createUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	var reqUser types.User
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqUser); err != nil {
		_ = encodeBadRequestResponse(w)
		return
	}

	user := dbconverter.ConvertUserToDB(&reqUser)
	user.OrganizationID = o.ID

	if err := h.d.CreateUser(&user); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) readUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	userID := sanitizeParameter(vars[userIDKey])

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

	user := dbconverter.ConvertUser(u)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &user,
	})
}

func (h *handler) readUsers(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	users := make([]types.User, 0, len(o.Users))
	for i := range o.Users {
		u, err := h.d.ReadUser(o.Users[i].Name)
		if err != nil {
			continue
		}

		user := dbconverter.ConvertUser(u)
		users = append(users, user)
	}

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &users,
	})
}

func (h *handler) updateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	userID := sanitizeParameter(vars[userIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	u, err := h.d.ReadUser(userID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	if u.OrganizationID != o.ID {
		_ = encodeNotFoundResponse(w)
		return
	}

	var reqUser types.User
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqUser); err != nil {
		_ = encodeBadRequestResponse(w)
		return
	}

	u.Name = reqUser.Name
	u.Email = reqUser.Email
	u.Role = reqUser.Role

	if err := h.d.UpdateUser(u); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) deleteUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	userID := sanitizeParameter(vars[userIDKey])

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

	if err := h.d.DeleteUser(u.Name); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) createUserToken(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	userID := sanitizeParameter(vars[userIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	u, err := h.d.ReadUser(userID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	if u.OrganizationID != o.ID {
		_ = encodeNotFoundResponse(w)
		return
	}

	genUUID, err := h.a.GenerateUUID()
	if err != nil {
		_ = encodeFailure(w)
		return
	}

	// TODO: add option to define token expiration via a body
	expiration := time.Time{} // zero time means no expiry

	ut := db.UserToken{
		Value:      db.StringToUUID(genUUID),
		Expiration: expiration,
		UserID:     u.ID,
		User:       *u,
	}

	if err := h.d.CreateUserToken(&ut); err != nil {
		_ = encodeFailure(w)
		return
	}

	returnedToken := dbconverter.ConvertUserToken(&ut)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &returnedToken,
	})
}

func (*handler) deleteUserToken(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}
