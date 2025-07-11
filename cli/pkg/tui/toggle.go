package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
)

// ToggleSelect represents a toggle select input component, e.g. yes/no or yes/no/maybe
type ToggleSelect struct {
	items    []string
	cursor   int
	selected int
	title    string
	style    ToggleSelectStyle
}

// ToggleSelectStyle defines the styling for the select component
type ToggleSelectStyle struct {
	Title    lipgloss.Style
	Item     lipgloss.Style
	Selected lipgloss.Style
	Faint    lipgloss.Style
}

// NewToggleSelect creates a new toggle select component
func NewToggleSelect(items []string, title string) *ToggleSelect {
	return &ToggleSelect{
		items:    items,
		title:    title,
		selected: -1,
		style: ToggleSelectStyle{
			Title: lipgloss.NewStyle().
				Foreground(colors.Teal).
				Bold(true),
			Item: lipgloss.NewStyle().
				Foreground(colors.White),
			Selected: lipgloss.NewStyle().
				Foreground(colors.Teal).
				Bold(true),
			Faint: lipgloss.NewStyle().
				Faint(true),
		},
	}
}

// Init initializes the select component
func (s *ToggleSelect) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the component state
func (s *ToggleSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "left", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "right", "j":
			if s.cursor < len(s.items)-1 {
				s.cursor++
			}
		case "enter", " ":
			s.selected = s.cursor
			return s, tea.Quit
		case "q", "ctrl+c", "ctrl+\\", "esc":
			return s, tea.Quit
		}
	}

	return s, nil
}

// View renders the select component
func (s *ToggleSelect) View() string {
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
		b.WriteString(" ")
	}

	// Items
	shownItems := make([]string, 0, len(s.items))
	for i, item := range s.items {
		itemStyle := s.style.Item
		if s.cursor == i {
			itemStyle = s.style.Selected
		}

		shownItems = append(shownItems, itemStyle.Render(item))
	}

	b.WriteString(strings.Join(shownItems, "/"))

	return b.String()
}

// GetSelected returns the selected item and its index
func (s *ToggleSelect) GetSelected() (string, int) {
	if s.selected >= 0 && s.selected < len(s.items) {
		return s.items[s.selected], s.selected
	}
	return "", -1
}

// SetItems updates the items in the select component
func (s *ToggleSelect) SetItems(items []string) {
	s.items = items
	if s.cursor >= len(items) {
		s.cursor = len(items) - 1
	}
	if s.cursor < 0 && len(items) > 0 {
		s.cursor = 0
	}
}

// SetTitle updates the title of the select component
func (s *ToggleSelect) SetTitle(title string) {
	s.title = title
}

// RunToggleSelect runs the select component and returns the selected item and index
func RunToggleSelect(items []string, title string) (string, int, error) {
	if len(items) == 0 {
		return "", -1, fmt.Errorf("no items provided")
	}

	s := NewToggleSelect(items, title)
	p := tea.NewProgram(s)

	_, err := p.Run()
	if err != nil {
		return "", -1, err
	}

	selected, index := s.GetSelected()
	return selected, index, nil
}
