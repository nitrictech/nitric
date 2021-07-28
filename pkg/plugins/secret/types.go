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

package secret

// Secret - Represents a container for secret versions
type Secret struct {
	Name string
}

// SecretVersion - A version of a secret
type SecretVersion struct {
	Secret  *Secret
	Version string
}

// SecretAccessResponse - Return value for a secret access request
type SecretAccessResponse struct {
	SecretVersion *SecretVersion
	Value         []byte
}

// SecretPutResponse - Return value for a secret put request
type SecretPutResponse struct {
	SecretVersion *SecretVersion
}
