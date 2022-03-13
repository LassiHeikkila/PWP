package api

import (
	"html"
	"strings"
)

const (
	orgIDKey     = "organization_id"
	userIDKey    = "user_id"
	machineIDKey = "machine_id"
	recordIDKey  = "record_id"
	taskIDKey    = "task_id"
	tokenKey     = "token"
)

func sanitizeParameter(input string) string {
	output := strings.ReplaceAll(input, "\n", "")
	output = strings.ReplaceAll(output, "\r", "")
	output = strings.ReplaceAll(output, "\t", "")
	output = strings.ReplaceAll(output, " ", "")
	output = html.EscapeString(output)

	return output
}
