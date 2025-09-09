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
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"golang.org/x/exp/maps"
)

type ImageArgs struct {
	SourceImage   string
	Runtime       []byte
	RepositoryUrl pulumi.StringInput
}

type LocalImageArgs struct {
	SourceImage   string
	SourceImageID string
	RepositoryUrl pulumi.StringInput
}

type Image struct {
	pulumi.ResourceState

	Name        string
	DockerImage *docker.RegistryImage
}

type WrappedBuildInput struct {
	Args       map[string]string
	Dockerfile string
}

var (
	//go:embed wrapper.dockerfile
	imageWrapper string
	//go:embed dummy.dockerfile
	dummyImageWrapper string
)

func NewLocalImage(ctx *pulumi.Context, name string, args *LocalImageArgs, opts ...pulumi.ResourceOption) (*Image, error) {
	res := &Image{Name: name}

	defaultOpts := append([]pulumi.ResourceOption{pulumi.Parent(res)}, opts...)

	err := ctx.RegisterComponentResource("nitriccommon:LocalImage", name, res, opts...)
	if err != nil {
		return nil, err
	}

	buildContext := fmt.Sprintf("%s/build-local-%s", os.TempDir(), name)
	err = os.MkdirAll(buildContext, 0o750)
	if err != nil {
		return nil, err
	}

	dockerfilePath := filepath.Clean(path.Join(buildContext, "Dockerfile"))
	if !strings.HasPrefix(dockerfilePath, os.TempDir()) {
		return nil, fmt.Errorf("unsafe dockerfile location")
	}
	dockerfile, err := os.Create(dockerfilePath)
	if err != nil {
		return nil, err
	}

	_, err = dockerfile.Write([]byte(dummyImageWrapper))
	if err != nil {
		return nil, err
	}
	err = dockerfile.Close()
	if err != nil {
		return nil, err
	}

	image, err := docker.NewImage(ctx, name+"-image", &docker.ImageArgs{
		ImageName: args.RepositoryUrl,
		Build: docker.DockerBuildArgs{
			Context:    pulumi.String(buildContext),
			Dockerfile: pulumi.String(path.Join(buildContext, "Dockerfile")),
			Args: pulumi.StringMap{
				"SOURCE_IMAGE":    pulumi.String(args.SourceImage),
				"SOURCE_IMAGE_ID": pulumi.String(args.SourceImageID),
			},
			Platform: pulumi.String("linux/amd64"),
		},
		SkipPush: pulumi.Bool(true),
	}, defaultOpts...)
	if err != nil {
		return nil, err
	}

	res.DockerImage, err = docker.NewRegistryImage(ctx, name+"-image", &docker.RegistryImageArgs{
		Name: image.ImageName,
		Triggers: pulumi.StringMap{
			"hash": image.RepoDigest,
		},
	}, defaultOpts...)
	if err != nil {
		return nil, err
	}

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":     pulumi.String(res.Name),
		"imageUri": pulumi.Sprintf("%s@%s", res.DockerImage.Name, res.DockerImage.Sha256Digest),
	})
}

