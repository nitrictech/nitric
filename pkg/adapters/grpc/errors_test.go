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

package grpc_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/pkg/adapters/grpc"
	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
)

type SecretValue struct {
	Type   string `log:"Type"`
	Factor int    `log:"-"`
	Value  string
}

type Secret struct {
	Name    string       `json:"Name" log:"Name"`
	Version string       `json:"Version" log:"Version"`
	Value   *SecretValue `json:"Value" log:"Value"`
}

var _ = Describe("GRPC Errors", func() {
	Context("GrpcError", func() {
		When("plugin.errors.InvalidArgument", func() {
			It("Should report GRPC IllegalArgument error", func() {
				newErr := errors.ErrorsWithScope("test", nil)
				err := newErr(
					codes.InvalidArgument,
					"bad param",
					nil,
				)
				grpcErr := grpc.NewGrpcError("BadServer.BadCall", err)
				Expect(grpcErr.Error()).To(ContainSubstring("rpc error: code = InvalidArgument desc = bad param"))
			})
		})
		When("plugin.errors.InvalidArgument args", func() {
			It("Should report GRPC IllegalArgument error with args", func() {
				args := map[string]interface{}{"key": "value"}
				newErr := errors.ErrorsWithScope("test", args)
				err := newErr(
					codes.InvalidArgument,
					"bad param",
					nil,
				)
				grpcErr := grpc.NewGrpcError("BadServer.BadCall", err)
				Expect(grpcErr.Error()).To(ContainSubstring("rpc error: code = InvalidArgument desc = bad param"))
			})
		})
		When("Standard Error", func() {
			It("Should report GRPC Internal error", func() {
				err := fmt.Errorf("internal error")
				err = grpc.NewGrpcError("BadServer.BadCall", err)
				Expect(err.Error()).To(ContainSubstring("rpc error: code = Internal desc = internal error"))
			})
		})
	})

	Context("PluginNotRegisteredError", func() {
		When("Creating a New PluginNotRegisteredError", func() {
			It("Should contain the name of the plugin", func() {
				err := grpc.NewPluginNotRegisteredError("Document")
				Expect(err.Error()).To(ContainSubstring("rpc error: code = Unimplemented desc = Document plugin not registered"))
			})
		})
	})

	Context("Logging Arg", func() {
		When("string", func() {
			It("return string value", func() {
				Expect(grpc.LogArg("string")).To(BeEquivalentTo("string"))
			})
		})

		When("int", func() {
			It("return string value", func() {
				Expect(grpc.LogArg(123)).To(BeEquivalentTo("123"))
			})
		})

		When("bool", func() {
			It("return string value", func() {
				Expect(grpc.LogArg(true)).To(BeEquivalentTo("true"))
			})
		})

		When("float", func() {
			It("return string value", func() {
				Expect(grpc.LogArg(3.1415)).To(BeEquivalentTo("3.1415"))
			})
		})

		When("struct", func() {
			It("return string value", func() {
				data := Secret{
					Name:    "name",
					Version: "3",
					Value: &SecretValue{
						Type:   "key",
						Factor: 2,
						Value:  "2a4wijgPq0PpwJ76IjT7&lTBZ$5SGRcq",
					},
				}

				value := grpc.LogArg(data)
				Expect(value).To(BeEquivalentTo("{Name: name, Version: 3, Value: {Type: key}}"))
			})
		})

		When("map", func() {
			It("return string value", func() {
				secret := Secret{
					Name:    "name",
					Version: "3",
					Value: &SecretValue{
						Type:   "key",
						Factor: 2,
						Value:  "2a4wijgPq0PpwJ76IjT7&lTBZ$5SGRcq",
					},
				}

				valueMap := map[string]interface{}{
					"key":    "value",
					"secret": secret,
				}
				value := grpc.LogArg(valueMap)
				Expect(value).To(ContainSubstring("secret: {Name: name, Version: 3, Value: {Type: key}}"))
				Expect(value).To(ContainSubstring("key: value"))
			})
		})
	})
})
