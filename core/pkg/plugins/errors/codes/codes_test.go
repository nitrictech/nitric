// Copyright 2021 Nitric Technologies Pty Ltd.
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

package codes_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
)

var codeStringMap map[codes.Code]string = map[codes.Code]string{
	codes.OK:                 "OK",
	codes.Cancelled:          "Cancelled",
	codes.Unknown:            "Unknown",
	codes.InvalidArgument:    "Invalid Argument",
	codes.DeadlineExceeded:   "Deadline Exceeded",
	codes.NotFound:           "Not Found",
	codes.AlreadyExists:      "Already Exists",
	codes.PermissionDenied:   "Permission Denied",
	codes.ResourceExhausted:  "Resource Exhausted",
	codes.FailedPrecondition: "Failed Precondition",
	codes.Aborted:            "Aborted",
	codes.OutOfRange:         "Out of Range",
	codes.Unimplemented:      "Unimplemented",
	codes.Internal:           "Internal",
	codes.Unavailable:        "Unavailable",
	codes.DataLoss:           "Data Loss",
	codes.Unauthenticated:    "Unauthenticated",
}

var _ = Describe("gRPC Codes", func() {
	for code, desc := range codeStringMap {
		When(fmt.Sprintf("converting grpc code %d to String()", code), func() {
			It(fmt.Sprintf("should print %s", desc), func() {
				Expect(codes.Aborted.String()).To(Equal(desc))
			})
		})
	}
})
