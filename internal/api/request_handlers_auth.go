package api

import (
	"encoding/json"
	"net/http"

	"github.com/LassiHeikkila/taskey/internal/auth"
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func (h *handler) loginHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var loginRequest types.LoginInfo
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&loginRequest); err != nil {
		_ = encodeUnauthenticatedResponse(w)
		return
	}

	loginInfo, err := h.d.ReadLoginInfo(loginRequest.Username)
	if err != nil {
		_ = encodeUnauthenticatedResponse(w)
		return
	}

	if !auth.PasswordEqualsHashed(loginRequest.Password, loginInfo.Password) {
		_ = encodeUnauthenticatedResponse(w)
		return
	}

	// look up organization by ID
	org := db.Organization{}
	if err := h.d.LoadModel(&org, loginInfo.User.OrganizationID); err != nil {
		_ = encodeUnauthenticatedResponse(w)
		return
	}

	token, err := h.a.CreateJWT(auth.CreateUserClaims(
		loginInfo.User.Name,
		org.Name,
		int(loginInfo.User.Role),
	))
	if err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: map[string]interface{}{
			"token": token,
		},
	})
}

func (h *handler) loginChecker(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	// check that token is valid, contains a legit user & organization, and role is what is required
	scheme, value := auth.GetAuthenticationSchemeAndValue(req.Header.Get("Authorization"))
	var user *types.User
	if scheme == auth.AuthenticationSchemeBearer {
		// validate token and read user info from it,
		// if its valid, call handler with the User
		user = lookupUserByJWT(h.a, value)
	}
	if scheme == auth.AuthenticationSchemeKey {
		// look up token in db,
		// if it exists and isn't revoked
		// call handler with corresponding User
		user = lookupUserByToken(h.d, value)
	}
	if user == nil {
		_ = encodeUnauthenticatedResponse(w)
		return
	}
	_ = encodeSuccess(w)
}

func (*handler) passwordChangeHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}
