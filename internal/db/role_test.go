package db

import (
	"testing"
)

func TestHasRole(t *testing.T) {
	tests := map[string]struct {
		in   Role
		b    Role
		want bool
	}{
		"0": {
			in:   RoleUser | RoleMaintainer | RoleAdministrator,
			b:    RoleUser,
			want: true,
		},
		"1": {
			in:   RoleUser | RoleMaintainer | RoleAdministrator,
			b:    RoleRoot,
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := HasRole(tc.in, tc.b)
			if tc.want != got {
				t.Fatal("wanted", tc.want, "got", got)
			}
		})
	}
}
