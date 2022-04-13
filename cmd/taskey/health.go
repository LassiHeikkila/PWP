package main

import (
	"encoding/json"
	"net/http"
)

func serveHealth(w http.ResponseWriter, _ *http.Request) {
	m := map[string]interface{}{
		"ok": true,
	}

	enc := json.NewEncoder(w)
	_ = enc.Encode(&m)
}
