package auth

import (
	"strings"
)

type AuthenticationScheme int8

const (
	AuthenticationSchemeNone   AuthenticationScheme = 0
	AuthenticationSchemeBearer AuthenticationScheme = 1
	AuthenticationSchemeKey    AuthenticationScheme = 2
)

var authenticationSchemePrefixes = map[AuthenticationScheme]string{
	AuthenticationSchemeBearer: "Bearer ",
	AuthenticationSchemeKey:    "Key ",
}

func GetAuthenticationSchemeAndValue(header string) (AuthenticationScheme, string) {
	for scheme, prefix := range authenticationSchemePrefixes {
		if strings.HasPrefix(header, prefix) {
			return scheme, strings.TrimPrefix(header, prefix)
		}
	}
	return AuthenticationSchemeNone, ""
}
