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

package api

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type domainNameArgs struct {
	domainName string
	api        *apigatewayv2.Api
	stage      *apigatewayv2.Stage
}

type domainName struct {
	pulumi.ResourceState

	Name string
}

func newDomainName(ctx *pulumi.Context, name string, args domainNameArgs) (*domainName, error) {
	domainParts := strings.Split(args.domainName, ".")

	res := &domainName{Name: name}

	err := ctx.RegisterComponentResource("nitric:api:DomainName", fmt.Sprintf("%s-%s", name, args.domainName), res)
	if err != nil {
		return nil, err
	}

	defaultOptions := []pulumi.ResourceOption{pulumi.Parent(res)}

	// Treat this domain as root by default
	baseName := ""
	// attempt to find hosted zone as the root domain name
	hostedZone, err := route53.LookupZone(ctx, &route53.LookupZoneArgs{
		// The name is the base name for the domain
		Name: &args.domainName,
	})
	if err != nil {
		// try by parent domain instead
		parentDomain := strings.Join(domainParts[1:], ".")
		hostedZone, err = route53.LookupZone(ctx, &route53.LookupZoneArgs{
			// The name is the base name for the domain
			Name: &parentDomain,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to find Route53 hosted zone to create records in: %w", err)
		}

		baseName = domainParts[0]
	}

	cert, err := acm.NewCertificate(ctx, fmt.Sprintf("%s-%s-cert", name, args.domainName), &acm.CertificateArgs{
		DomainName:       pulumi.String(args.domainName),
		ValidationMethod: pulumi.String("DNS"),
	}, defaultOptions...)
	if err != nil {
		return nil, err
	}

	domainValidationOption := cert.DomainValidationOptions.ApplyT(func(options []acm.CertificateDomainValidationOption) interface{} {
		return options[0]
	})

	certValidationDns, err := route53.NewRecord(ctx, fmt.Sprintf("%s-%s-certvalidationdns", name, args.domainName), &route53.RecordArgs{
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
		ZoneId: pulumi.String(hostedZone.ZoneId),
	}, defaultOptions...)
	if err != nil {
		return nil, err
	}

	certValidation, err := acm.NewCertificateValidation(ctx, fmt.Sprintf("%s-%s-certvalidation", name, args.domainName), &acm.CertificateValidationArgs{
		CertificateArn: cert.Arn,
		ValidationRecordFqdns: pulumi.StringArray{
			certValidationDns.Fqdn,
		},
	}, defaultOptions...)
	if err != nil {
		return nil, err
	}

	// Create a domain name if one has been requested
	apiDomainName, err := apigatewayv2.NewDomainName(ctx, fmt.Sprintf("%s-%s", name, args.domainName), &apigatewayv2.DomainNameArgs{
		DomainName: pulumi.String(args.domainName),
		DomainNameConfiguration: &apigatewayv2.DomainNameDomainNameConfigurationArgs{
			EndpointType:   pulumi.String("REGIONAL"),
			SecurityPolicy: pulumi.String("TLS_1_2"),
			CertificateArn: certValidation.CertificateArn,
		},
	}, defaultOptions...)
	if err != nil {
		return nil, err
	}

	// Create an API mapping for the new domain name
	_, err = apigatewayv2.NewApiMapping(ctx, fmt.Sprintf("%s-%s", name, args.domainName), &apigatewayv2.ApiMappingArgs{
		ApiId:      args.api.ID(),
		DomainName: apiDomainName.DomainName,
		Stage:      args.stage.Name,
	}, append(defaultOptions, pulumi.DependsOn([]pulumi.Resource{args.stage}))...)
	if err != nil {
		return nil, err
	}

	// Create a DNS record for the domain name that maps to the APIs
	// regional endpoint
	_, err = route53.NewRecord(ctx, fmt.Sprintf("%s-%s-dnsrecord", name, args.domainName), &route53.RecordArgs{
		ZoneId: pulumi.String(hostedZone.ZoneId),
		Type:   pulumi.String("A"),
		Name:   pulumi.String(baseName),
		Aliases: &route53.RecordAliasArray{
			&route53.RecordAliasArgs{
				// The target of the A record
				Name:                 apiDomainName.DomainNameConfiguration.TargetDomainName().Elem(),
				ZoneId:               apiDomainName.DomainNameConfiguration.HostedZoneId().Elem(),
				EvaluateTargetHealth: pulumi.Bool(false),
			},
		},
	}, defaultOptions...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
