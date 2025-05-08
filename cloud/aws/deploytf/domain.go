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

package deploytf

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v19/acmcertificate"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v19/acmcertificatevalidation"
	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v19/provider"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v19/route53record"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/common/resources"
)

type Domain struct {
	Name           string
	ZoneId         string
	CertificateArn *string
}

func newTerraformDomain(tfstack cdktf.TerraformStack, domainName string) (*Domain, error) {
	zoneLookup, err := resources.GetZoneID(domainName)
	if err != nil {
		return nil, err
	}

	// ACM Provider in us-east-1
	acmProvider := awsprovider.NewAwsProvider(tfstack, jsii.String("AWSUsEast1"), &awsprovider.AwsProviderConfig{
		Region: jsii.String("us-east-1"),
		Alias:  jsii.String("us-east-1"),
	})

	// ACM Certificate (must be in us-east-1)
	cert := acmcertificate.NewAcmCertificate(tfstack, jsii.String("CdnCert"), &acmcertificate.AcmCertificateConfig{
		DomainName:       jsii.String(domainName),
		ValidationMethod: jsii.String("DNS"),
		Provider:         acmProvider, // Ensure ACM is deployed in us-east-1
		Lifecycle: &cdktf.TerraformResourceLifecycle{
			CreateBeforeDestroy: jsii.Bool(true),
		},
	})

	// Route 53 Record for DNS validation (remains in the main region)
	validationRecord := route53record.NewRoute53Record(tfstack, jsii.String("CdnCertValidation"), &route53record.Route53RecordConfig{
		ZoneId: jsii.String(zoneLookup.ZoneID),
		Name:   cert.DomainValidationOptions().Get(jsii.Number(0)).ResourceRecordName(),
		Type:   cert.DomainValidationOptions().Get(jsii.Number(0)).ResourceRecordType(),
		Records: &[]*string{
			cert.DomainValidationOptions().Get(jsii.Number(0)).ResourceRecordValue(),
		},
		Ttl: jsii.Number(10 * 60),
		DependsOn: &[]cdktf.ITerraformDependable{
			cert,
		},
	})

	// ACM Certificate Validation (must be in us-east-1)
	validation := acmcertificatevalidation.NewAcmCertificateValidation(tfstack, jsii.String("CertValidation"), &acmcertificatevalidation.AcmCertificateValidationConfig{
		CertificateArn: cert.Arn(),
		ValidationRecordFqdns: &[]*string{
			validationRecord.Fqdn(),
		},
		Provider: acmProvider, // Use us-east-1 provider
	})

	return &Domain{
		Name:           domainName,
		ZoneId:         zoneLookup.ZoneID,
		CertificateArn: validation.CertificateArn(),
	}, nil
}
