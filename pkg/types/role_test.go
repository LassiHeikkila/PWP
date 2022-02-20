package types

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
			b:    RoleMaintainer,
			want: true,
		},
		"2": {
			in:   RoleUser | RoleMaintainer | RoleAdministrator,
			b:    RoleAdministrator,
			want: true,
		},
		"3": {
			in:   RoleUser | RoleMaintainer | RoleAdministrator,
			b:    RoleRoot,
			want: false,
		},
		"4": {
			in:   RoleUser | RoleMaintainer,
			b:    RoleUser,
			want: true,
		},
		"5": {
			in:   RoleUser | RoleMaintainer,
			b:    RoleMaintainer,
			want: true,
		},
		"6": {
			in:   RoleUser | RoleMaintainer,
			b:    RoleAdministrator,
			want: false,
		},
		"7": {
			in:   RoleUser | RoleMaintainer,
			b:    RoleRoot,
			want: false,
		},
		"8": {
			in:   RoleRoot,
			b:    RoleUser,
			want: false,
		},
		"9": {
			in:   RoleRoot,
			b:    RoleMaintainer,
			want: false,
		},
		"10": {
			in:   RoleRoot,
			b:    RoleAdministrator,
			want: false,
		},
		"11": {
			in:   RoleRoot,
			b:    RoleRoot,
			want: true,
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
