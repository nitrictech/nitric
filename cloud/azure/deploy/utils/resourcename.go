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

package utils

import (
	"regexp"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	autoNameLength = 7
)

type ResouceType struct {
	Abbreviation   string
	MaxLen         int
	AllowUpperCase bool
	AllowHyphen    bool
	UseName        bool
}

// https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations
// https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules
var (
	alphanumeric = regexp.MustCompile("[^a-zA-Z0-9-]+")

	// Alphanumerics, underscores, parentheses, hyphens, periods, and unicode characters that match the regex documentation.
	// Can't end with period. Regex pattern: ^[-\w\._\(\)]+$
	ResourceGroupRT = ResouceType{Abbreviation: "rg", MaxLen: 90, AllowUpperCase: true, AllowHyphen: true}

	ContainerAppRT = ResouceType{Abbreviation: "app", MaxLen: 64, UseName: true, AllowHyphen: true}
	// Alphanumerics
	RegistryRT = ResouceType{Abbreviation: "cr", MaxLen: 50, AllowUpperCase: true}
	// Alphanumerics and hyphens. Start and end with alphanumeric.
	AnalyticsWorkspaceRT = ResouceType{Abbreviation: "log", MaxLen: 24, AllowHyphen: true}
	AssignmentRT         = ResouceType{Abbreviation: "assign", MaxLen: 64, UseName: true}
	// TODO find docs on this..
	KubeRT = ResouceType{Abbreviation: "kube", MaxLen: 64, AllowUpperCase: true}
	// lowercase letters, numbers, and the '-' character, and must be between 3 and 50 characters.
	CosmosDBAccountRT = ResouceType{Abbreviation: "cosmos", MaxLen: 50, AllowHyphen: true}
	// TODO find requirements
	MongoDBRT = ResouceType{Abbreviation: "mongo", MaxLen: 24, AllowUpperCase: true}
	// TODO find requirements
	MongoCollectionRT            = ResouceType{Abbreviation: "coll", MaxLen: 24, AllowUpperCase: true, UseName: true}
	ADApplicationRT              = ResouceType{Abbreviation: "aad-app", MaxLen: 64, UseName: true}
	ADServicePrincipalRT         = ResouceType{Abbreviation: "aad-sp", MaxLen: 64, UseName: true}
	ADServicePrincipalPasswordRT = ResouceType{Abbreviation: "aad-spp", MaxLen: 64, UseName: true}
	// Lowercase letters and numbers.
	StorageAccountRT = ResouceType{Abbreviation: "st", MaxLen: 24}
	// 	Lowercase letters, numbers, and hyphens.
	// Start with lowercase letter or number. Can't use consecutive hyphens.
	StorageContainerRT = ResouceType{MaxLen: 63, AllowHyphen: true, UseName: true}
	// Lowercase letters, numbers, and hyphens.
	// Can't start or end with hyphen. Can't use consecutive hyphens.
	StorageQueueRT = ResouceType{MaxLen: 63, AllowHyphen: true, UseName: true}

	// Alphanumerics and hyphens. Start with letter. End with letter or digit. Can't contain consecutive hyphens.
	KeyVaultRT = ResouceType{Abbreviation: "kv", MaxLen: 14, AllowUpperCase: true}

	// Alphanumerics and hyphens.
	EventGridRT = ResouceType{Abbreviation: "evgt", MaxLen: 24, AllowUpperCase: true, AllowHyphen: true, UseName: true}

	// Alphanumerics and hyphens.
	EventSubscriptionRT = ResouceType{Abbreviation: "sub", MaxLen: 24, AllowUpperCase: true, AllowHyphen: true, UseName: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiRT = ResouceType{Abbreviation: "api", MaxLen: 80, AllowHyphen: true, AllowUpperCase: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiHttpProxyRT = ResouceType{Abbreviation: "httpproxy", MaxLen: 80, AllowHyphen: true, AllowUpperCase: true, UseName: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiManagementRT = ResouceType{Abbreviation: "api-mgmt", MaxLen: 80, AllowHyphen: true, AllowUpperCase: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiManagementServiceRT = ResouceType{Abbreviation: "api-mgmt", MaxLen: 50, AllowHyphen: true, AllowUpperCase: true}

	ApiManagementProxyRT = ResouceType{Abbreviation: "httpproxy-mgmt", MaxLen: 80, AllowHyphen: true, AllowUpperCase: true, UseName: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiOperationPolicyRT = ResouceType{Abbreviation: "api-op-pol", MaxLen: 80, AllowUpperCase: true, AllowHyphen: true, UseName: true}
)

func cleanPart(p string, rt ResouceType) string {
	r := alphanumeric.ReplaceAllString(p, "")
	if !rt.AllowHyphen {
		r = strings.ReplaceAll(r, "-", "")
	}

	return r
}

func ignoreEmpty(strs []string) []string {
	newStrs := []string{}

	for _, s := range strs {
		if s != "" {
			newStrs = append(newStrs, s)
		}
	}

	return newStrs
}

func ResourceName(ctx *pulumi.Context, name string, rt ResouceType) string {
	var parts []string

	maxLen := rt.MaxLen - autoNameLength
	abbrLen := len(rt.Abbreviation)

	if rt.AllowHyphen {
		abbrLen += 1
	}

	if rt.UseName {
		parts = []string{
			StringTrunc(cleanPart(name, rt), maxLen-abbrLen),
			rt.Abbreviation,
		}
	} else {
		deployName := strings.TrimPrefix(ctx.Stack(), ctx.Project()+"-")
		partLen := (maxLen - abbrLen) / 2
		parts = []string{
			StringTrunc(cleanPart(ctx.Project(), rt), partLen),
			StringTrunc(cleanPart(deployName, rt), partLen),
			rt.Abbreviation,
		}
	}

	parts = ignoreEmpty(parts)

	// first char must be a letter
	parts[0] = strings.TrimLeft(parts[0], "0123456789-")

	var s string

	if rt.AllowHyphen {
		s = strings.Join(parts, "-")
		s = strings.ReplaceAll(s, "--", "-")
	} else if rt.AllowUpperCase {
		s = JoinCamelCase(parts)
	} else {
		s = strings.Join(parts, "")
	}

	if !rt.AllowHyphen {
		s = strings.ReplaceAll(s, "-", "")
	}

	if !rt.AllowUpperCase {
		s = strings.ToLower(s)
	}

	return StringTrunc(s, maxLen)
}
