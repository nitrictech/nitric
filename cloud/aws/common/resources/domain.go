package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func GetSubdomainNameLabel(domainName string) string {
	domainParts := strings.Split(domainName, ".")
	if len(domainParts) > 2 {
		return domainParts[0]
	}

	return ""
}

func GetZoneID(domainName string) (string, error) {
	zoneIds := GetZoneIDs([]string{domainName})
	if zoneIds[domainName] == "" {
		return "", fmt.Errorf("zone ID not found for domain name: %s", domainName)
	}

	return zoneIds[domainName], nil
}

func GetZoneIDs(domainNames []string) map[string]string {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		return nil
	}

	client := route53.NewFromConfig(cfg)

	zoneMap := make(map[string]string)

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
			zoneMap[domain] = id
			continue
		}

		// Try parent/root domain
		parts := strings.Split(domain, ".")
		if len(parts) > 2 {
			root := strings.Join(parts[len(parts)-2:], ".")
			if id, ok := hostedZones[root]; ok {
				zoneMap[domain] = id
				continue
			}
		}
	}

	return zoneMap
}
