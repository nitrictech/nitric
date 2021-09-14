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
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/golang/mock/gomock"
	mock_sesiface "github.com/nitric-dev/membrane/mocks/aws_ses"
	"github.com/nitric-dev/membrane/pkg/plugins/emails"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SES", func() {
	When("Send", func() {
		When("Given a valid set of email request attributes", func() {
			from := "from@example.com"
			to := "to@example.com"
			subject := "the subject"
			textBody := "the body"
			htmlBody := "<p>the html body</p>"

			When("Sending the email", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSes := mock_sesiface.NewMockSESAPI(crtl)
				emailsPlugin, _ := NewWithClient(mockSes)
				It("Should successfully provide the to address to the SES API", func() {
					By("Calling send with the correct arguments")
					mockSes.EXPECT().SendEmail(&ses.SendEmailInput{
						Source: aws.String("from@example.com"),
						Destination: &ses.Destination{
							ToAddresses:  []*string{aws.String("to@example.com")},
							CcAddresses:  nil,
							BccAddresses: nil,
						},
						Message: &ses.Message{
							Body: &ses.Body{
								Html: &ses.Content{
									Charset: aws.String(CharSet),
									Data:    aws.String("<p>the html body</p>"),
								},
								Text: &ses.Content{
									Charset: aws.String(CharSet),
									Data:    aws.String("the body"),
								},
							},
							Subject: &ses.Content{
								Charset: aws.String(CharSet),
								Data:    aws.String("the subject"),
							},
						},
					}).Times(1).Return(nil, nil)

					err := emailsPlugin.Send(from, emails.EmailDestination{To: []string{to}}, subject, emails.EmailBody{
						Text: &textBody,
						Html: &htmlBody,
					})
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		When("Given no destination addresses", func() {
			from := "from@example.com"
			subject := "the subject"
			textBody := "the body"
			htmlBody := "<p>the html body</p>"

			When("Sending the email", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSes := mock_sesiface.NewMockSESAPI(crtl)
				emailsPlugin, _ := NewWithClient(mockSes)
				It("Should successfully provide the to address to the SES API", func() {
					By("Calling not calling send on the SES API")

					err := emailsPlugin.Send(from, emails.EmailDestination{To: []string{}}, subject, emails.EmailBody{
						Text: &textBody,
						Html: &htmlBody,
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})

		When("AWS SES returns an error", func() {
			from := "from@example.com"
			to := "to@example.com"
			subject := "the subject"
			textBody := "the body"
			htmlBody := "<p>the html body</p>"

			When("Sending the email", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSes := mock_sesiface.NewMockSESAPI(crtl)
				emailsPlugin, _ := NewWithClient(mockSes)
				It("Should successfully provide the to address to the SES API", func() {
					By("Calling SES API")
					mockSes.EXPECT().SendEmail(&ses.SendEmailInput{
						Source: aws.String("from@example.com"),
						Destination: &ses.Destination{
							ToAddresses:  []*string{aws.String("to@example.com")},
							CcAddresses:  nil,
							BccAddresses: nil,
						},
						Message: &ses.Message{
							Body: &ses.Body{
								Html: &ses.Content{
									Charset: aws.String(CharSet),
									Data:    aws.String("<p>the html body</p>"),
								},
								Text: &ses.Content{
									Charset: aws.String(CharSet),
									Data:    aws.String("the body"),
								},
							},
							Subject: &ses.Content{
								Charset: aws.String(CharSet),
								Data:    aws.String("the subject"),
							},
						},
					}).Times(1).Return(nil, awserr.New(ses.ErrCodeMessageRejected, "test error", fmt.Errorf("internal test error")))

					err := emailsPlugin.Send(from, emails.EmailDestination{To: []string{to}}, subject, emails.EmailBody{
						Text: &textBody,
						Html: &htmlBody,
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})
})
