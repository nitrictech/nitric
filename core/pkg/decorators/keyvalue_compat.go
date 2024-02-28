package decorators

import (
	"context"

	keyvaluepb "github.com/nitrictech/nitric/core/pkg/proto/keyvalue/v1"
	kvstorepb "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
)

// KvStoreCompat - Compatibility layer for the KvStore service
type KvStoreCompat struct {
	kvstorepb.KvStoreServer
}

var _ keyvaluepb.KeyValueServer = (*KvStoreCompat)(nil)

// Get an existing value
func (k *KvStoreCompat) Get(ctx context.Context, req *keyvaluepb.KeyValueGetRequest) (*keyvaluepb.KeyValueGetResponse, error) {
	resp, err := k.GetKey(ctx, &kvstorepb.KvStoreGetRequest{
		Ref: &kvstorepb.ValueRef{
			Key:   req.Ref.Key,
			Store: req.Ref.Store,
		},
	})

	if err != nil {
		return nil, err
	}

	return &keyvaluepb.KeyValueGetResponse{
		Value: &keyvaluepb.Value{
			Ref: &keyvaluepb.ValueRef{
				Key:   resp.Value.Ref.Key,
				Store: resp.Value.Ref.Store,
			},
			Content: resp.Value.Content,
		},
	}, nil
}

// Create a new or overwrite an existing value
func (k *KvStoreCompat) Set(ctx context.Context, req *keyvaluepb.KeyValueSetRequest) (*keyvaluepb.KeyValueSetResponse, error) {
	_, err := k.SetKey(ctx, &kvstorepb.KvStoreSetRequest{
		Ref: &kvstorepb.ValueRef{
			Key:   req.Ref.Key,
			Store: req.Ref.Store,
		},
	})

	if err != nil {
		return nil, err
	}

	return &keyvaluepb.KeyValueSetResponse{}, nil
}

// Delete a key and its value
func (k *KvStoreCompat) Delete(ctx context.Context, req *keyvaluepb.KeyValueDeleteRequest) (*keyvaluepb.KeyValueDeleteResponse, error) {
	_, err := k.DeleteKey(ctx, &kvstorepb.KvStoreDeleteRequest{
		Ref: &kvstorepb.ValueRef{
			Key:   req.Ref.Key,
			Store: req.Ref.Store,
		},
	})

	if err != nil {
		return nil, err
	}

	return &keyvaluepb.KeyValueDeleteResponse{}, nil
}

func KeyValueServerWithCompat(srv kvstorepb.KvStoreServer) *KvStoreCompat {
	return &KvStoreCompat{
		KvStoreServer: srv,
	}
}
