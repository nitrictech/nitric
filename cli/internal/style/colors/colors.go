package colors

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
)

//

var (
	White  = lipgloss.ANSIColor(7)
	Gray   = lipgloss.ANSIColor(8)
	Black  = lipgloss.ANSIColor(0)
	Red    = lipgloss.ANSIColor(1)
	Orange = lipgloss.ANSIColor(3)
	Yellow = lipgloss.ANSIColor(11)
	Green  = lipgloss.ANSIColor(10)
	Teal   = lipgloss.ANSIColor(14)
	Blue   = lipgloss.ANSIColor(4)
	Purple = lipgloss.ANSIColor(13)
)

// var (
// 	AdaptiveWhite  = lipgloss.CompleteColor{TrueColor: "#FFFFFF", ANSI256: "255", ANSI: "15"}
// 	AdaptiveGray   = lipgloss.CompleteColor{TrueColor: "#696969", ANSI256: "250", ANSI: "7"}
// 	AdaptiveBlack  = lipgloss.CompleteColor{TrueColor: "#000000", ANSI256: "16", ANSI: "0"}
// 	AdaptiveRed    = lipgloss.CompleteColor{TrueColor: "#ff4499", ANSI256: "197", ANSI: "1"}
// 	AdaptiveOrange = lipgloss.CompleteColor{TrueColor: "#F97316", ANSI256: "208", ANSI: "3"}
// 	AdaptiveYellow = lipgloss.CompleteColor{TrueColor: "#FDE047", ANSI256: "220", ANSI: "11"}
// 	AdaptiveGreen  = lipgloss.CompleteColor{TrueColor: "#00ffd2", ANSI256: "47", ANSI: "10"}
// 	AdaptiveTeal   = lipgloss.CompleteColor{TrueColor: "#32D0D1", ANSI256: "51", ANSI: "14"}
// 	AdaptiveBlue   = lipgloss.CompleteColor{TrueColor: "#2563EB", ANSI256: "21", ANSI: "4"}
// 	AdaptivePurple = lipgloss.CompleteColor{TrueColor: "#C27AFA", ANSI256: "99", ANSI: "13"}
// )

// var (
// 	Text = lipgloss.CompleteAdaptiveColor{
// 		Light: AdaptiveBlack,
// 		Dark:  AdaptiveWhite,
// 	}
// 	TextMuted = lipgloss.CompleteAdaptiveColor{
// 		Light: lipgloss.CompleteColor{TrueColor: "#4E4E4E", ANSI256: "242", ANSI: "8"},
// 		Dark:  lipgloss.CompleteColor{TrueColor: "#5E5E5E", ANSI256: "249", ANSI: "7"},
// 	}
// 	TextHighlight = lipgloss.CompleteAdaptiveColor{
// 		Light: AdaptiveBlue,
// 		Dark:  AdaptivePurple,
// 	}
// 	TextActive = lipgloss.CompleteAdaptiveColor{
// 		Light: AdaptiveBlue,
// 		Dark:  AdaptiveBlue,
// 	}
// )

var Theme = huh.Theme{
	Form: huh.FormStyles{
		Base: lipgloss.NewStyle(),
	},
	Group: huh.GroupStyles{
		Base:        lipgloss.NewStyle(),
		Title:       lipgloss.NewStyle(),
		Description: lipgloss.NewStyle(),
	},
	FieldSeparator: lipgloss.NewStyle(),
	Blurred: huh.FieldStyles{
		Base:           lipgloss.NewStyle(),
		Title:          lipgloss.NewStyle().Faint(true),
		Description:    lipgloss.NewStyle().Faint(true),
		ErrorIndicator: lipgloss.NewStyle(),
		ErrorMessage:   lipgloss.NewStyle(),

		// Select styles.
		SelectSelector: lipgloss.NewStyle(), // Selection indicator
		Option:         lipgloss.NewStyle(), // Select options
		NextIndicator:  lipgloss.NewStyle(),
		PrevIndicator:  lipgloss.NewStyle(),

		// FilePicker styles.
		Directory: lipgloss.NewStyle(),
		File:      lipgloss.NewStyle(),

		// Multi-select styles.
		MultiSelectSelector: lipgloss.NewStyle(),
		SelectedOption:      lipgloss.NewStyle(),
		SelectedPrefix:      lipgloss.NewStyle(),
		UnselectedOption:    lipgloss.NewStyle(),
		UnselectedPrefix:    lipgloss.NewStyle(),

		// Textinput and teatarea styles.
		TextInput: huh.TextInputStyles{
			Cursor:      lipgloss.NewStyle(),
			CursorText:  lipgloss.NewStyle(),
			Placeholder: lipgloss.NewStyle(),
			Prompt:      lipgloss.NewStyle(),
			Text:        lipgloss.NewStyle(),
		},

		// Confirm styles.
		FocusedButton: lipgloss.NewStyle(),
		BlurredButton: lipgloss.NewStyle(),

		// Card styles.
		Card:      lipgloss.NewStyle(),
		NoteTitle: lipgloss.NewStyle(),
		Next:      lipgloss.NewStyle(),
	},
	Focused: huh.FieldStyles{
		Base:           lipgloss.NewStyle(),
		Title:          lipgloss.NewStyle().Foreground(Teal).Bold(true),
		Description:    lipgloss.NewStyle().Faint(true),
		ErrorIndicator: lipgloss.NewStyle(),
		ErrorMessage:   lipgloss.NewStyle(),

		// Select styles.
		SelectSelector: lipgloss.NewStyle().SetString(icons.RightArrow + " ").Foreground(Teal).Bold(true), // Selection indicator
		Option:         lipgloss.NewStyle(),                                                               // Select options
		NextIndicator:  lipgloss.NewStyle(),
		PrevIndicator:  lipgloss.NewStyle(),

		// FilePicker styles.
		Directory: lipgloss.NewStyle(),
		File:      lipgloss.NewStyle(),

		// Multi-select styles.
		MultiSelectSelector: lipgloss.NewStyle(),
		SelectedOption:      lipgloss.NewStyle().Foreground(Teal).Bold(true),
		SelectedPrefix:      lipgloss.NewStyle().SetString(icons.RightArrow).Foreground(Teal).Bold(true),
		UnselectedOption:    lipgloss.NewStyle(),
		UnselectedPrefix:    lipgloss.NewStyle(),

		// Textinput and teatarea styles.
		TextInput: huh.TextInputStyles{
			Cursor:      lipgloss.NewStyle(),
			CursorText:  lipgloss.NewStyle(),
			Placeholder: lipgloss.NewStyle().Faint(true),
			Prompt:      lipgloss.NewStyle(),
			Text:        lipgloss.NewStyle(),
		},

		// Confirm styles.
		FocusedButton: lipgloss.NewStyle(),
		BlurredButton: lipgloss.NewStyle(),

		// Card styles.
		Card:      lipgloss.NewStyle(),
		NoteTitle: lipgloss.NewStyle(),
		Next:      lipgloss.NewStyle(),
	},
}
