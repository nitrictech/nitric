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

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

type ZoneLookup struct {
	// The domain that matched the Hosted Zone lookup
	Domain string
	// The Hosted Zone ID
	ZoneID string
	// If the zone matched the domain (false) or matched the parent (true)
	IsParent bool
}

func GetARecordLabel(zoneLookup *ZoneLookup) string {
	if !zoneLookup.IsParent {
		return ""
	}

	return getSubdomainLabel(zoneLookup.Domain)
}

func getSubdomainLabel(domain string) string {
	domainParts := strings.Split(domain, ".")
	if len(domainParts) > 2 {
		return domainParts[0]
	}

	return ""
}

func GetZoneID(domainName string) (*ZoneLookup, error) {
	zoneIds := GetZoneIDs([]string{domainName})
	if zoneIds[domainName] == nil {
		return nil, fmt.Errorf("zone ID not found for domain name: %s", domainName)
	}

	return zoneIds[domainName], nil
}

func GetZoneIDs(domainNames []string) map[string]*ZoneLookup {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		return nil
	}

	client := route53.NewFromConfig(cfg)

	zoneMap := make(map[string]*ZoneLookup)

	normalizedDomains := make(map[string]string)
	for _, d := range domainNames {
		d = strings.ToLower(strings.TrimSuffix(d, "."))
		normalizedDomains[d] = d + "."
	}

	paginator := route53.NewListHostedZonesPaginator(client, &route53.ListHostedZonesInput{})
	hostedZones := make(map[string]string) // map of zone name -> zone ID

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil
		}

		for _, hz := range page.HostedZones {
			name := strings.ToLower(strings.TrimSuffix(*hz.Name, "."))
			hostedZones[name] = strings.TrimPrefix(*hz.Id, "/hostedzone/")
		}
	}

	// Resolve each domain name
	for domain, normalized := range normalizedDomains {
		// Check full domain
		if id, ok := hostedZones[strings.TrimSuffix(normalized, ".")]; ok {
			zoneMap[domain] = &ZoneLookup{
				Domain:   domain,
				ZoneID:   id,
				IsParent: false,
			}
			continue
		}

		// Try parent/root domain
		parts := strings.Split(domain, ".")
		if len(parts) > 2 {
			root := strings.Join(parts[len(parts)-2:], ".")
			if id, ok := hostedZones[root]; ok {
				zoneMap[domain] = &ZoneLookup{
					Domain:   domain,
					ZoneID:   id,
					IsParent: true,
				}
				continue
			}
		}

		zoneMap[domain] = nil
	}

	return zoneMap
}
