package api

import (
	"net/http"

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
	requiredRole   int
}

func NewAuthUserMiddleware(
	next AuthenticatedUserHandler,
	authController auth.Controller,
	dbController db.Controller,
	requiredRole int,
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
	scheme := auth.GetAuthenticationScheme(r.Header.Get("Authorization"))
	if scheme == auth.AuthenticationSchemeNone {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	a.handler(w, r, &types.User{})
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
	a.handler(w, r, &types.Machine{})
}
