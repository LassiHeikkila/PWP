package auth

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoundTripUserToken(t *testing.T) {
	a := NewController([]byte("my-test-key"))

	user := "user"
	organization := "organization"
	role := 1
	claims := CreateUserClaims(user, organization, role)

	token, err := a.CreateJWT(claims)
	if err != nil {
		t.Fatal("error creating token:", err)
	}

	var u string
	var o string
	var r int
	if !a.ValidateUserToken(token, &u, &o, &r) {
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

func TestRoundTripMachineToken(t *testing.T) {
	a := NewController([]byte("my-test-key"))

	machine := "machine"
	organization := "organization"

	claims := CreateMachineClaims(machine, organization)

	token, err := a.CreateJWT(claims)
	if err != nil {
		t.Fatal("error creating token:", err)
	}

	var m string
	var o string
	if !a.ValidateMachineToken(token, &m, &o) {
		t.Fatal("validation failed")
	}
	if !cmp.Equal(machine, m) {
		t.Fatal(cmp.Diff(machine, m))
	}
	if !cmp.Equal(organization, o) {
		t.Fatal(cmp.Diff(organization, o))
	}
}

func TestBadTokenParsing(t *testing.T) {
	a := NewController([]byte("bad-token-parsing-key"))
	b := NewController([]byte("other-key"))

	mustCreateTokenB := func(user string, org string, role int) string {
		token, err := b.CreateJWT(CreateUserClaims(user, org, role))
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
			if a.ValidateUserToken(tc.token, &tc.user, &tc.organization, &tc.role) {
				t.Fatal("token validation passed for bad inputs")
			}
		})
	}
}

func TestControllerConstructor(t *testing.T) {
	t.Run("with key", func(t *testing.T) {
		key := []byte("my-key")
		a := NewController(key)
		if a == nil {
			t.Fatal("nil controller returned")
		}
		if !cmp.Equal(key, a.(*authController).key) {
			t.Fatal("wrong key")
		}
	})
	t.Run("without key", func(t *testing.T) {
		a := NewController(nil)
		if a == nil {
			t.Fatal("nil controller returned")
		}
		if len(a.(*authController).key) == 0 {
			t.Fatal("no key generated")
		}
	})
}
