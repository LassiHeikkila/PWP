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

func (*handler) deleteMachineToken(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}
