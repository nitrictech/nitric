package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
)

type textInputModel struct {
	textinput textinput.Model
	title     string
	style     textInputStyle
	showError bool

	value string
}

type textInputStyle struct {
	Title   lipgloss.Style
	Input   lipgloss.Style
	Faint   lipgloss.Style
	Invalid lipgloss.Style
	Help    lipgloss.Style
}

func (m textInputModel) Init() tea.Cmd {
	m.textinput.Focus()
	return textinput.Blink
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "ctrl+\\":
			return m, tea.Quit
		case "enter":
			if m.textinput.Err != nil {
				m.showError = true
				return m, nil
			}
			m.value = m.textinput.Value()
			m.textinput.PromptStyle = m.style.Faint.MarginRight(1)
			m.textinput.Blur()
			return m, tea.Quit
		}
	}

	// Always pass the message to the internal textinput component
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m textInputModel) View() string {
	var b strings.Builder

	m.textinput.TextStyle = m.style.Input
	if m.showError && m.textinput.Err != nil {
		m.textinput.TextStyle = m.style.Invalid
	}

	b.WriteString(m.textinput.View())
	b.WriteString("\n")

	if m.textinput.Err != nil && m.showError {
		b.WriteString(m.style.Help.Render(m.textinput.Err.Error()))
		b.WriteString("\n")
	}

	if m.textinput.Value() == "" {
		b.WriteString("\n")
	}

	return b.String()
}

func RunTextInput(title string, validate func(string) error) (string, error) {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.CharLimit = 200
	ti.Width = 50
	ti.Prompt = title
	ti.PromptStyle = lipgloss.NewStyle().Foreground(colors.Teal).Bold(true).MarginRight(1)
	ti.Validate = validate
	ti.Focus()

	style := textInputStyle{
		Title:   lipgloss.NewStyle().Foreground(colors.Teal).Bold(true),
		Input:   lipgloss.NewStyle().Foreground(colors.White),
		Faint:   lipgloss.NewStyle().Faint(true),
		Invalid: lipgloss.NewStyle().Foreground(colors.Red),
		Help:    lipgloss.NewStyle().Faint(true).Italic(true),
	}

	model := textInputModel{
		textinput: ti,
		title:     title,
		style:     style,
	}

	p := tea.NewProgram(model)
	m, err := p.Run()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return "", err
	}

	return m.(textInputModel).value, nil
}
