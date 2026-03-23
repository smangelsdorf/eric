package tui

import (
	"database/sql"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smangelsdorf/eric/internal/db"
	"github.com/smangelsdorf/eric/internal/storage"
)

type viewState int

const (
	viewTable viewState = iota
	viewModal
)

type Model struct {
	db           *sql.DB
	tasks        []db.Task
	cursor       int
	view         viewState
	modalTask    *db.Task
	modalContent string
	modalScroll  int
	width        int
	height       int
	dataVersion  int
	err          error
}

func NewModel(database *sql.DB) Model {
	return Model{
		db:     database,
		width:  80,
		height: 24,
	}
}

// Messages

type tickMsg time.Time
type tasksLoadedMsg struct {
	tasks []db.Task
	err   error
}
type dataVersionMsg struct {
	version int
}
type modalContentMsg struct {
	task    *db.Task
	content string
	err     error
}

// Commands

func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func loadTasksCmd(database *sql.DB) tea.Cmd {
	return func() tea.Msg {
		tasks, err := db.ListTasks(database, db.TaskFilter{})
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

func checkDataVersionCmd(database *sql.DB) tea.Cmd {
	return func() tea.Msg {
		var version int
		database.QueryRow("PRAGMA data_version").Scan(&version)
		return dataVersionMsg{version: version}
	}
}

func loadModalContentCmd(task db.Task) tea.Cmd {
	return func() tea.Msg {
		content, err := storage.ReadTaskFile(task.FilePath)
		return modalContentMsg{task: &task, content: content, err: err}
	}
}

// Bubble Tea interface

func (m Model) Init() tea.Cmd {
	return tea.Batch(loadTasksCmd(m.db), checkDataVersionCmd(m.db), tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tasksLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		selectedID := ""
		if m.cursor < len(m.tasks) {
			selectedID = m.tasks[m.cursor].ID
		}
		m.tasks = sortTasks(msg.tasks)
		m.err = nil
		// Restore cursor position
		if selectedID != "" {
			for i, t := range m.tasks {
				if t.ID == selectedID {
					m.cursor = i
					break
				}
			}
		}
		m.clampCursor()
		return m, nil

	case dataVersionMsg:
		if msg.version != m.dataVersion {
			m.dataVersion = msg.version
			return m, loadTasksCmd(m.db)
		}
		return m, nil

	case tickMsg:
		return m, tea.Batch(checkDataVersionCmd(m.db), tickCmd())

	case modalContentMsg:
		if msg.err != nil {
			m.err = msg.err
			m.view = viewTable
			return m, nil
		}
		m.modalContent = msg.content
		m.modalScroll = 0
		m.view = viewModal
		m.modalTask = msg.task
		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.view == viewModal {
		return m.handleModalKey(msg)
	}
	return m.handleTableKey(msg)
}

func (m Model) handleTableKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		m.cursor++
		m.clampCursor()
	case "k", "up":
		m.cursor--
		m.clampCursor()
	case " ":
		if m.cursor < len(m.tasks) {
			return m, loadModalContentCmd(m.tasks[m.cursor])
		}
	case "s":
		if m.cursor < len(m.tasks) && m.tasks[m.cursor].Status == "open" {
			task := m.tasks[m.cursor]
			db.StartTask(m.db, task.ID)
			return m, loadTasksCmd(m.db)
		}
	case "c":
		if m.cursor < len(m.tasks) && (m.tasks[m.cursor].Status == "open" || m.tasks[m.cursor].Status == "in_progress") {
			task := m.tasks[m.cursor]
			db.UpdateTaskStatus(m.db, task.ID, "closed")
			return m, loadTasksCmd(m.db)
		}
	case "o":
		if m.cursor < len(m.tasks) && m.tasks[m.cursor].Status != "open" {
			task := m.tasks[m.cursor]
			db.UpdateTaskStatus(m.db, task.ID, "open")
			return m, loadTasksCmd(m.db)
		}
	}
	return m, nil
}

func (m Model) handleModalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case " ", "esc", "q":
		m.view = viewTable
		m.modalContent = ""
		m.modalTask = nil
		m.modalScroll = 0
	case "j", "down":
		m.modalScroll++
	case "k", "up":
		m.modalScroll--
		if m.modalScroll < 0 {
			m.modalScroll = 0
		}
	}
	return m, nil
}

func (m *Model) clampCursor() {
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.tasks) {
		m.cursor = len(m.tasks) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m Model) View() string {
	if m.view == viewModal {
		modal := renderModal(m.modalTask, m.modalContent, m.modalScroll, m.width, m.height-2) // reserve hint bar
		hints := renderHints(modalHints)
		return lipgloss.JoinVertical(lipgloss.Left, modal, hints)
	}

	pageSize := m.height - 6 // borders + header + separator + hint bar + padding
	if pageSize < 1 {
		pageSize = 1
	}

	extra := 0
	var errLine string
	if m.err != nil {
		errLine = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("Error: " + m.err.Error())
		extra = 1
	}

	table := renderTable(m.tasks, m.cursor, m.width-2, pageSize)
	// Border adds 2 lines (top+bottom) + hints 1 line
	tableHeight := m.height - 3 - extra
	bordered := tableBorder.Width(m.width - 2).Height(tableHeight).Render(table)
	bordered = replaceTopBorder(bordered, "eric-tui", m.width, tableBorderColor)
	hints := renderHints(tableHints)

	parts := []string{}
	if errLine != "" {
		parts = append(parts, errLine)
	}
	parts = append(parts, bordered, hints)
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func statusOrder(s string) int {
	switch s {
	case "in_progress":
		return 0
	case "open":
		return 1
	default:
		return 2
	}
}

func sortTasks(tasks []db.Task) []db.Task {
	sort.SliceStable(tasks, func(i, j int) bool {
		oi, oj := statusOrder(tasks[i].Status), statusOrder(tasks[j].Status)
		if oi != oj {
			return oi < oj
		}
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt) // newest first within group
	})
	return tasks
}
