package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/auth"
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/internal/db/dbconverter"
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

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	org := dbconverter.ConvertOrganization(o)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &org,
	})
}

func (h *handler) updateOrganization(w http.ResponseWriter, req *http.Request) {
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
	user.Organization = *o

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
	if u.Organization.ID != o.ID {
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

func (h *handler) deleteUserToken(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) createMachine(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	var reqMachine types.Machine
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqMachine); err != nil {
		_ = encodeBadRequestResponse(w)
		return
	}

	machine := dbconverter.ConvertMachineToDB(&reqMachine)
	machine.OrganizationID = o.ID

	if err := h.d.CreateMachine(&machine); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) readMachine(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])

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

	machine := dbconverter.ConvertMachine(m)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &machine,
	})
}

func (h *handler) readMachines(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	machines := make([]types.Machine, 0, len(o.Machines))
	for i := range o.Machines {
		mchn, err := h.d.ReadMachine(o.Machines[i].Name)
		if err != nil {
			continue
		}

		machine := dbconverter.ConvertMachine(mchn)
		machines = append(machines, machine)
	}

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &machines,
	})
}

func (h *handler) updateMachine(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])

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

	dec := json.NewDecoder(req.Body)
	var updatedMachine types.Machine
	if err := dec.Decode(&updatedMachine); err != nil {
		_ = encodeBadRequestResponse(w)
		return
	}

	m.Name = updatedMachine.Name
	m.Description = updatedMachine.Description
	m.OS = updatedMachine.OS
	m.Arch = updatedMachine.Arch

	if err := h.d.UpdateMachine(m); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) deleteMachine(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])

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

	if err := h.d.DeleteMachine(m.Name); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) createMachineToken(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	m, err := h.d.ReadMachine(machineID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	if m.OrganizationID != o.ID {
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

	mt := db.MachineToken{
		Value:      db.StringToUUID(genUUID),
		Expiration: expiration,
		MachineID:  m.ID,
		Machine:    *m,
	}

	if err := h.d.CreateMachineToken(&mt); err != nil {
		_ = encodeFailure(w)
		return
	}

	returnedToken := dbconverter.ConvertMachineToken(&mt)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &returnedToken,
	})
}

func (h *handler) deleteMachineToken(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) createMachineSchedule(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) readMachineSchedule(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])

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

	sched, err := h.d.ReadSchedule(m.Name)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	schedule := dbconverter.ConvertSchedule(sched)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &schedule,
	})
}

func (h *handler) updateMachineSchedule(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])

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

	var reqSched types.Schedule
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqSched); err != nil {
		_ = encodeFailure(w)
		return
	}

	sched, err := h.d.ReadSchedule(m.Name)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	schedule := dbconverter.ConvertSchedule(sched)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &schedule,
	})
}

func (h *handler) deleteMachineSchedule(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])

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

	if err := h.d.DeleteSchedule(m.Name); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) addRecord(w http.ResponseWriter, req *http.Request, machine *types.Machine) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) readRecord(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])
	recordID := sanitizeParameter(vars[recordIDKey])

	rid, err := strconv.ParseInt(recordID, 10, 64)
	if err != nil {
		_ = encodeBadRequestResponse(w)
	}

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

	r, err := h.d.ReadRecords(machineID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	for i := range r {
		rec := r[i]
		if int64(rec.ID) == rid {
			record := dbconverter.ConvertRecord(&rec)
			_ = encodeResponse(w, Response{
				Code:    http.StatusOK,
				Message: "ok",
				Payload: &record,
			})
			return
		}
	}
	_ = encodeNotFoundResponse(w)
}

func (h *handler) readRecords(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])

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

	r, err := h.d.ReadRecords(machineID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	records := make([]types.Record, 0, len(r))
	for i := range r {
		record := dbconverter.ConvertRecord(&r[i])
		records = append(records, record)
	}

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &records,
	})
}

func (h *handler) deleteRecord(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	machineID := sanitizeParameter(vars[machineIDKey])
	recordID := sanitizeParameter(vars[recordIDKey])
	rid, err := strconv.ParseInt(recordID, 10, 64)
	if err != nil {
		_ = encodeBadRequestResponse(w)
	}

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

	if err := h.d.DeleteRecord(m.Name, uint(rid)); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) createTask(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	var reqTask types.Task
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqTask); err != nil {
		_ = encodeBadRequestResponse(w)
		return
	}

	task := dbconverter.ConvertTaskToDB(&reqTask)
	task.OrganizationID = o.ID

	if err := h.d.CreateTask(&task); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}

func (h *handler) readTask(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	taskID := sanitizeParameter(vars[taskIDKey])

	tsk, err := h.d.ReadTask(taskID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	if tsk.OrganizationID != o.ID {
		_ = encodeNotFoundResponse(w)
		return
	}

	task := dbconverter.ConvertTask(tsk)

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &task,
	})
}

func (h *handler) readTasks(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	tasks := make([]types.Task, 0, len(o.Tasks))
	for i := range o.Tasks {
		t, err := h.d.ReadTask(o.Tasks[i].Name)
		if err != nil {
			continue
		}

		task := dbconverter.ConvertTask(t)
		tasks = append(tasks, task)
	}

	_ = encodeResponse(w, Response{
		Code:    http.StatusOK,
		Message: "ok",
		Payload: &tasks,
	})
}

func (h *handler) updateTask(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) deleteTask(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	taskID := sanitizeParameter(vars[taskIDKey])

	tsk, err := h.d.ReadTask(taskID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	if tsk.OrganizationID != o.ID {
		_ = encodeNotFoundResponse(w)
		return
	}

	if err := h.d.DeleteTask(tsk.Name); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
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
