package config

import (
	"reflect"
	"testing"
)

func TestSplitList(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"simple", "a,b,c", []string{"a", "b", "c"}},
		{"trims spaces", " a , b ,c ", []string{"a", "b", "c"}},
		{"skips blanks", "a,,b,", []string{"a", "b"}},
		{"empty", "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitList(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitList(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}