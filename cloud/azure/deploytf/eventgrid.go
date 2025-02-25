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

package deploytf

type EventGridSubscriber struct {
	Url                       *string `json:"url"`
	ActiveDirectoryAppIdOrUri *string `json:"active_directory_app_id_or_uri"`
	ActiveDirectoryTenantId   *string `json:"active_directory_tenant_id"`
	EventToken                *string `json:"event_token"`
}
