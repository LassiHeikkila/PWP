package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/auth"
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func (h *handler) signupHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var signupRequest SignUpRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&signupRequest); err != nil || !validateSignUpRequest(&signupRequest) {
		_ = encodeFailure(w)
		return
	}

	org := db.Organization{
		Name: signupRequest.OrganizationName,
	}
	if err := h.d.CreateOrganization(&org); err != nil {
		log.Println("failed to create organization:", err)
		_ = encodeFailure(w)
		return
	}

	user := db.User{
		Name:           signupRequest.Username,
		Email:          signupRequest.Email,
		Role:           types.RoleRoot | types.RoleAdministrator | types.RoleMaintainer | types.RoleUser,
		OrganizationID: org.ID,
	}
	if err := h.d.CreateUser(&user); err != nil {
		log.Println("failed to create user:", err)
		_ = encodeFailure(w)
		return
	}

	li := db.LoginInfo{
		Username: signupRequest.Username,
		Password: auth.HashPassword(signupRequest.Password),
		UserID:   user.ID,
	}
	if err := h.d.CreateLoginInfo(&li); err != nil {
		log.Println("failed to create login info:", err)
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) readOrganization(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	org := lookupOrganizationByID(h.d, orgID)
	if org == nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &org,
	})
}

func (*handler) updateOrganization(w http.ResponseWriter, req *http.Request) {
	// TODO: implement when / if organization has more details than name
	// with just name property, it doesn't make sense to implement anything
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteOrganization(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	if err := h.d.DeleteOrganization(o.Name); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}
