package email_service

import (
	"fmt"

	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"

	"github.com/nitric-dev/membrane/pkg/plugins/emails"
)

type DevEmailService struct {
	emails.UnimplementedEmailService
}

func isNilOrEmpty(strArray []string) bool {
	return strArray == nil || len(strArray) == 0
}

func (s *DevEmailService) Send(from string, dest emails.EmailDestination, subject string, body emails.EmailBody) error {
	newErr := errors.ErrorsWithScope(
		"SesEmailService.Send",
		fmt.Sprintf("from=%s", from),
	)

	// At least one destination email address is required.
	if isNilOrEmpty(dest.To) && isNilOrEmpty(dest.Cc) && isNilOrEmpty(dest.Bcc) {
		return newErr(codes.InvalidArgument, "no recipients specified", fmt.Errorf("one of to, cc or bcc must be specified"))
	}

	// Just print the request so basic debugging/testing can be performed.
	// Future enhancement may allow SMTP setting or similar be used to send real emails.
	fmt.Println(fmt.Sprintf("Email sent from: %s, destination: %v, subject: %s, body text: %s, body html: %s", from, dest, subject, *body.Text, *body.Html))
	return nil
}

func New() (emails.EmailService, error) {
	return &DevEmailService{}, nil
}
