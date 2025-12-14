package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Bold(true).
			Foreground(lipgloss.Color("212"))
	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("212"))
)

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return ""
	}
	if m.quitting {
		return "Operation cancelled.\n"
	}
	return "\n" + m.list.View()
}

func Select(title string, options []string) string {
	items := []list.Item{}
	for _, opt := range options {
		items = append(items, item{title: opt, desc: ""})
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)
	delegate.Styles.NormalTitle = itemStyle
	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.SelectedDesc = selectedItemStyle

	const defaultWidth = 20
	const chromeHeight = 6
	listHeight := min(len(items)+chromeHeight, 20)

	l := list.New(items, delegate, defaultWidth, listHeight)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle

	m := model{list: l}
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running UI: ", err)
		os.Exit(1)
	}

	if finalM, ok := finalModel.(model); ok {
		if finalM.choice == "" && finalM.quitting {
			fmt.Println("Selection cancelled.")
			os.Exit(0)
		}
		return finalM.choice
	}

	return ""
}
