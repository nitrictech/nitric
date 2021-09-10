// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func isEmpty(strArray []string) bool {
	return len(strArray) == 0
}

func (s *DevEmailService) Send(from string, dest emails.EmailDestination, subject string, body emails.EmailBody) error {
	newErr := errors.ErrorsWithScope(
		"DevEmailService.Send",
		fmt.Sprintf("from=%s", from),
	)

	// At least one destination email address is required.
	if isEmpty(dest.To) && isEmpty(dest.Cc) && isEmpty(dest.Bcc) {
		return newErr(codes.InvalidArgument, "no recipients specified", fmt.Errorf("one of to, cc or bcc must be specified"))
	}

	// Just print the request so basic debugging/testing can be performed.
	// Future enhancement may allow SMTP setting or similar be used to send real emails.
	fmt.Printf("Email sent from: %s, destination: %v, subject: %s, body text: %s, body html: %s \n", from, dest, subject, *body.Text, *body.Html)

	return nil
}

func New() (emails.EmailService, error) {
	return &DevEmailService{}, nil
}
