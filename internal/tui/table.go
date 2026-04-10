package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/smangelsdorf/eric/internal/db"
)

const (
	colIDWidth     = 10
	colStatusWidth = 13
	colDateWidth   = 12
	colOriginWidth = 15
	colDestWidth   = 15
	colPadding     = 2 // padding per cell (0,1 on each side)

	// Minimum terminal widths to show optional columns
	originMinWidth = 120
	destMinWidth   = 150
)

func renderTable(tasks []db.Task, cursor, width, pageSize int) string {
	showOrigin := width >= originMinWidth
	showDest := width >= destMinWidth

	colCount := 4
	fixedWidth := colIDWidth + colStatusWidth + colDateWidth
	if showOrigin {
		colCount++
		fixedWidth += colOriginWidth
	}
	if showDest {
		colCount++
		fixedWidth += colDestWidth
	}

	summaryWidth := width - fixedWidth - colCount*colPadding - 2 // 2 for borders
	if summaryWidth < 10 {
		summaryWidth = 10
	}

	var b strings.Builder

	// Header
	h := headerStyle.Render(pad("ID", colIDWidth)) +
		headerStyle.Render(pad("Summary", summaryWidth)) +
		headerStyle.Render(pad("Status", colStatusWidth))
	if showOrigin {
		h += headerStyle.Render(pad("Origin", colOriginWidth))
	}
	if showDest {
		h += headerStyle.Render(pad("Dest", colDestWidth))
	}
	h += headerStyle.Render(pad("Created", colDateWidth))
	b.WriteString(h)
	b.WriteString("\n")

	// Separator
	b.WriteString(strings.Repeat("─", width-2))
	b.WriteString("\n")

	if len(tasks) == 0 {
		b.WriteString(cellStyle.Render("No tasks found."))
		b.WriteString("\n")
		return b.String()
	}

	page := cursor / pageSize
	start := page * pageSize
	end := start + pageSize
	if end > len(tasks) {
		end = len(tasks)
	}

	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, -7)

	for i := start; i < end; i++ {
		t := tasks[i]
		isClosed := t.Status == "closed"

		id := pad(t.ID, colIDWidth)
		summary := truncate(t.Summary, summaryWidth)
		summary = pad(summary, summaryWidth)
		status := pad(t.Status, colStatusWidth)
		date := pad(t.CreatedAt.Format("2006-01-02"), colDateWidth)

		var origin, dest string
		if showOrigin {
			origin = truncate(t.Origin, colOriginWidth)
			origin = pad(origin, colOriginWidth)
		}
		if showDest {
			dest = truncate(t.Destination, colDestWidth)
			dest = pad(dest, colDestWidth)
		}

		// Per-cell styling
		if isClosed {
			id = closedRowStyle.Render(id)
			summary = closedRowStyle.Render(summary)
			status = statusClosed.Render(status)
			date = closedRowStyle.Render(date)
			if showOrigin {
				origin = closedRowStyle.Render(origin)
			}
			if showDest {
				dest = closedRowStyle.Render(dest)
			}
		} else if t.Status == "in_progress" {
			id = inProgressRowStyle.Render(id)
			summary = inProgressRowStyle.Render(summary)
			status = statusInProgress.Render(status)
			date = inProgressRowStyle.Render(date)
			if showOrigin {
				origin = inProgressRowStyle.Render(origin)
			}
			if showDest {
				dest = inProgressRowStyle.Render(dest)
			}
		} else {
			if t.Status == "open" {
				status = statusOpen.Render(status)
			}
			if t.CreatedAt.Before(oneWeekAgo) {
				date = dateOld.Render(date)
			}
		}

		row := cellStyle.Render(id) +
			cellStyle.Render(summary) +
			cellStyle.Render(status)
		if showOrigin {
			row += cellStyle.Render(origin)
		}
		if showDest {
			row += cellStyle.Render(dest)
		}
		row += cellStyle.Render(date)

		if i == cursor {
			row = cursorStyle.Render(row)
		}

		b.WriteString(row)
		b.WriteString("\n")
	}

	// Page indicator
	totalPages := (len(tasks) + pageSize - 1) / pageSize
	if totalPages > 1 {
		indicator := fmt.Sprintf(" Page %d/%d ", page+1, totalPages)
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(indicator))
		b.WriteString("\n")
	}

	return b.String()
}

// replaceTopBorder swaps the first line of a lipgloss-rendered bordered box
// with a title line like: ╭─ eric-tui ────...──╮ styled in the border color.
func replaceTopBorder(rendered, title string, width int, borderColor lipgloss.Color) string {
	label := "─ " + title + " "
	fillLen := width - 2 - len([]rune(label))
	if fillLen < 0 {
		fillLen = 0
	}
	line := "╭" + label + strings.Repeat("─", fillLen) + "╮"
	styled := lipgloss.NewStyle().Foreground(borderColor).Render(line)

	if i := strings.IndexByte(rendered, '\n'); i >= 0 {
		return styled + rendered[i:]
	}
	return styled
}

func pad(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
