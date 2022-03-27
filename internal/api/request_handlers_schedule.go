package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/db/dbconverter"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func (h *handler) createMachineSchedule(w http.ResponseWriter, req *http.Request) {
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

	var reqSchedule types.Schedule
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqSchedule); err != nil {
		_ = encodeBadRequestResponse(w)
		return
	}

	schedule := dbconverter.ConvertScheduleToDB(&reqSchedule)
	schedule.MachineID = m.ID

	if err := h.d.CreateSchedule(&schedule); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
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
