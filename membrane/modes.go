package membrane

// SourceType enum
type Mode int

const (
	Mode_Faas Mode = iota
	// PROXY Mode is designed for integration of monoliths into a nitric application
	Mode_HttpProxy
)

func (m Mode) String() string {
	return []string{"FAAS", "HTTP_PROXY"}[m]
}
