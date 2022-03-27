package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/internal/db/dbconverter"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

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
	defer req.Body.Close()

	vars := mux.Vars(req)
	orgID := sanitizeParameter(vars[orgIDKey])
	taskID := sanitizeParameter(vars[taskIDKey])

	o, err := h.d.ReadOrganization(orgID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	t, err := h.d.ReadTask(taskID)
	if err != nil {
		_ = encodeNotFoundResponse(w)
		return
	}
	if t.OrganizationID != o.ID {
		_ = encodeNotFoundResponse(w)
		return
	}

	var reqTask types.Task
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqTask); err != nil {
		_ = encodeBadRequestResponse(w)
		return
	}

	t.Name = reqTask.Name
	t.Description = reqTask.Description
	b, _ := json.Marshal(&reqTask.Content)
	t.Content = db.StringToJSON(string(b))

	if err := h.d.UpdateTask(t); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
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
