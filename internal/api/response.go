package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Payload interface{} `json:"payload,omitempty"`
}

func encodeResponse(w http.ResponseWriter, r Response) error {
	enc := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")

	if r.Code != http.StatusOK {
		w.WriteHeader(r.Code)
	}

	return enc.Encode(&r)
}

func encodeSuccess(w http.ResponseWriter) error {
	return encodeResponse(w, Response{Code: http.StatusOK, Message: "ok"})
}

func encodeNotFoundResponse(w http.ResponseWriter) error {
	return encodeResponse(w, Response{Code: http.StatusNotFound, Message: "not found"})
}

func encodeForbiddenResponse(w http.ResponseWriter) error {
	return encodeResponse(w, Response{Code: http.StatusForbidden, Message: "forbidden"})
}

func encodeUnauthenticatedResponse(w http.ResponseWriter) error {
	return encodeResponse(w, Response{Code: http.StatusUnauthorized, Message: "unauthorized"})
}

func encodeUnimplementedResponse(w http.ResponseWriter) error {
	return encodeResponse(w, Response{Code: http.StatusNotImplemented, Message: "not implemented yet"})
}
