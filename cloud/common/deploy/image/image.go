// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package image

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ImageArgs struct {
	SourceImage   string
	Runtime       []byte
	RepositoryUrl pulumi.StringInput
	Server        pulumi.StringInput
	Username      pulumi.StringInput
	Password      pulumi.StringInput
}

type Image struct {
	pulumi.ResourceState

	Name        string
	DockerImage *docker.Image
}

//go:embed wrapper.dockerfile
var imageWrapper string

func wrapDockerImage(wrapper string, sourceImage string) (string, error) {
	if sourceImage == "" {
		return "", fmt.Errorf("blank sourceImage provided")
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	ii, _, err := cli.ImageInspectWithRaw(context.Background(), sourceImage)
	if err != nil {
		return "", errors.WithMessage(err, fmt.Sprintf("could not inspect image: %s", sourceImage))
	}

	// Get the original command of the source image
	// and inject it into the wrapper image
	cmd := append(ii.Config.Entrypoint, ii.Config.Cmd...)

	// Wrap each command in string quotes
	cmdStr := []string{}
	for _, c := range cmd {
		cmdStr = append(cmdStr, "\""+c+"\"")
	}

	return fmt.Sprintf(wrapper, strings.Join(cmdStr, ",")), nil
}

func NewImage(ctx *pulumi.Context, name string, args *ImageArgs, opts ...pulumi.ResourceOption) (*Image, error) {
	res := &Image{Name: name}

	err := ctx.RegisterComponentResource("nitric:Image", name, res, opts...)
	if err != nil {
		return nil, err
	}

	// TODO: Need to re-add support for telemetry wrappers as well
	dockerfileContent, err := wrapDockerImage(imageWrapper, args.SourceImage)
	if err != nil {
		return nil, err
	}
	buildContext := fmt.Sprintf("%s/build-%s", os.TempDir(), name)
	os.MkdirAll(buildContext, os.ModePerm)

	dockerfile, err := os.Create(path.Join(buildContext, "Dockerfile"))
	if err != nil {
		return nil, err
	}
	dockerfile.Write([]byte(dockerfileContent))
	dockerfile.Close()

	runtimefile, err := os.Create(path.Join(buildContext, "runtime"))
	if err != nil {
		return nil, err
	}
	runtimefile.Write(args.Runtime)
	runtimefile.Close()

	res.DockerImage, err = docker.NewImage(ctx, name+"-image", &docker.ImageArgs{
		ImageName: args.RepositoryUrl,
		Registry: &docker.RegistryArgs{
			Server:   args.Server,
			Username: args.Username,
			Password: args.Password,
		},
		Build: docker.DockerBuildArgs{
			Context:    pulumi.String(buildContext),
			Dockerfile: pulumi.String(path.Join(buildContext, "Dockerfile")),
			Platform:   pulumi.String("linux/amd64"),
			Args: pulumi.StringMap{
				"BASE_IMAGE":   pulumi.String(args.SourceImage),
				"RUNTIME_FILE": pulumi.String("runtime"),
			},
		},
	}, append(opts, pulumi.Parent(res))...)
	if err != nil {
		return nil, err
	}

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":     pulumi.String(res.Name),
		"imageUri": res.DockerImage.ImageName,
	})
}

func (d *Image) URI() pulumi.StringOutput {
	return d.DockerImage.RepoDigest.Elem().ToStringOutput()
}
