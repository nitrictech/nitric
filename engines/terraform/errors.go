package terraform

import "fmt"

var (
	ErrPlatformNotFound = fmt.Errorf("platform not found")
	ErrUnauthenticated  = fmt.Errorf("unauthenticated")
)
