package deploy

// Stack - represents a collection of related functions and their shared dependencies.
type Stack struct {
	// A stack can be composed of one or more applications
	functions []*Function
}

// Listen - Starts server to listen for new resource registrations for the stack
func (s *Stack) Listen() {

}
