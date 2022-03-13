package main

import (
	_ "embed"
	"net/http"

    "github.com/LassiHeikkila/taskey/docs"
)

//go:embed swaggerui.html
var swaggerui []byte

func ServeSwaggerUI(w http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()
    _, _ = w.Write(swaggerui)
}


func ServeOpenAPI(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	_, _ = w.Write(docs.OpenAPI)
}
