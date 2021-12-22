package deploy

// Stack
type Stack struct {
	// A stack can be composed of one or more applications
	apps []*App
}

// Listen - Starts server to listen for new resource registrations for the stack
func (s *Stack) Listen() {

}
