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
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/firestore"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type FirestoreCollectionDatabase struct {
	pulumi.ResourceState

	Name     string
	Database *firestore.Database
}

type FirestoreCollectionDatabaseArgs struct {
	Location      string
	DefaultExists bool
}

func NewFirestoreCollectionDatabase(ctx *pulumi.Context, name string, args *FirestoreCollectionDatabaseArgs, opts ...pulumi.ResourceOption) (*FirestoreCollectionDatabase, error) {
	res := &FirestoreCollectionDatabase{
		Name: name,
	}
	err := ctx.RegisterComponentResource("nitric:collection:GCPFirestoreDatabase", name, res, opts...)
	if err != nil {
		return nil, err
	}
	defaultFirestoreId := pulumi.ID("(default)")

	// Attempt to locate the default database
	// This get appears to actually create a stack resource, so just proceed if we find it
	if args.DefaultExists {
		res.Database, err = firestore.GetDatabase(ctx, "default", defaultFirestoreId, nil)
		if err != nil {
			return nil, err
		}
	} else {
		res.Database, err = firestore.NewDatabase(ctx, "default", &firestore.DatabaseArgs{
			Name:                     defaultFirestoreId,
			AppEngineIntegrationMode: pulumi.String("DISABLED"),
			LocationId:               pulumi.String(args.Location),
			Type:                     pulumi.String("FIRESTORE_NATIVE"),
		}, pulumi.RetainOnDelete(true))
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
