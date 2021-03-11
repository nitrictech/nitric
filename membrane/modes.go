package membrane

import (
	"fmt"
	"strings"
)

// SourceType enum
type Mode int

const (
	Mode_Faas Mode = iota
	// PROXY Mode is designed for integration of monoliths into a nitric application
	Mode_HttpProxy
)

var modes = [...]string{"FAAS", "HTTP_PROXY"}

func (m Mode) String() string {
	return modes[m]
}

func ModeFromString(modeString string) (Mode, error) {
	for i, mode := range modes {
		if mode == modeString {
			return Mode(i), nil
		}
	}
	return -1, fmt.Errorf("Invalid mode %s, supported modes are: %s", modeString, strings.Join(modes[:], ", "))
}
