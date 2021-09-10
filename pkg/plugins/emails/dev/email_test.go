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

package email_service_test

import (
	"github.com/nitric-dev/membrane/pkg/plugins/emails"
	email_service "github.com/nitric-dev/membrane/pkg/plugins/emails/dev"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Email", func() {

	emailPlugin, err := email_service.New()
	if err != nil {
		panic(err)
	}

	Context("Send", func() {
		text := "some text"
		html := "<h1>some html</h1>"
		body := emails.EmailBody{
			Text: &text,
			Html: &html,
		}

		When("A 'to' address is set", func() {
			dest := emails.EmailDestination{
				To: []string{"someoneelse@example.com"},
			}

			It("Should send the email", func() {
				err := emailPlugin.Send("someone@example.com", dest, "some subject", body)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("A 'cc' address is set", func() {
			dest := emails.EmailDestination{
				Cc: []string{"someoneelse@example.com"},
			}

			It("Should send the email", func() {
				err := emailPlugin.Send("someone@example.com", dest, "some subject", body)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("A 'bcc' address is set", func() {
			dest := emails.EmailDestination{
				Bcc: []string{"someoneelse@example.com"},
			}

			It("Should send the email", func() {
				err := emailPlugin.Send("someone@example.com", dest, "some subject", body)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("A no destination addresses are set", func() {
			dest := emails.EmailDestination{}

			It("Should error", func() {
				err := emailPlugin.Send("someone@example.com", dest, "some subject", body)
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