func NewImage(ctx *pulumi.Context, name string, args *ImageArgs, opts ...pulumi.ResourceOption) (*Image, error) {
	res := &Image{Name: name}

	defaultOpts := append([]pulumi.ResourceOption{pulumi.Parent(res)}, opts...)

	err := ctx.RegisterComponentResource("nitriccommon:Image", name, res, opts...)
	if err != nil {
		return nil, err
	}

	imageWrapper, err := getWrapperDockerfile()
	if err != nil {
		return nil, err
	}

	dockerfileContent, sourceImageID, err := wrapDockerImage(imageWrapper.Dockerfile, args.SourceImage)
	if err != nil {
		return nil, err
	}

	buildContext := fmt.Sprintf("%s/build-%s", os.TempDir(), name)
	// Set Read/Write/Execute permissions for owner and group in compliance with https://securego.io/docs/rules/g301.html
	err = os.MkdirAll(buildContext, 0o750)
	if err != nil {
		return nil, err
	}

	// Address: https://securego.io/docs/rules/g304.html
	dockerfilePath := filepath.Clean(path.Join(buildContext, "Dockerfile"))
	if !strings.HasPrefix(dockerfilePath, os.TempDir()) {
		return nil, fmt.Errorf("unsafe dockerfile location")
	}
	dockerfile, err := os.Create(dockerfilePath)
	if err != nil {
		return nil, err
	}

	_, err = dockerfile.Write([]byte(dockerfileContent))
	if err != nil {
		return nil, err
	}
	err = dockerfile.Close()
	if err != nil {
		return nil, err
	}

	runtimefilePath := filepath.Clean(path.Join(buildContext, "runtime"))
	if !strings.HasPrefix(dockerfilePath, os.TempDir()) {
		return nil, fmt.Errorf("unsafe runtime location")
	}
	runtimefile, err := os.Create(runtimefilePath)
	if err != nil {
		return nil, err
	}

	_, err = runtimefile.Write(args.Runtime)
	if err != nil {
		return nil, err
	}
	err = runtimefile.Close()
	if err != nil {
		return nil, err
	}

	buildArgs := combineBuildArgs(map[string]string{
		"BASE_IMAGE":    args.SourceImage,
		"RUNTIME_FILE":  "runtime",
		"BASE_IMAGE_ID": sourceImageID,
	}, imageWrapper.Args)

	image, err := docker.NewImage(ctx, name+"-image", &docker.ImageArgs{
		ImageName: args.RepositoryUrl,
		Build: docker.DockerBuildArgs{
			Context:    pulumi.String(buildContext),
			Dockerfile: pulumi.String(path.Join(buildContext, "Dockerfile")),
			Args:       buildArgs,
			Platform:   pulumi.String("linux/amd64"),
		},
		SkipPush: pulumi.Bool(true),
	}, defaultOpts...)
	if err != nil {
		return nil, err
	}

	res.DockerImage, err = docker.NewRegistryImage(ctx, name+"-image", &docker.RegistryImageArgs{
		Name: image.ImageName,
		Triggers: pulumi.StringMap{
			"hash": image.RepoDigest,
		},
	}, defaultOpts...)
	if err != nil {
		return nil, err
	}

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":     pulumi.String(res.Name),
		"imageUri": pulumi.Sprintf("%s@%s", res.DockerImage.Name, res.DockerImage.Sha256Digest),
	})
}

func (d *Image) URI() pulumi.StringOutput {
	return pulumi.Sprintf("%s@%s", d.DockerImage.Name, d.DockerImage.Sha256Digest)
}

// Returns the default docker file if telemetry sampling is disabled for this service. Otherwise, will return a wrapped telemetry image.
func getWrapperDockerfile() (*WrappedBuildInput, error) {
	return &WrappedBuildInput{
		Dockerfile: imageWrapper,
		Args:       map[string]string{},
	}, nil
}

func combineBuildArgs(baseArgs, wrapperArgs map[string]string) pulumi.StringMap {
	maps.Copy(wrapperArgs, baseArgs)

	return pulumi.ToStringMap(wrapperArgs)
}

// Wraps the source image with the wrapper image, acknowledging the command from the source image
func wrapDockerImage(wrapper, sourceImage string) (string, string, error) {
	if sourceImage == "" {
		return "", "", fmt.Errorf("blank sourceImage provided")
	}

	inspectResult, err := CommandFromImageInspect(sourceImage, ",")
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf(wrapper, inspectResult.Cmd), inspectResult.ID, nil
}

type ImageInspect struct {
	ID      string
	Cmd     string
	WorkDir string
}

// Gets the command from the source image and returns as a comma separated string
func CommandFromImageInspect(sourceImage string, delimeter string) (*ImageInspect, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	imageInspect, _, err := client.ImageInspectWithRaw(context.Background(), sourceImage)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("could not inspect image: %s", sourceImage))
	}

	// Get the original command of the source image
	// and inject it into the wrapper image
	cmd := append(imageInspect.Config.Entrypoint, imageInspect.Config.Cmd...)

	// Wrap each command in string quotes
	cmdStr := []string{}
	for _, c := range cmd {
		cmdStr = append(cmdStr, "\""+c+"\"")
	}

	return &ImageInspect{
		ID:      imageInspect.ID,
		Cmd:     strings.Join(cmdStr, delimeter),
		WorkDir: imageInspect.Config.WorkingDir,
	}, nil
}
