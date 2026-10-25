package matcher

import "testing"

func TestMatch(t *testing.T) {
	m := New([]string{"київська", "київ"})

	tests := []struct {
		name string
		text string
		want bool
	}{
		{"keyword present", "Загроза для Київська область", true},
		{"different case", "КИЇВСЬКА ОБЛАСТЬ", true},
		{"second keyword", "Шахеди курсом на Київ", true},
		{"no keyword", "Одеська область, відбій", false},
		{"empty text", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.Match(tt.text); got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}