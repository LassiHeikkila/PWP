package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/auth"
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

// AuthUserMW implements middleware pattern.
// It should be chained to match routes requiring a specific role.
// It also checks that caller is a member of the organization owning the resource
type AuthUserMW struct {
	handler func(http.ResponseWriter, *http.Request)

	authController auth.Controller
	dbController   db.Controller
	requiredRole   types.Role
}

func NewAuthUserMiddleware(
	next func(http.ResponseWriter, *http.Request),
	authController auth.Controller,
	dbController db.Controller,
	requiredRole types.Role,
) *AuthUserMW {
	return &AuthUserMW{
		handler:        next,
		authController: authController,
		dbController:   dbController,
		requiredRole:   requiredRole,
	}
}

func (a *AuthUserMW) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// check that token is valid, contains a legit user & organization, and role is what is required
	scheme, value := auth.GetAuthenticationSchemeAndValue(req.Header.Get("Authorization"))
	var user *types.User
	if scheme == auth.AuthenticationSchemeBearer {
		// validate token and read user info from it,
		// if its valid, call handler with the User
		user = lookupUserByJWT(a.authController, value)
	}
	if scheme == auth.AuthenticationSchemeKey {
		// look up token in db,
		// if it exists and isn't revoked
		// call handler with corresponding User
		user = lookupUserByToken(a.dbController, value)
	}
	if user == nil {
		_ = encodeUnauthenticatedResponse(w)
		return
	}

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	if orgID != "" {
		// path contains organization_id,
		// so we need to check that the caller is a member of that organization
		org, err := a.dbController.ReadOrganization(orgID)
		if err != nil {
			_ = encodeNotFoundResponse(w)
			return
		}

		usr, err := a.dbController.ReadUser(user.Name)
		if err != nil {
			_ = encodeNotFoundResponse(w)
			return
		}

		if usr.OrganizationID != org.ID {
			_ = encodeForbiddenResponse(w)
			return
		}
	}

	if !types.HasRole(user.Role, a.requiredRole) {
		_ = encodeForbiddenResponse(w)
		return
	}
	a.handler(w, req)
}

type AuthMachineMW struct {
	handler AuthenticatedMachineHandler

	authController auth.Controller
	dbController   db.Controller
}

func NewAuthMachineMiddleware(
	next AuthenticatedMachineHandler,
	authController auth.Controller,
	dbController db.Controller,
) *AuthMachineMW {
	return &AuthMachineMW{
		handler:        next,
		authController: authController,
		dbController:   dbController,
	}
}

func (a *AuthMachineMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check that token is valid, contains a legit machine key
	scheme, value := auth.GetAuthenticationSchemeAndValue(r.Header.Get("Authorization"))
	var machine *types.Machine
	if scheme == auth.AuthenticationSchemeKey {
		// look up token in db,
		// if it exists and isn't revoked
		// call handler with corresponding Machine
		machine = lookupMachineByToken(a.dbController, value)
	}
	if machine == nil {
		_ = encodeUnauthenticatedResponse(w)
		return
	}
	a.handler(w, r, machine)
}

func (h *handler) requiresAdmin(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return NewAuthUserMiddleware(next, h.a, h.d, types.RoleAdministrator)
}

func (h *handler) requiresMaintainer(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return NewAuthUserMiddleware(next, h.a, h.d, types.RoleMaintainer)
}

func (h *handler) requiresRoot(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return NewAuthUserMiddleware(next, h.a, h.d, types.RoleRoot)
}

func (h *handler) requiresUser(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return NewAuthUserMiddleware(next, h.a, h.d, types.RoleUser)
}

func (h *handler) requiresMachine(next AuthenticatedMachineHandler) http.Handler {
	return NewAuthMachineMiddleware(next, h.a, h.d)
}
