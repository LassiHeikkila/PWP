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

func GetAuthenticationScheme(header string) AuthenticationScheme {
	for scheme, prefix := range authenticationSchemePrefixes {
		if strings.HasPrefix(header, prefix) {
			return scheme
		}
	}
	return AuthenticationSchemeNone
}
