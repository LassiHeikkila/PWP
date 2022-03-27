package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/db/dbconverter"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func (h *handler) readMachineOwnSchedule(w http.ResponseWriter, req *http.Request, _ *types.Machine) {
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

func (*handler) addRecord(w http.ResponseWriter, req *http.Request, _ *types.Machine) {
	// TODO: implement
	defer req.Body.Close()
	_ = encodeUnimplementedResponse(w)
}

func (h *handler) readMachineTasks(w http.ResponseWriter, req *http.Request, _ *types.Machine) {
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
