package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/runeutil"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
)

type TextInput struct {
	charLimit int

	title    string
	validate func(string) error
	value    string
	style    TextInputStyle

	selected bool

	startValidation bool

	rsan runeutil.Sanitizer
}

type TextInputStyle struct {
	Title   lipgloss.Style
	Input   lipgloss.Style
	Faint   lipgloss.Style
	Invalid lipgloss.Style
	Help    lipgloss.Style
}

func NewTextInput(title string, validate func(string) error) *TextInput {
	return &TextInput{
		charLimit:       200,
		title:           title,
		validate:        validate,
		value:           "",
		startValidation: false,
		selected:        false,
		style: TextInputStyle{
			Title:   lipgloss.NewStyle().Foreground(colors.Teal).Bold(true),
			Input:   lipgloss.NewStyle().Foreground(colors.White),
			Faint:   lipgloss.NewStyle().Faint(true),
			Invalid: lipgloss.NewStyle().Foreground(colors.Red),
			Help:    lipgloss.NewStyle().Faint(true).Italic(true),
		},
	}
}

func (t *TextInput) Init() tea.Cmd {
	return nil
}

func (m *TextInput) san() runeutil.Sanitizer {
	if m.rsan == nil {
		// Textinput has all its input on a single line so collapse
		// newlines/tabs to single spaces.
		m.rsan = runeutil.NewSanitizer(
			runeutil.ReplaceTabs(" "), runeutil.ReplaceNewlines(" "))
	}
	return m.rsan
}

func (m *TextInput) insertRunesFromUserInput(v []rune) {
	// Clean up any special characters in the input provided by the
	// clipboard. This avoids bugs due to e.g. tab characters and
	// whatnot.
	paste := m.san().Sanitize(v)

	var availSpace int
	if m.charLimit > 0 {
		availSpace = m.charLimit - len(m.value)

		// If the char limit's been reached, cancel.
		if availSpace <= 0 {
			return
		}

		// If there's not enough space to paste the whole thing cut the pasted
		// runes down so they'll fit.
		if availSpace < len(paste) {
			paste = paste[:availSpace]
		}
	}

	// Put it all back together
	m.value = m.value + string(paste)
}

func (t *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "ctrl+\\":
			return t, tea.Quit
		case "enter":
			if err := t.validate(t.value); err != nil {
				t.startValidation = true
				return t, nil
			}
			return t, tea.Quit
		case "backspace":
			if len(t.value) > 0 {
				t.value = t.value[:len(t.value)-1]
			}
		case "delete":
			t.value = ""
		default:
			// Input one or more regular characters.
			t.insertRunesFromUserInput(msg.Runes)
		}
	}

	return t, nil
}

func (t *TextInput) View() string {
	var b strings.Builder

	titleStyle := t.style.Title
	if t.selected {
		titleStyle = t.style.Faint
	}

	showHelp := t.value != "" && t.startValidation
	validError := t.validate(t.value)

	inputStyle := t.style.Input
	if t.startValidation && validError != nil {
		inputStyle = t.style.Invalid
	}

	b.WriteString("\n")
	b.WriteString(titleStyle.Render(t.title))
	b.WriteString(" ")
	// TODO: Add cursor
	b.WriteString(inputStyle.Render(t.value))
	b.WriteString("\n")

	if validError != nil && showHelp {
		b.WriteString(t.style.Help.Render(validError.Error()))
		b.WriteString("\n")
	}

	if !t.selected {
		b.WriteString("\n")
	}

	return b.String()
}

func (t *TextInput) GetValue() string {
	return t.value
}

func RunTextInput(title string, validate func(string) error) (string, error) {
	t := NewTextInput(title, validate)
	p := tea.NewProgram(t)

	_, err := p.Run()
	if err != nil {
		return "", err
	}

	return t.GetValue(), nil
}
