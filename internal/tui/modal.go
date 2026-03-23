package tui

import (
	"fmt"
	"strings"

	"github.com/smangelsdorf/eric/internal/db"
)

func renderModal(task *db.Task, content string, scroll, width, height int) string {
	// Reserve space for border (2 top/bottom) and padding (1 top/bottom)
	innerWidth := width - 6
	innerHeight := height - 4
	if innerWidth < 20 {
		innerWidth = 20
	}
	if innerHeight < 5 {
		innerHeight = 5
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
		header = append(header, strings.Repeat("─", innerWidth))
	}

	lines := strings.Split(content, "\n")

	// Wrap long lines
	var wrapped []string
	for _, line := range append(header, lines...) {
		if len(line) <= innerWidth {
			wrapped = append(wrapped, line)
		} else {
			for len(line) > innerWidth {
				wrapped = append(wrapped, line[:innerWidth])
				line = line[innerWidth:]
			}
			wrapped = append(wrapped, line)
		}
	}

	// Apply scroll
	if scroll > len(wrapped)-innerHeight {
		scroll = len(wrapped) - innerHeight
	}
	if scroll < 0 {
		scroll = 0
	}

	end := scroll + innerHeight
	if end > len(wrapped) {
		end = len(wrapped)
	}
	visible := wrapped[scroll:end]

	body := strings.Join(visible, "\n")

	return modalBorder.Width(innerWidth).Height(innerHeight).Render(body)
}
