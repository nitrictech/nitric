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

package codes

import "fmt"

type Code int

const (
	OK                 Code = 0
	Cancelled          Code = 1
	Unknown            Code = 2
	InvalidArgument    Code = 3
	DeadlineExceeded   Code = 4
	NotFound           Code = 5
	AlreadyExists      Code = 6
	PermissionDenied   Code = 7
	ResourceExhausted  Code = 8
	FailedPrecondition Code = 9
	Aborted            Code = 10
	OutOfRange         Code = 11
	Unimplemented      Code = 12
	Internal           Code = 13
	Unavailable        Code = 14
	DataLoss           Code = 15
	Unauthenticated    Code = 16
)

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case Cancelled:
		return "Cancelled"
	case Unknown:
		return "Unknown"
	case InvalidArgument:
		return "Invalid Argument"
	case DeadlineExceeded:
		return "Deadline Exceeded"
	case NotFound:
		return "Not Found"
	case AlreadyExists:
		return "Already Exists"
	case PermissionDenied:
		return "Permission Denied"
	case ResourceExhausted:
		return "Resource Exhausted"
	case FailedPrecondition:
		return "Failed Precondition"
	case Aborted:
		return "Aborted"
	case OutOfRange:
		return "Out of Range"
	case Unimplemented:
		return "Unimplemented"
	case Internal:
		return "Internal"
	case Unavailable:
		return "Unavailable"
	case DataLoss:
		return "Data Loss"
	case Unauthenticated:
		return "Unauthenticated"
	default:
		return fmt.Sprintf("Unknown error code: %d", c)
	}
}
