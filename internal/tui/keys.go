package tui

import "strings"

type hint struct {
	key   string
	label string
}

func renderHints(hints []hint) string {
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = hintKey.Render(h.key) + " " + hintLabel.Render(h.label)
	}
	return strings.Join(parts, "  ")
}

var tableHints = []hint{
	{"j/k", "navigate"},
	{"space", "view"},
	{"c", "close"},
	{"o", "reopen"},
	{"q", "quit"},
}

var modalHints = []hint{
	{"j/k", "scroll"},
	{"space/esc", "back"},
}
