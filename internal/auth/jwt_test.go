package auth

import (
	"fmt"
	"testing"
)

func TestRoundRobin(t *testing.T) {
	a := NewController([]byte("my-test-key"))

	user := "user"
	organization := "organization"
	claims := CreateClaims(user, organization)

	token, err := a.CreateJWT(claims)
	if err != nil {
		t.Fatal("error creating token:", err)
	}

	if !a.ValidateToken(token, user, organization) {
		t.Fatal("validation failed")
	}
}

func TestBadTokenParsing(t *testing.T) {
	a := NewController([]byte("bad-token-parsing-key"))
	b := NewController([]byte("other-key"))

	mustCreateToken := func(user string, org string) string {
		token, err := a.CreateJWT(CreateClaims(user, org))
		if err != nil {
			panic("error creating token: " + err.Error())
		}
		return token
	}
	mustCreateTokenB := func(user string, org string) string {
		token, err := b.CreateJWT(CreateClaims(user, org))
		if err != nil {
			panic("error creating token: " + err.Error())
		}
		return token
	}

	tests := []struct {
		token        string
		user         string
		organization string
	}{
		{
			token:        mustCreateToken("user1", "org1"),
			user:         "user",
			organization: "org1",
		},
		{
			token:        mustCreateToken("user1", "org1"),
			user:         "user1",
			organization: "org",
		},
		{
			token:        "09j92932d23d3",
			user:         "user1",
			organization: "org1",
		},
		{
			token:        "",
			user:         "",
			organization: "",
		},
		{
			token:        "09j92932d23d3.2d923d2903jd0j23d.892hd832hd",
			user:         "user1",
			organization: "org1",
		},
		{
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU",
			user:         "",
			organization: "",
		},
		{
			token:        mustCreateTokenB("userB", "orgB"), // signed with different key
			user:         "userB",
			organization: "orgB",
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if a.ValidateToken(tc.token, tc.user, tc.organization) {
				t.Fatal("token validation passed for bad inputs")
			}
		})
	}
}
