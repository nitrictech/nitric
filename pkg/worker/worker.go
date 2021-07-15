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
	"bytes"
	"fmt"
	triggers2 "github.com/nitric-dev/membrane/pkg/triggers"
	"io/ioutil"
	"net/http"
)

type Worker interface {
	HandleEvent(trigger *triggers2.Event) error
	HandleHttpRequest(trigger *triggers2.HttpRequest) (*triggers2.HttpResponse, error)
}

type UnimplementedWorker struct{}

func (*UnimplementedWorker) HandleEvent(trigger *triggers2.Event) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedWorker) HandleHttpRequest(trigger *triggers2.HttpRequest) *http.Response {
	return &http.Response{
		Status:     "Unimplemented",
		StatusCode: 501,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("HTTP Handler Unimplemented"))),
	}
}
