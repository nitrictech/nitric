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

package worker

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/valyala/fasthttp"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

var _ = Describe("HttpWorker", func() {
	Context("Http", func() {
		rWrkr := &HttpWorker{
			port: 3000,
		}

		When("calling HandlesTrigger with arbitrary path", func() {
			It("should return true", func() {
				Expect(rWrkr.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{
							Method: "GET",
							Path:   "/test/",
						},
					},
				})).To(BeTrue())
			})
		})

		When("calling HandlesTrigger with arbitrary method", func() {
			It("should return true", func() {
				Expect(rWrkr.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{
							Method: "POST",
							Path:   "/test/test",
						},
					},
				})).To(BeTrue())
			})
		})

		When("calling HandleHttpRequest", func() {
			// Run on a non-blocking thread
			go func() {
				defer GinkgoRecover()

				err := fasthttp.ListenAndServe(":3000", func(ctx *fasthttp.RequestCtx) {
					ctx.SuccessString("text/plain", "success")
				})

				Expect(err).ToNot(HaveOccurred())
			}()

			// Delay to allow the HTTP server to correctly start
			// FIXME: Should block on channels...
			time.Sleep(500 * time.Millisecond)

			It("should rewrite X-Forwarded-Authorization", func() {
				ctrl := gomock.NewController(GinkgoT())

				httpWrkr := &HttpWorker{
					port: 3000,
				}

				resp, err := httpWrkr.HandleTrigger(context.TODO(), &v1.TriggerRequest{
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{
							Method: "GET",
							Path:   "/test/name",
							Headers: map[string]*v1.HeaderValue{
								"X-Forwarded-Authorization": {
									Value: []string{"Bearer 1234"},
								},
							},
						},
					},
				})

				Expect(err).ShouldNot(HaveOccurred())

				Expect(string(resp.GetData())).To(Equal("success"))
				ctrl.Finish()
			})
		})
	})

	Context("Event", func() {
		eventTrigger := &v1.TriggerRequest{
			Context: &v1.TriggerRequest_Topic{
				Topic: &v1.TopicTriggerContext{},
			},
		}
		When("calling HandlesTrigger wth an Event", func() {
			rWrkr := &HttpWorker{}

			It("should return false", func() {
				Expect(rWrkr.HandlesTrigger(eventTrigger)).To(BeFalse())
			})
		})

		When("calling HandleEvent", func() {
			subWrkr := &HttpWorker{}

			It("should return an error", func() {
				_, err := subWrkr.HandleTrigger(context.TODO(), eventTrigger)

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
