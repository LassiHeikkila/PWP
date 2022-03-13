package docs

import (
	_ "embed"
)

//go:embed openapi.yml
var OpenAPI []byte
