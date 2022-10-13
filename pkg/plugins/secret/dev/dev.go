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

package secret_service

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/pkg/plugins/secret"
	"github.com/nitrictech/nitric/pkg/utils"
)

const DEV_SUB_DIRECTORY = "./secrets/"

type DevSecretService struct {
	secret.UnimplementedSecretPlugin
	secDir string
}

func (s *DevSecretService) secretFileName(sec *secret.Secret, v string) string {
	filename := fmt.Sprintf("%s_%s.txt", sec.Name, v)
	return filepath.Join(s.secDir, filename)
}

func (s *DevSecretService) Put(sec *secret.Secret, val []byte) (*secret.SecretPutResponse, error) {
	newErr := errors.ErrorsWithScope(
		"DevSecretService.Put",
		map[string]interface{}{
			"secret": sec,
		},
	)

	if sec == nil {
		return nil, newErr(codes.InvalidArgument, "provide non-empty secret", nil)
	}
	if len(sec.Name) == 0 {
		return nil, newErr(codes.InvalidArgument, "provide non-blank secret name", nil)
	}
	if len(val) == 0 {
		return nil, newErr(codes.InvalidArgument, "provide non-blank secret value", nil)
	}

	versionId := uuid.New().String()
	// Creates a new file in the form:
	// DIR/Name_Version.txt
	file, err := os.Create(s.secretFileName(sec, versionId))
	if err != nil {
		return nil, newErr(
			codes.FailedPrecondition,
			"error creating secret store",
			err,
		)
	}

	sVal := base64.StdEncoding.EncodeToString(val)

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(sVal)
	if err != nil {
		return nil, newErr(
			codes.FailedPrecondition,
			"error writing secret value",
			err,
		)
	}
	writer.Flush()

	// Creates a new file as latest
	latestFile, err := os.Create(s.secretFileName(sec, "latest"))
	if err != nil {
		return nil, newErr(
			codes.FailedPrecondition,
			"error creating latest secret",
			err,
		)
	}

	latestWriter := bufio.NewWriter(latestFile)
	_, err = latestWriter.WriteString(sVal + "," + versionId)
	if err != nil {
		return nil, newErr(
			codes.FailedPrecondition,
			"error writing secret value",
			err,
		)
	}

	latestWriter.Flush()

	return &secret.SecretPutResponse{
		SecretVersion: &secret.SecretVersion{
			Secret: &secret.Secret{
				Name: sec.Name,
			},
			Version: versionId,
		},
	}, nil
}

func (s *DevSecretService) Access(sv *secret.SecretVersion) (*secret.SecretAccessResponse, error) {
	newErr := errors.ErrorsWithScope(
		"DevSecretService.Access",
		map[string]interface{}{
			"version": sv,
		},
	)

	if sv.Secret.Name == "" {
		return nil, newErr(
			codes.InvalidArgument,
			"provide non-blank name",
			nil,
		)
	}
	if sv.Version == "" {
		return nil, newErr(
			codes.InvalidArgument,
			"provide non-blank version",
			nil,
		)
	}

	content, err := os.ReadFile(s.secretFileName(sv.Secret, sv.Version))
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"error reading secret store",
			err,
		)
	}

	splitContent := strings.Split(string(content), ",")
	version := sv.Version
	// check whether a version number is stored in the file, this indicates the 'latest' version file.
	if len(splitContent) == 2 {
		version = splitContent[1]
	}

	sVal, err := base64.StdEncoding.DecodeString(splitContent[0])
	if err != nil {
		return nil, err
	}

	return &secret.SecretAccessResponse{
		SecretVersion: &secret.SecretVersion{
			Secret: &secret.Secret{
				Name: sv.Secret.Name,
			},
			Version: version,
		},
		Value: sVal,
	}, nil
}

// Create new secret store
func New() (secret.SecretService, error) {
	secDir := utils.GetEnv("LOCAL_SEC_DIR", utils.GetRelativeDevPath(DEV_SUB_DIRECTORY))

	// Check whether file exists
	_, err := os.Stat(secDir)
	if os.IsNotExist(err) {
		// Make directory if not present
		err := os.MkdirAll(secDir, 0o777)
		if err != nil {
			return nil, err
		}
	}
	return &DevSecretService{
		secDir: secDir,
	}, nil
}
