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

package uri

import (
	"fmt"
	"net/url"
	"strings"
)

type NitricUri struct {
	uri url.URL
}

func (n NitricUri) ResourceType() string {
	return strings.Split(n.uri.Opaque, ":")[0]
}

func (n NitricUri) ResourceName() string {
	return strings.Split(n.uri.Opaque, ":")[1]
}

func (n NitricUri) Query() map[string][]string {
	return n.uri.Query()
}

func New(uri url.URL) (*NitricUri, error) {
	if uri.Scheme != "nitric" {
		return nil, fmt.Errorf("provided url is not a nitric uri")
	}

	if uri.Opaque != "" && len(strings.Split(uri.Opaque, ":")) != 2 {
		return nil, fmt.Errorf("malformed nitric uri")
	}

	return &NitricUri{
		uri: uri,
	}, nil
}
