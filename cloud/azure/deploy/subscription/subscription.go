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

package subscription

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/nitrictech/nitric/cloud/azure/deploy/policy"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	pulumiNativeEventgrid "github.com/pulumi/pulumi-azure-native-sdk/eventgrid"
	pulumiEventgrid "github.com/pulumi/pulumi-azure/sdk/v4/go/azure/eventgrid"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SubscriptionArgs struct {
	ResourceGroupName  pulumi.StringInput
	Subscriptions      map[string]*pulumiNativeEventgrid.Topic
	Sp                 *policy.ServicePrincipal
	AppName            string
	LatestRevisionFqdn pulumi.StringOutput
}

type Subscription struct {
	pulumi.ResourceState

	Name string
}

func NewSubscription(ctx *pulumi.Context, name string, args *SubscriptionArgs, opts ...pulumi.ResourceOption) (*Subscription, error) {
	res := &Subscription{Name: name}

	err := ctx.RegisterComponentResource("nitric:api:AzureSubscription", name, res, opts...)
	if err != nil {
		return nil, err
	}

	if len(args.Subscriptions) == 0 {
		return nil, nil
	}

	hostUrl := args.LatestRevisionFqdn.ApplyT(func(fqdn string) (string, error) {
		_ = ctx.Log.Info("waiting for "+args.AppName+" to start before creating subscriptions", &pulumi.LogArgs{Ephemeral: true})

		// Get the full URL of the deployed container
		hostUrl := "https://" + fqdn

		hCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// Poll the URL until the host has started.
		for {
			// Provide data in the expected shape. The content is current not important.
			empty := ""
			dummyEvgt := eventgrid.Event{
				ID:          &empty,
				Data:        &empty,
				EventType:   &empty,
				Subject:     &empty,
				DataVersion: &empty,
			}

			jsonStr, err := dummyEvgt.MarshalJSON()
			if err != nil {
				return "", err
			}

			body := bytes.NewBuffer(jsonStr)

			req, err := http.NewRequestWithContext(hCtx, "POST", hostUrl, body)
			if err != nil {
				return "", err
			}

			// TODO: Implement a membrane health check handler in the Membrane and trigger that instead.
			// Set event type header to simulate a subscription validation event.
			// These events are automatically resolved by the Membrane and won't be processed by handlers.
			req.Header.Set("aeg-event-type", "SubscriptionValidation")
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			resp, err := client.Do(req)
			if err == nil {
				resp.Body.Close()
				break
			}
		}

		return hostUrl, nil
	}).(pulumi.StringOutput)

	_ = ctx.Log.Info("creating subscriptions for "+args.AppName, &pulumi.LogArgs{})

	for subName, sub := range args.Subscriptions {
		_, err = pulumiEventgrid.NewEventSubscription(ctx, utils.ResourceName(ctx, args.AppName+"-"+subName, utils.EventSubscriptionRT), &pulumiEventgrid.EventSubscriptionArgs{
			Scope: sub.ID(),
			WebhookEndpoint: pulumiEventgrid.EventSubscriptionWebhookEndpointArgs{
				Url: hostUrl,
				// TODO: Reduce event chattiness here and handle internally in the Azure AppService HTTP Gateway?
				MaxEventsPerBatch:         pulumi.Int(1),
				ActiveDirectoryAppIdOrUri: args.Sp.ClientID,
				ActiveDirectoryTenantId:   args.Sp.TenantID,
			},
			RetryPolicy: pulumiEventgrid.EventSubscriptionRetryPolicyArgs{
				MaxDeliveryAttempts: pulumi.Int(30),
				EventTimeToLive:     pulumi.Int(5),
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
