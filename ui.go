package main

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

type model struct {
	spinner      spinner.Model
	list         list.Model
	table        table.Model
	choice       string
	serverStatus string
	quitting     bool
	width        int
	height       int
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.serverStatus = string(msg)
		m.table.SetRows([]table.Row{
			{"1", string(msg), "Japan", "37,274,000"},
			{"2", string(msg), "Japan", "37,274,000"},
			{"2", string(msg), "Japan", "37,274,000"}})
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
				return m, tea.Batch(m.checkServer)
			}
			return m, nil
		case "esc":
			m.choice = ""
			m.serverStatus = ""
			return m, nil
		default:
			var spinnerCmd, tableCmd, listCmd tea.Cmd
			m.spinner, spinnerCmd = m.spinner.Update(msg)
			m.table, tableCmd = m.table.Update(msg)
			m.list, listCmd = m.list.Update(msg)
			return m, tea.Batch(spinnerCmd, tableCmd, listCmd)
		}
	default:
		var spinnerCmd, tableCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		m.table, tableCmd = m.table.Update(msg)
		return m, tea.Batch(spinnerCmd, tableCmd)

	}
}

func (m model) View() string {

	if m.serverStatus != "" {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, baseStyle.Render(m.table.View())))
	}
	if m.quitting {
		return quitTextStyle.Render("Leaving so early?")
	}

	if m.choice != "" {
		return lipgloss.Place(m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, fmt.Sprintf("%s Loading ...", m.spinner.View())))
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, m.list.View()))
}

func (m model) checkServer() tea.Msg {
	_, err := http.Get(m.choice)
	if err != nil {
		return statusMsg(fmt.Sprintf("%s, might be down", m.choice))
	}

	return statusMsg(fmt.Sprintf("%s, might is up", m.choice))
}

type statusMsg string
