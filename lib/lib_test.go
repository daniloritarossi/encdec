package lib

import "testing"

func TestValidateMachineID(t *testing.T) {
	cases := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"empty", "", true},
		{"whitespace", "  \n\t\n", true},
		{"placeholder", "uninitialized", true},
		{"placeholder with newline", "uninitialized\n", true},
		{"placeholder upper case", "UNINITIALIZED", true},
		{"real id", "0123456789abcdef0123456789abcdef", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := validateMachineID(c.id); (err != nil) != c.wantErr {
				t.Fatalf("validateMachineID(%q) error = %v, wantErr %v", c.id, err, c.wantErr)
			}
		})
	}
}
