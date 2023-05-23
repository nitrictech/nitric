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

package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"

	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

// Up - Deploy requested infrastructure for a stack
func (d *DeployServer) Up(request *deploy.DeployUpRequest, stream deploy.DeployService_UpServer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			err = fmt.Errorf("recovered panic: %+v\n Stack: %s", r, stack)
		}
	}()

	reqJson, err := json.MarshalIndent(request, "", "    ")
	if err != nil {
		return err
	}

	out, ok := request.Attributes.AsMap()["output"]

	if ok {
		if outFile, isString := out.(string); isString {
			absPath, err := filepath.Abs(outFile)
			if err != nil {
				return err
			}

			p := path.Dir(outFile)
			err = os.MkdirAll(p, os.ModePerm)

			if err != nil {
				return err
			}

			err = os.WriteFile(outFile, reqJson, os.ModePerm)
			if err != nil {
				return err
			}

			err = stream.Send(&deploy.DeployUpEvent{
				Content: &deploy.DeployUpEvent_Result{
					Result: &deploy.DeployUpEventResult{
						Success: true,
						Result: &deploy.UpResult{
							Content: &deploy.UpResult_StringResult{
								StringResult: fmt.Sprintf("spec written to: %s", absPath),
							},
						},
					},
				},
			})
			if err != nil {
				return err
			}
		}
	} else {
		err = stream.Send(&deploy.DeployUpEvent{
			Content: &deploy.DeployUpEvent_Result{
				Result: &deploy.DeployUpEventResult{
					Success: true,
					Result: &deploy.UpResult{
						Content: &deploy.UpResult_StringResult{
							StringResult: string(reqJson),
						},
					},
				},
			},
		})

		if err != nil {
			return err
		}
	}

	return nil
}
