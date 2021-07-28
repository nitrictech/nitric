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

	"github.com/nitric-dev/membrane/pkg/adapters/grpc"
	"github.com/nitric-dev/membrane/pkg/plugins"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GRPC Utils", func() {
	Context("GrpcError", func() {
		When("sdk.IllegalArgumentError", func() {
			It("Should report GRPC IllegalArgument error", func() {
				err := plugins.NewInvalidArgError("bad param")
				err = grpc.NewGrpcError("BadServer.BadCall", err)
				Expect(err.Error()).To(ContainSubstring("rpc error: code = InvalidArgument desc = BadServer.BadCall: bad param"))
			})
		})
		When("Standard Error", func() {
			It("Should report GRPC Internal error", func() {
				err := fmt.Errorf("internal error")
				err = grpc.NewGrpcError("BadServer.BadCall", err)
				Expect(err.Error()).To(ContainSubstring("rpc error: code = Internal desc = BadServer.BadCall: internal error"))
			})
		})
	})
})
