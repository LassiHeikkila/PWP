package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/db/dbconverter"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func (h *handler) readMachineOwnSchedule(w http.ResponseWriter, req *http.Request, self *types.Machine) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	m, err := h.d.ReadMachine(self.Name)
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

func (h *handler) addRecord(w http.ResponseWriter, req *http.Request, self *types.Machine) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	m, err := h.d.ReadMachine(self.Name)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	if m.OrganizationID != o.ID {
		_ = encodeNotFoundResponse(w)
		return
	}

	var reqRecord types.Record
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqRecord); err != nil {
		_ = encodeBadRequestResponse(w)
		return
	}

	t, err := h.d.ReadTask(reqRecord.TaskName)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}

	record := dbconverter.ConvertRecordToDB(&reqRecord)

	record.MachineID = m.ID
	record.TaskID = t.ID

	if err := h.d.CreateRecord(&record); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
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
	// let's keep it simple for now and just return all task definitions in the same organization

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
