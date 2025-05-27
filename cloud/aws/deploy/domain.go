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
	"fmt"
	"slices"

	awsprovider "github.com/pulumi/pulumi-aws/sdk/v5/go/aws"

	"github.com/nitrictech/nitric/cloud/aws/common/resources"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Domain struct {
	pulumi.ResourceState

	Name                  string
	ZoneLookup            *resources.ZoneLookup
	CertificateValidation *acm.CertificateValidation
}

type domainArgs struct {
	DomainName string
	// Required for backwards compatibility with provider versions < 1.26.1
	AliasName string
	// If the domain is used for a CDN, it will be deployed in the us-east-1 region
	IsCDNDomain bool
}

func (a *NitricAwsPulumiProvider) newPulumiDomainName(ctx *pulumi.Context, args domainArgs) (*Domain, error) {
	var err error
	res := &Domain{Name: args.DomainName}

	res.ZoneLookup, err = resources.GetZoneID(args.DomainName)
	if err != nil {
		return nil, err
	}

	err = ctx.RegisterComponentResource("nitric:api:DomainName", fmt.Sprintf("%s-%s", args.DomainName, a.StackId), res)
	if err != nil {
		return nil, err
	}

	defaultOptions := []pulumi.ResourceOption{pulumi.Parent(res)}

	// Create an AWS provider for the us-east-1 region as the acm certificates require being deployed in us-east-1 region
	if args.IsCDNDomain && a.Region != "us-east-1" {
		useast1, err := awsprovider.NewProvider(ctx, "us-east-1", &awsprovider.ProviderArgs{
			Region: pulumi.String("us-east-1"),
		})
		if err != nil {
			return nil, err
		}

		defaultOptions = append(defaultOptions, pulumi.Provider(useast1))
	}

	cert, err := acm.NewCertificate(ctx, fmt.Sprintf("cert-%s", a.StackId), &acm.CertificateArgs{
		DomainName:       pulumi.String(args.DomainName),
		ValidationMethod: pulumi.String("DNS"),
	},
		slices.Concat(defaultOptions, []pulumi.ResourceOption{pulumi.Aliases([]pulumi.Alias{
			// Required for backwards compatibility with provider versions < 1.26.1
			{Name: pulumi.String(fmt.Sprintf("%s-%s-cert", args.AliasName, args.DomainName))},
		})})...,
	)
	if err != nil {
		return nil, err
	}

	domainValidationOption := cert.DomainValidationOptions.ApplyT(func(options []acm.CertificateDomainValidationOption) interface{} {
		return options[0]
	})

	cdnRecord, err := route53.NewRecord(ctx, fmt.Sprintf("cdn-record-%s", a.StackId), &route53.RecordArgs{
		Name: domainValidationOption.ApplyT(func(option interface{}) string {
			return *option.(acm.CertificateDomainValidationOption).ResourceRecordName
		}).(pulumi.StringOutput),
		Type: domainValidationOption.ApplyT(func(option interface{}) string {
			return *option.(acm.CertificateDomainValidationOption).ResourceRecordType
		}).(pulumi.StringOutput),
		Records: pulumi.StringArray{
			domainValidationOption.ApplyT(func(option interface{}) string {
				return *option.(acm.CertificateDomainValidationOption).ResourceRecordValue
			}).(pulumi.StringOutput),
		},
		Ttl:    pulumi.Int(10 * 60),
		ZoneId: pulumi.String(res.ZoneLookup.ZoneID),
	}, []pulumi.ResourceOption{
		pulumi.Parent(res),
		pulumi.Aliases([]pulumi.Alias{
			// Required for backwards compatibility with provider versions < 1.26.1
			{Name: pulumi.String(fmt.Sprintf("%s-%s-certvalidationdns", args.AliasName, args.DomainName))},
		}),
	}...)
	if err != nil {
		return nil, err
	}

	res.CertificateValidation, err = acm.NewCertificateValidation(ctx, fmt.Sprintf("cert-validation-%s", a.StackId), &acm.CertificateValidationArgs{
		CertificateArn: cert.Arn,
		ValidationRecordFqdns: pulumi.StringArray{
			cdnRecord.Fqdn,
		},
	},
		slices.Concat(defaultOptions, []pulumi.ResourceOption{
			pulumi.Aliases([]pulumi.Alias{
				// Required for backwards compatibility with provider versions < 1.26.1
				{Name: pulumi.String(fmt.Sprintf("%s-%s-certvalidation", args.AliasName, args.DomainName))},
			}),
		})...,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
