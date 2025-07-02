package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
)

// Select represents a select input component
type Select struct {
	items        []string
	cursor       int
	selected     int
	title        string
	width        int
	height       int
	showHelp     bool
	helpText     string
	style        SelectStyle
	maxDisplayed int
	viewCursor   int
}

func (s *Select) numToDisplay() int {
	TitleAndHelpHeight := 3
	return max(min(min(s.maxDisplayed, len(s.items)), s.height-TitleAndHelpHeight), 1)
}

// SelectStyle defines the styling for the select component
type SelectStyle struct {
	Title    lipgloss.Style
	Item     lipgloss.Style
	Selected lipgloss.Style
	Faint    lipgloss.Style
	Cursor   lipgloss.Style
	Help     lipgloss.Style
}

// NewSelect creates a new select component
func NewSelect(items []string, title string, maxDisplayed int) *Select {
	return &Select{
		items:        items,
		title:        title,
		selected:     -1,
		showHelp:     true,
		helpText:     "↑/↓: navigate • enter: select",
		maxDisplayed: maxDisplayed,
		viewCursor:   0,
		style: SelectStyle{
			Title: lipgloss.NewStyle().
				Foreground(colors.Teal).
				Bold(true),
			Item: lipgloss.NewStyle().
				Foreground(colors.White).
				MarginLeft(1),
			Selected: lipgloss.NewStyle().
				Foreground(colors.Teal).
				Bold(true).
				MarginLeft(1),
			Cursor: lipgloss.NewStyle().
				Foreground(colors.Teal).
				Bold(true).
				MarginLeft(1),
			Faint: lipgloss.NewStyle().
				Faint(true),
			Help: lipgloss.NewStyle().
				Foreground(colors.Gray).
				Italic(true),
		},
	}
}

// Init initializes the select component
func (s *Select) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the component state
func (s *Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
				if s.cursor < s.viewCursor {
					s.viewCursor--
				}
			}
		case "down", "j":
			if s.cursor < len(s.items)-1 {
				s.cursor++
				if s.cursor >= s.viewCursor+s.numToDisplay() {
					s.viewCursor++
				}
			}
		case "enter", " ":
			s.selected = s.cursor
			return s, tea.Quit
		case "q", "ctrl+c", "ctrl+\\", "esc":
			return s, tea.Quit
		}
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		// Center the view, but don't go below 0
		s.viewCursor = max(0, s.cursor-(s.numToDisplay()/2))
		return s, nil
	}

	return s, nil
}

func (s *Select) isInView(itemIndex int) bool {
	return itemIndex >= s.viewCursor && itemIndex < s.viewCursor+s.numToDisplay()
}

// View renders the select component
func (s *Select) View() string {
	if len(s.items) == 0 {
		return "No items available"
	}

	var b strings.Builder

	if s.selected != -1 {
		return s.style.Faint.Render(s.title) + " " + s.items[s.selected] + "\n"
	}

	// Title
	if s.title != "" {
		b.WriteString(s.style.Title.Render(s.title))
		b.WriteString("\n")
	}

	// Help text
	if s.showHelp && s.helpText != "" && s.height >= 3 {
		b.WriteString(s.style.Help.Render(s.helpText))
		b.WriteString("\n")
		if s.height >= 4 {
			b.WriteString("\n")
		}
	}

	// Items
	shownItems := make([]string, 0, s.numToDisplay())
	for i, item := range s.items {
		if !s.isInView(i) {
			continue
		}

		cursor := " "
		if s.cursor == i {
			cursor = s.style.Cursor.Render(icons.RightArrow)
		}

		itemStyle := s.style.Item
		if s.cursor == i {
			itemStyle = s.style.Selected
		}

		shownItems = append(shownItems, fmt.Sprintf("%s%s", cursor, itemStyle.Render(item)))
	}

	b.WriteString(strings.Join(shownItems, "\n"))

	return b.String()
}

// GetSelected returns the selected item and its index
func (s *Select) GetSelected() (string, int) {
	if s.selected >= 0 && s.selected < len(s.items) {
		return s.items[s.selected], s.selected
	}
	return "", -1
}

// SetItems updates the items in the select component
func (s *Select) SetItems(items []string) {
	s.items = items
	if s.cursor >= len(items) {
		s.cursor = len(items) - 1
	}
	if s.cursor < 0 && len(items) > 0 {
		s.cursor = 0
	}
}

// SetTitle updates the title of the select component
func (s *Select) SetTitle(title string) {
	s.title = title
}

// SetHelpText updates the help text
func (s *Select) SetHelpText(helpText string) {
	s.helpText = helpText
}

// SetShowHelp toggles the visibility of help text
func (s *Select) SetShowHelp(show bool) {
	s.showHelp = show
}

// RunSelect runs the select component and returns the selected item and index
func RunSelect(items []string, title string) (string, int, error) {
	if len(items) == 0 {
		return "", -1, fmt.Errorf("no items provided")
	}

	s := NewSelect(items, title, 5)
	p := tea.NewProgram(s)

	_, err := p.Run()
	if err != nil {
		return "", -1, err
	}

	selected, index := s.GetSelected()
	return selected, index, nil
}
