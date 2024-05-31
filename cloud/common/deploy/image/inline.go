// Copyright Nitric Pty Ltd.

// SPDX-License-Identifier: Apache-2.0

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package image

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type BuildWrappedImageArgs struct {
	ServiceName string
	SourceImage string
	TargetImage string
	Runtime     []byte
}

var runtimeFileName = "runtime"

// Build a wrapped image on the current machine inline
func BuildWrappedImage(args *BuildWrappedImageArgs) (string, error) {
	imageWrapper, err := getWrapperDockerfile(nil)
	if err != nil {
		return "", err
	}

	dockerfileContent, sourceImageID, err := wrapDockerImage(imageWrapper.Dockerfile, args.SourceImage)
	if err != nil {
		return "", err
	}

	buildContext := fmt.Sprintf("%s/build-%s", os.TempDir(), args.ServiceName)
	err = os.MkdirAll(buildContext, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Address: https://securego.io/docs/rules/g304.html
	dockerfilePath := filepath.Clean(path.Join(buildContext, "Dockerfile"))
	if !strings.HasPrefix(dockerfilePath, os.TempDir()) {
		return "", fmt.Errorf("unsafe dockerfile location")
	}
	dockerfile, err := os.Create(dockerfilePath)
	if err != nil {
		return "", err
	}

	_, err = dockerfile.Write([]byte(dockerfileContent))
	if err != nil {
		return "", err
	}
	err = dockerfile.Close()
	if err != nil {
		return "", err
	}

	runtimefilePath := filepath.Clean(path.Join(buildContext, runtimeFileName))
	if !strings.HasPrefix(dockerfilePath, os.TempDir()) {
		return "", fmt.Errorf("unsafe runtime location")
	}
	runtimefile, err := os.Create(runtimefilePath)
	if err != nil {
		return "", err
	}

	_, err = runtimefile.Write(args.Runtime)
	if err != nil {
		return "", err
	}
	err = runtimefile.Close()
	if err != nil {
		return "", err
	}

	tarBuildContext, err := archive.TarWithOptions(buildContext, &archive.TarOptions{})
	if err != nil {
		return "", err
	}

	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	buildArgs := map[string]*string{
		"BASE_IMAGE":    &args.SourceImage,
		"RUNTIME_FILE":  &runtimeFileName,
		"BASE_IMAGE_ID": &sourceImageID,
	}

	response, err := client.ImageBuild(context.TODO(), tarBuildContext, types.ImageBuildOptions{
		BuildArgs:  buildArgs,
		Tags:       []string{args.TargetImage},
		Dockerfile: "Dockerfile",
	})
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		return "", err
	}

	// Return the target image ID
	inspect, _, err := client.ImageInspectWithRaw(context.Background(), args.TargetImage)
	if err != nil {
		return "", err
	}

	return inspect.ID, nil
}
