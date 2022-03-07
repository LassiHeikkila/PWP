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
type AuthUserMW struct {
	handler AuthenticatedUserHandler

	authController auth.Controller
	dbController   db.Controller
	requiredRole   types.Role
}

func NewAuthUserMiddleware(
	next AuthenticatedUserHandler,
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

func (a *AuthUserMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check that token is valid, contains a legit user & organization, and role is what is required
	scheme, value := auth.GetAuthenticationSchemeAndValue(r.Header.Get("Authorization"))
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
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if !types.HasRole(user.Role, a.requiredRole) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	a.handler(w, r, user)
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
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	a.handler(w, r, machine)
}

// should be used when accessing resources belonging to an organization
type MemberOfOrganizationMW struct {
	handler AuthenticatedUserHandler

	dbController db.Controller
}

func NewMemberOfOrganizationMW(
	next AuthenticatedUserHandler,
	dbController db.Controller,
) *MemberOfOrganizationMW {
	return &MemberOfOrganizationMW{
		handler:      next,
		dbController: dbController,
	}
}

func (m *MemberOfOrganizationMW) ServeHTTP(
	w http.ResponseWriter,
	req *http.Request,
	user *types.User,
) {
	vars := mux.Vars(req)
	orgID := vars[orgIDKey]

	org := lookupOrganizationByID(m.dbController, orgID)
	if org == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if user.Organization != org.Name {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	m.handler(w, req, user)
}

func (h *handler) requiresAdmin(next AuthenticatedUserHandler) http.Handler {
	return NewAuthUserMiddleware(next, h.a, h.d, types.RoleAdministrator)
}

func (h *handler) requiresMaintainer(next AuthenticatedUserHandler) http.Handler {
	return NewAuthUserMiddleware(next, h.a, h.d, types.RoleMaintainer)
}

func (h *handler) requiresRoot(next AuthenticatedUserHandler) http.Handler {
	return NewAuthUserMiddleware(next, h.a, h.d, types.RoleRoot)
}

func (h *handler) requiresUser(next AuthenticatedUserHandler) http.Handler {
	return NewAuthUserMiddleware(next, h.a, h.d, types.RoleUser)
}

func (h *handler) requiresMachine(next AuthenticatedMachineHandler) http.Handler {
	return NewAuthMachineMiddleware(next, h.a, h.d)
}

func (h *handler) mustBeMember(next AuthenticatedUserHandler) AuthenticatedUserHandler {
	mw := NewMemberOfOrganizationMW(next, h.d)
	return mw.ServeHTTP
}
