package dev

import (
	"fmt"

	"github.com/nitric-dev/membrane/pkg/plugins/emails"
)

type DevEmailService struct {
	emails.UnimplementedEmailService
}

func (s *DevEmailService) Send(from string, dest emails.EmailDestination, subject string, body emails.EmailBody) error {
	// Just print the request so basic debugging/testing can be performed.
	// Future enhancement may allow SMTP setting or similar be used to send real emails.
	fmt.Println(fmt.Sprintf("Send Email from: %s, destination: %v, subject: %s, body text: %s, body html: %s", from, dest, subject, body.Text, body.Html))
	// TODO: consider mechanism to simulate errors
	return nil
}

func New() (emails.EmailService, error) {
	return &DevEmailService{}, nil
}
