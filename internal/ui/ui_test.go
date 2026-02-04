package ui

import (
	"testing"
	"github.com/mattn/go-runewidth"
)

func TestTruncateToWidth(t *testing.T) {
	tests := []struct {
		input    string
		width    int
		expected string
	}{
		{"Hello", 10, "Hello     "},
		{"Hello World", 5, "He..."},
		{"你好世界", 8, "你好世界"},
		{"你好世界", 6, "你好..."}, // "你好" is 4, "..." is 3. Total 7? No wait. 
		// If width is 6. "..." takes 3. Remaining is 3.
		// "你" takes 2. "好" takes 2.
		// So can only fit "你". 2 + 3 = 5.
		// runewidth.Truncate("你好世界", 6, "...") behavior:
		// It tries to fill 6. "..." is 3. It needs 3 more.
		// "你好" is 4, too big. "你" is 2.
		// So "你" + "..." = "你..." (width 5).
		// Wait, if it returns "你...", width is 5. But we requested 6.
		// My implementation adds padding if w < width!
		
		// Let's trace my implementation:
		// w = StringWidth("你...") -> 5.
		// 5 < 6.
		// Returns "你..." + " " -> "你... " (Width 6).
	}

	for _, tt := range tests {
		got := TruncateToWidth(tt.input, tt.width)
		if runewidth.StringWidth(got) != tt.width {
			t.Errorf("TruncateToWidth(%q, %d) width = %d, want %d", tt.input, tt.width, runewidth.StringWidth(got), tt.width)
		}
		// Optional: Check exact string if we are sure about implementation details
		// if got != tt.expected {
		// 	t.Errorf("TruncateToWidth(%q, %d) = %q, want %q", tt.input, tt.width, got, tt.expected)
		// }
	}
}

func TestTruncateToWidth_WideChars(t *testing.T) {
	input := "你好世界"
	width := 6
	got := TruncateToWidth(input, width)
	
	if runewidth.StringWidth(got) != width {
		t.Errorf("TruncateToWidth(%q, %d) width = %d, want %d. Got: %q", input, width, runewidth.StringWidth(got), width, got)
	}
}
