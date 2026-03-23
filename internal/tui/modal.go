package tui

import (
	"fmt"
	"strings"

	"github.com/smangelsdorf/eric/internal/db"
)

func renderModal(task *db.Task, content string, scroll, width, height int) string {
	// lipgloss v1.1.0 semantics:
	//   .Width(w)  → content+padding = w, rendered = w + 2 (borders)
	//   .Height(h) → content+padding = h (floor, no truncation), rendered = h + 2
	// modalBorder has Padding(1, 2) → 2 vertical, 4 horizontal padding chars
	//
	// Caller passes height = m.height - 2. We need total rendered ≤ m.height - 1
	// (1 line for hints). So: heightParam + 2 = m.height - 1 → heightParam = height - 1.
	// Content lines = heightParam - 2 (padding) = height - 3.
	// Text width = width - 2 (borders) - 4 (padding) = width - 6.

	textWidth := width - 6
	contentLines := height - 3
	if textWidth < 20 {
		textWidth = 20
	}
	if contentLines < 5 {
		contentLines = 5
	}

	// Build header with task metadata
	var header []string
	if task != nil {
		header = append(header, fmt.Sprintf("%s  %s", task.ID, task.Summary))
		if task.Origin != "" || task.Destination != "" {
			origin := task.Origin
			if origin == "" {
				origin = "?"
			}
			dest := task.Destination
			if dest == "" {
				dest = "?"
			}
			header = append(header, fmt.Sprintf("%s → %s", origin, dest))
		}
		header = append(header, strings.Repeat("-", textWidth))
	}

	lines := strings.Split(content, "\n")

	// Wrap long lines (rune-aware to handle multi-byte chars like →)
	var wrapped []string
	for _, line := range append(header, lines...) {
		runes := []rune(line)
		if len(runes) <= textWidth {
			wrapped = append(wrapped, line)
		} else {
			for len(runes) > textWidth {
				wrapped = append(wrapped, string(runes[:textWidth]))
				runes = runes[textWidth:]
			}
			wrapped = append(wrapped, string(runes))
		}
	}

	// Apply scroll
	if scroll > len(wrapped)-contentLines {
		scroll = len(wrapped) - contentLines
	}
	if scroll < 0 {
		scroll = 0
	}

	end := scroll + contentLines
	if end > len(wrapped) {
		end = len(wrapped)
	}
	visible := wrapped[scroll:end]

	body := strings.Join(visible, "\n")

	heightParam := height - 1
	widthParam := width - 2
	return modalBorder.Width(widthParam).Height(heightParam).Render(body)
}
