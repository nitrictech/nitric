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

package deploy

import (
	"regexp"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	autoNameLength = 8
)

type ResourceType struct {
	Abbreviation   string
	MaxLen         int
	AllowUpperCase bool
	AllowHyphen    bool
	UseName        bool
}

// https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations
// https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules
var (
	notAlphaNumericRegexp = regexp.MustCompile("[^a-zA-Z0-9-]+")

	// Alphanumerics, underscores, parentheses, hyphens, periods, and unicode characters that match the regex documentation.
	// Can't end with period. Regex pattern: ^[-\w\._\(\)]+$
	ResourceGroupRT = ResourceType{Abbreviation: "rg", MaxLen: 90, AllowUpperCase: true, AllowHyphen: true}

	ContainerAppRT = ResourceType{Abbreviation: "app", MaxLen: 32, UseName: true, AllowHyphen: true}
	// Alphanumerics
	RegistryRT = ResourceType{Abbreviation: "cr", MaxLen: 50, AllowUpperCase: true}
	// Alphanumerics and hyphens. Start and end with alphanumeric.
	AnalyticsWorkspaceRT = ResourceType{Abbreviation: "log", MaxLen: 24, AllowHyphen: true}
	AssignmentRT         = ResourceType{Abbreviation: "assign", MaxLen: 64, UseName: true}
	// TODO find docs on this..
	KubeRT = ResourceType{Abbreviation: "kube", MaxLen: 64, AllowUpperCase: true}
	// lowercase letters, numbers, and the '-' character, and must be between 3 and 50 characters.
	CosmosDBAccountRT = ResourceType{Abbreviation: "cosmos", MaxLen: 50, AllowHyphen: true}
	// TODO find requirements
	MongoDBRT = ResourceType{Abbreviation: "mongo", MaxLen: 24, AllowUpperCase: true}
	// TODO find requirements
	MongoCollectionRT            = ResourceType{Abbreviation: "coll", MaxLen: 24, AllowUpperCase: true, UseName: true}
	ADApplicationRT              = ResourceType{Abbreviation: "aad-app", MaxLen: 64, UseName: true}
	ADServicePrincipalRT         = ResourceType{Abbreviation: "aad-sp", MaxLen: 64, UseName: true}
	ADServicePrincipalPasswordRT = ResourceType{Abbreviation: "aad-spp", MaxLen: 64, UseName: true}
	// Lowercase letters and numbers.
	StorageAccountRT = ResourceType{Abbreviation: "st", MaxLen: 24}
	// 	Lowercase letters, numbers, and hyphens.
	// Start with lowercase letter or number. Can't use consecutive hyphens.
	StorageContainerRT = ResourceType{MaxLen: 63, AllowHyphen: true, UseName: true}
	// Lowercase letters, numbers, and hyphens.
	// Can't start or end with hyphen. Can't use consecutive hyphens.
	StorageQueueRT = ResourceType{MaxLen: 63, AllowHyphen: true, UseName: true}

	// Alphanumerics and hyphens. Start with letter. End with letter or digit. Can't contain consecutive hyphens.
	KeyVaultRT = ResourceType{Abbreviation: "kv", MaxLen: 14, AllowUpperCase: true}

	// Alphanumerics and hyphens.
	EventGridRT = ResourceType{Abbreviation: "evgt", MaxLen: 24, AllowUpperCase: true, AllowHyphen: true, UseName: true}

	// Alphanumerics and hyphens.
	EventSubscriptionRT = ResourceType{Abbreviation: "sub", MaxLen: 24, AllowUpperCase: true, AllowHyphen: true, UseName: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiRT = ResourceType{Abbreviation: "api", MaxLen: 80, AllowHyphen: true, AllowUpperCase: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiHttpProxyRT = ResourceType{Abbreviation: "httpproxy", MaxLen: 80, AllowHyphen: true, AllowUpperCase: true, UseName: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiManagementRT = ResourceType{Abbreviation: "api-mgmt", MaxLen: 80, AllowHyphen: true, AllowUpperCase: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiManagementServiceRT = ResourceType{Abbreviation: "api-mgmt", MaxLen: 50, AllowHyphen: true, AllowUpperCase: true}

	ApiManagementProxyRT = ResourceType{Abbreviation: "httpproxy-mgmt", MaxLen: 80, AllowHyphen: true, AllowUpperCase: true, UseName: true}

	// Alphanumerics and hyphens, Start with letter and end with alphanumeric.
	ApiOperationPolicyRT = ResourceType{Abbreviation: "api-op-pol", MaxLen: 80, AllowUpperCase: true, AllowHyphen: true, UseName: true}
)

// cleanNameSegment removes all non-alphanumeric characters from a string.
// also removes hyphens if they're not permitted by the particular resource type.
func cleanNameSegment(p string, rt ResourceType) string {
	r := notAlphaNumericRegexp.ReplaceAllString(p, "")
	if !rt.AllowHyphen {
		r = strings.ReplaceAll(r, "-", "")
	}

	return r
}

// withoutBlanks returns a new copy of the string array with all blank strings removed.
func withoutBlanks(strs []string) []string {
	newStrs := []string{}

	for _, s := range strs {
		if s != "" {
			newStrs = append(newStrs, s)
		}
	}

	return newStrs
}

// ResourceName generates a name for the deployed version of a resource in Azure.
// follows restrictions like max length, hyphenation, etc.
func ResourceName(ctx *pulumi.Context, name string, rt ResourceType) string {
	var parts []string

	maxLen := rt.MaxLen - autoNameLength
	abbrLen := len(rt.Abbreviation)

	if rt.AllowHyphen {
		abbrLen += 1
	}

	if rt.UseName {
		parts = []string{
			StringTrunc(cleanNameSegment(name, rt), maxLen-abbrLen),
			rt.Abbreviation,
		}
	} else {
		deployName := strings.TrimPrefix(ctx.Stack(), ctx.Project()+"-")
		partLen := (maxLen - abbrLen) / 2
		parts = []string{
			StringTrunc(cleanNameSegment(ctx.Project(), rt), partLen),
			StringTrunc(cleanNameSegment(deployName, rt), partLen),
			rt.Abbreviation,
		}
	}

	parts = withoutBlanks(parts)

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

	rname := StringTrunc(s, maxLen)
	if strings.Trim(rname, "") == "" {
		panic("generated blank resource name")
	}

	return rname
}
