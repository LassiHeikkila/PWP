package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/taskey/internal/db/dbconverter"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

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

	rid, err := strconv.ParseUint(recordID, 10, 64)
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

	if err := h.d.DeleteRecord(m.Name, rid); err != nil {
		_ = encodeFailure(w)
		return
	}

	_ = encodeSuccess(w)
}
