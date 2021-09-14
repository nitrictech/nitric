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

package ses_service

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/nitric-dev/membrane/pkg/plugins/emails"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/utils"
)

const CharSet = "UTF-8"

type SesEmailService struct {
	emails.UnimplementedEmailService
	client sesiface.SESAPI
}

func destToAwsDest(dest emails.EmailDestination) (*ses.Destination, error) {
	if len(dest.To)+len(dest.Cc)+len(dest.Bcc) == 0 {
		return nil, fmt.Errorf("at least one destination email address (to, cc, or bcc) is required")
	}

	destination := ses.Destination{}
	// To Addresses
	if len(dest.To) > 0 {
		destination.ToAddresses = make([]*string, 0, len(dest.To))
		for i := 0; i < len(dest.To); i++ {
			destination.ToAddresses = append(destination.ToAddresses, aws.String(dest.To[i]))
		}
	}
	// Cc Addresses
	if len(dest.Cc) > 0 {
		destination.CcAddresses = make([]*string, 0, len(dest.Cc))
		for i := 0; i < len(dest.Cc); i++ {
			destination.CcAddresses = append(destination.CcAddresses, aws.String(dest.Cc[i]))
		}
	}
	// Bcc Addresses
	if len(dest.Bcc) > 0 {
		destination.BccAddresses = make([]*string, 0, len(dest.Bcc))
		for i := 0; i < len(dest.Bcc); i++ {
			destination.BccAddresses = append(destination.BccAddresses, aws.String(dest.Bcc[i]))
		}
	}

	return &destination, nil
}

// Send an email
func (s *SesEmailService) Send(from string, dest emails.EmailDestination, subject string, body emails.EmailBody) error {
	newErr := errors.ErrorsWithScope(
		"SesEmailService.Send",
		fmt.Sprintf("from=%s", from),
	)

	destination, err := destToAwsDest(dest)
	if err != nil {
		return newErr(
			codes.InvalidArgument,
			"failed to determine message destination",
			err,
		)
	}

	sendInput := &ses.SendEmailInput{
		Source:      aws.String(from),
		Destination: destination,
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    body.Html,
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    body.Text,
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
	}

	_, err = s.client.SendEmail(sendInput)

	// TODO: Consider sending message id in response.
	//result, err := s.client.SendEmail(sendInput)
	//result.MessageId

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return newErr(
					codes.InvalidArgument,
					"email message rejected",
					aerr,
				)
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return newErr(
					codes.FailedPrecondition,
					"email from domain not verified",
					aerr,
				)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				// TODO: consider FailedPrecondition.
				return newErr(
					codes.Internal,
					"configuration set does not exist",
					aerr,
				)
			default:
				return newErr(
					codes.Internal,
					"failed to send email",
					aerr,
				)
			}
		} else {
			return newErr(
				codes.Unknown,
				"failed to send email",
				aerr,
			)
		}
	}
	return nil
}

// New - Create a new SES email service plugin
func New() (emails.EmailService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	sess, sessionError := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %v", sessionError)
	}

	sesClient := ses.New(sess)

	return &SesEmailService{
		client: sesClient,
	}, nil
}

func NewWithClient(client sesiface.SESAPI) (emails.EmailService, error) {
	return &SesEmailService{
		client: client,
	}, nil
}
