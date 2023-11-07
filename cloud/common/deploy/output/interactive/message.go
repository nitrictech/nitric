// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interactive

type ResourceState string

const (
	Created  ResourceState = "Created"
	Creating ResourceState = "Creating"
	Updating ResourceState = "Updating"
	Updated  ResourceState = "Updated"
	Deleting ResourceState = "Deleting"
	Deleted  ResourceState = "Deleted"
)

type LogMessage struct {
	Message string
}

type LogMessageSubscriptionWriter struct {
	Sub chan LogMessage
}

func (l LogMessageSubscriptionWriter) Write(b []byte) (int, error) {
	l.Sub <- LogMessage{
		Message: string(b),
	}

	return len(b), nil
}
