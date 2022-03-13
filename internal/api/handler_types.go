package api

import (
	"net/http"

	"github.com/LassiHeikkila/taskey/pkg/types"
)

//type AuthenticatedUserHandler func(http.ResponseWriter, *http.Request, *types.User)
type AuthenticatedMachineHandler func(http.ResponseWriter, *http.Request, *types.Machine)
