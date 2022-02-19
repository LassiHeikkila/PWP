package auth

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoundRobin(t *testing.T) {
	a := NewController([]byte("my-test-key"))

	user := "user"
	organization := "organization"
	role := 1
	claims := CreateClaims(user, organization, role)

	token, err := a.CreateJWT(claims)
	if err != nil {
		t.Fatal("error creating token:", err)
	}

	var u string
	var o string
	var r int
	if !a.ValidateToken(token, &u, &o, &r) {
		t.Fatal("validation failed")
	}
	if !cmp.Equal(user, u) {
		t.Fatal(cmp.Diff(user, u))
	}
	if !cmp.Equal(organization, o) {
		t.Fatal(cmp.Diff(organization, o))
	}
	if !cmp.Equal(role, r) {
		t.Fatal(cmp.Diff(role, r))
	}
}

func TestBadTokenParsing(t *testing.T) {
	a := NewController([]byte("bad-token-parsing-key"))
	b := NewController([]byte("other-key"))

	mustCreateTokenB := func(user string, org string, role int) string {
		token, err := b.CreateJWT(CreateClaims(user, org, role))
		if err != nil {
			panic("error creating token: " + err.Error())
		}
		return token
	}

	tests := []struct {
		token        string
		user         string
		organization string
		role         int
	}{
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
			token:        mustCreateTokenB("userB", "orgB", 12), // signed with different key
			user:         "userB",
			organization: "orgB",
			role:         12,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if a.ValidateToken(tc.token, &tc.user, &tc.organization, &tc.role) {
				t.Fatal("token validation passed for bad inputs")
			}
		})
	}
}
