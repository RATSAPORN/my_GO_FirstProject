package configs

import "testing"

func TestAtoiDefault(t *testing.T) {
	tests := []struct {
		name string
		in   string
		def  int
		want int
	}{
		{name: "empty uses default", in: "", def: 20, want: 20},
		{name: "invalid uses default", in: "abc", def: 20, want: 20},
		{name: "valid int", in: "5", def: 20, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := atoiDefault(tt.in, tt.def); got != tt.want {
				t.Fatalf("atoiDefault(%q, %d) = %d, want %d", tt.in, tt.def, got, tt.want)
			}
		})
	}
}
