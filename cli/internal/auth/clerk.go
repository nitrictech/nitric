package auth

type ClerkAuth struct {
}

var _ Auth = (*ClerkAuth)(nil)

func NewClerkAuth() *ClerkAuth {
	return &ClerkAuth{}
}
