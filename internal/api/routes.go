package api

import (
	"net/http"
)

/*
   ${base}/api/v1.0/auth -> auth service
   ${base}/api/v1.0/signup -> signup service
   ${base}/api/v1.0/${org}/users -> user management
   ${base}/api/v1.0/${org}/machines -> machine management
   ${base}/api/v1.0/${org}/machines/${machine}/schedule -> control machine schedule
   ${base}/api/v1.0/${org}/machines/${machine}/records/ -> get and post machine records

*/

func (h *handler) setOrgRoutesV1() {
	// create organization
	h.router.HandleFunc("/api/v1/organizations/", h.createOrganization).Methods(http.MethodPost)
	// read organization
	h.router.Handle("/api/v1/organizations/{organization_id}/", h.requiresAdmin(h.readOrganization)).Methods(http.MethodGet)
	// update organization
	h.router.Handle("/api/v1/organizations/{organization_id}/", h.requiresAdmin(h.updateOrganization)).Methods(http.MethodPut)
	// delete organization
	h.router.Handle("/api/v1/organizations/{organization_id}/", h.requiresRoot(h.deleteOrganization)).Methods(http.MethodDelete)
}

func (h *handler) setUserRoutesV1() {
	// create user
	h.router.Handle("/api/v1/{organization_id}/users/", h.requiresAdmin(h.createUser)).Methods(http.MethodPost)
	// read users
	h.router.Handle("/api/v1/{organization_id}/users/", h.requiresAdmin(h.readUsers)).Methods(http.MethodGet)
	// read user
	h.router.Handle("/api/v1/{organization_id}/users/{user_id}/", h.requiresUser(h.readUser)).Methods(http.MethodGet)
	// update user
	h.router.Handle("/api/v1/{organization_id}/users/{user_id}/", h.requiresAdmin(h.updateUser)).Methods(http.MethodPut)
	// delete user
	h.router.Handle("/api/v1/{organization_id}/users/{user_id}/", h.requiresAdmin(h.deleteUser)).Methods(http.MethodDelete)
	// create token
	h.router.Handle("/api/v1/{organization_id}/users/{user_id}/tokens/", h.requiresAdmin(h.createUserToken)).Methods(http.MethodPost)
	// delete / revoke token
	h.router.Handle("/api/v1/{organization_id}/users/{user_id}/tokens/{token}/", h.requiresAdmin(h.deleteUserToken)).Methods(http.MethodDelete)
}

func (h *handler) setMachineRoutesV1() {
	// create, read, update and delete machine(s)
	h.router.Handle("/api/v1/{organization_id}/machines/", h.requiresAdmin(h.createMachine)).Methods(http.MethodPost)
	h.router.Handle("/api/v1/{organization_id}/machines/", h.requiresAdmin(h.readMachines)).Methods(http.MethodGet)
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/", h.requiresAdmin(h.readMachine)).Methods(http.MethodGet)
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/", h.requiresAdmin(h.updateMachine)).Methods(http.MethodPut)
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/", h.requiresAdmin(h.deleteMachine)).Methods(http.MethodDelete)

	// create token
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/tokens/", h.requiresAdmin(h.createMachineToken)).Methods(http.MethodPost)
	// delete token
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/tokens/{token}/", h.requiresAdmin(h.deleteMachineToken)).Methods(http.MethodDelete)
}

func (h *handler) setScheduleRoutesV1() {
	// create, read, update or delete machine schedule
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/schedule/", h.requiresUser(h.createMachineSchedule)).Methods(http.MethodPost)
	// TODO: think about how to route this so machine can also get their own schedule
	// or should there be a separate endpoint like /api/v1/{organization_id}/machines/self/schedule/
	// and handler will identify machine based on token ?
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/schedule/", h.requiresUser(h.readMachineSchedule)).Methods(http.MethodGet)
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/schedule/", h.requiresUser(h.updateMachineSchedule)).Methods(http.MethodPut)
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/schedule/", h.requiresUser(h.deleteMachineSchedule)).Methods(http.MethodDelete)
}

func (h *handler) setRecordRoutesV1() {
	// create and read records
	// TODO: this about adding machine member middleware
	// TODO: also think about adding endpoint where machine can just post record without specifying its own id
	// e.g. POST /api/v1/{organization_id}/addrecord/
	// and server will figure out who it came from based on Authorization header
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/records/", h.requiresMachine(h.addRecord)).Methods(http.MethodPost)
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/records/", h.requiresUser(h.readRecords)).Methods(http.MethodGet)

	// records are immutable so modifying them via PUT is not allowed

	// get and delete a particular record
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/records/{record_id}/", h.requiresUser(h.readRecord)).Methods(http.MethodGet)
	h.router.Handle("/api/v1/{organization_id}/machines/{machine_id}/records/{record_id}/", h.requiresAdmin(h.deleteRecord)).Methods(http.MethodDelete)
}

func (h *handler) setTaskRoutesV1() {
	// create, read, update and delete tasks
	h.router.Handle("/api/v1/{organization_id}/tasks/", h.requiresUser(h.createTask)).Methods(http.MethodPost)
	h.router.Handle("/api/v1/{organization_id}/tasks/", h.requiresUser(h.readTasks)).Methods(http.MethodGet)
	h.router.Handle("/api/v1/{organization_id}/tasks/{task_id}/", h.requiresUser(h.readTask)).Methods(http.MethodGet)
	h.router.Handle("/api/v1/{organization_id}/tasks/{task_id}/", h.requiresUser(h.updateTask)).Methods(http.MethodPut)
	h.router.Handle("/api/v1/{organization_id}/tasks/{task_id}/", h.requiresUser(h.deleteTask)).Methods(http.MethodDelete)
	// TODO: this about how machines will retrieve task definitions
}

func (h *handler) setAuthRoutesV1() {
	// login with username & password or token, get JWT
	h.router.HandleFunc("/api/v1/auth/", h.loginHandler).Methods(http.MethodPost)
	// check if JWT is OK
	h.router.HandleFunc("/api/v1/auth/", h.loginChecker).Methods(http.MethodGet)
	// change password
	h.router.HandleFunc("/api/v1/auth/{username}/changepassword/", h.passwordChangeHandler).Methods(http.MethodPost)
}
