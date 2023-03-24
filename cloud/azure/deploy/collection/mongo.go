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

package collection

import (
	"fmt"

	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/documentdb"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type MongoCollectionsArgs struct {
	ResourceGroup *resources.ResourceGroup
	Collections   []*v1.Resource
}

type MongoCollections struct {
	pulumi.ResourceState

	Name             string
	Account          *documentdb.DatabaseAccount
	MongoDB          *documentdb.MongoDBResourceMongoDBDatabase
	ConnectionString pulumi.StringOutput
	Collections      map[string]*documentdb.MongoDBResourceMongoDBCollection
}

func NewMongoCollections(ctx *pulumi.Context, name string, args *MongoCollectionsArgs, opts ...pulumi.ResourceOption) (*MongoCollections, error) {
	res := &MongoCollections{
		Name:        name,
		Collections: map[string]*documentdb.MongoDBResourceMongoDBCollection{},
	}

	err := ctx.RegisterComponentResource("nitric:collections:CosmosMongo", name, res, opts...)
	if err != nil {
		return nil, err
	}

	primaryGeo := documentdb.LocationArgs{
		FailoverPriority: pulumi.Int(0),
		IsZoneRedundant:  pulumi.Bool(false),
		LocationName:     args.ResourceGroup.Location,
	}
	secondaryGeo := documentdb.LocationArgs{
		FailoverPriority: pulumi.Int(1),
		IsZoneRedundant:  pulumi.Bool(false),
		// FIXME: Don't use defaults like this
		LocationName: pulumi.String("canadacentral"),
	}

	if primaryGeo.LocationName == secondaryGeo.LocationName {
		// FIXME: Don't use defaults like this
		secondaryGeo.LocationName = pulumi.String("northeurope")
	}

	res.Account, err = documentdb.NewDatabaseAccount(ctx, utils.ResourceName(ctx, name, utils.CosmosDBAccountRT), &documentdb.DatabaseAccountArgs{
		ResourceGroupName: args.ResourceGroup.Name,
		Kind:              pulumi.String("MongoDB"),

		ApiProperties: &documentdb.ApiPropertiesArgs{
			ServerVersion: pulumi.String("4.0"),
		},
		Location:                 args.ResourceGroup.Location,
		DatabaseAccountOfferType: documentdb.DatabaseAccountOfferTypeStandard.ToDatabaseAccountOfferTypeOutput(),
		Locations: documentdb.LocationArray{documentdb.LocationArgs{
			FailoverPriority: pulumi.IntPtr(0),
			IsZoneRedundant:  pulumi.BoolPtr(false),
			LocationName:     args.ResourceGroup.Location,
		}, documentdb.LocationArgs{
			FailoverPriority: pulumi.IntPtr(1),
			IsZoneRedundant:  pulumi.BoolPtr(false),
			LocationName:     pulumi.String("eastus"),
		}},
	}, pulumi.Parent(res))
	if err != nil {
		return nil, errors.WithMessage(err, "cosmosdb account")
	}

	res.MongoDB, err = documentdb.NewMongoDBResourceMongoDBDatabase(ctx, utils.ResourceName(ctx, name, utils.MongoDBRT), &documentdb.MongoDBResourceMongoDBDatabaseArgs{
		ResourceGroupName: args.ResourceGroup.Name,
		AccountName:       res.Account.Name,
		DatabaseName:      pulumi.String(name),
		Location:          args.ResourceGroup.Location,
		Resource: documentdb.MongoDBDatabaseResourceArgs{
			Id: pulumi.String(name),
		},
	}, pulumi.Parent(res))
	if err != nil {
		return nil, errors.WithMessage(err, "mongo db")
	}

	for _, col := range args.Collections {
		res.Collections[col.Name], err = documentdb.NewMongoDBResourceMongoDBCollection(ctx, utils.ResourceName(ctx, col.Name, utils.MongoCollectionRT), &documentdb.MongoDBResourceMongoDBCollectionArgs{
			ResourceGroupName: args.ResourceGroup.Name,
			AccountName:       res.Account.Name,
			DatabaseName:      res.MongoDB.Name,
			CollectionName:    pulumi.String(col.Name),
			Location:          res.MongoDB.Location,
			Options:           &documentdb.CreateUpdateOptionsArgs{},
			Resource: documentdb.MongoDBCollectionResourceArgs{
				Id: pulumi.String(col.Name),
			},
		}, pulumi.Parent(res))
		if err != nil {
			return nil, errors.WithMessage(err, "mongo collection")
		}
	}

	connectionString := pulumi.All(args.ResourceGroup.Name, res.Account.Name).ApplyT(func(args []interface{}) (string, error) {
		rgName := args[0].(string)
		acctName := args[1].(string)
		connStr, err := documentdb.ListDatabaseAccountConnectionStrings(ctx, &documentdb.ListDatabaseAccountConnectionStringsArgs{
			ResourceGroupName: rgName,
			AccountName:       acctName,
		})
		if err != nil {
			return "", err
		}

		if len(connStr.ConnectionStrings) == 0 {
			return "", fmt.Errorf("no avaialable db connection strings")
		}

		return connStr.ConnectionStrings[0].ConnectionString, nil
	}).(pulumi.StringOutput)

	res.ConnectionString = connectionString

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":              pulumi.String(res.Name),
		"mongoDatabaseName": res.MongoDB.Name,
		"connectionString":  connectionString,
	})
}
