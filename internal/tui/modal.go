package tui

import "strings"

func renderModal(content string, scroll, width, height int) string {
	// Reserve space for border (2 top/bottom) and padding (1 top/bottom)
	innerWidth := width - 6
	innerHeight := height - 6
	if innerWidth < 20 {
		innerWidth = 20
	}
	if innerHeight < 5 {
		innerHeight = 5
	}

	lines := strings.Split(content, "\n")

	// Wrap long lines
	var wrapped []string
	for _, line := range lines {
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

	return modalBorder.Width(innerWidth).Render(body)
}
