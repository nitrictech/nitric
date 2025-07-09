// GENERATED CODE -- DO NOT EDIT!
'use strict';
var grpc = require('@grpc/grpc-js');
var proto_storage_v2_storage_pb = require('../../../proto/storage/v2/storage_pb.js');
var google_protobuf_duration_pb = require('google-protobuf/google/protobuf/duration_pb.js');
function serialize_nitric_proto_storage_v2_StorageDeleteRequest(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageDeleteRequest)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageDeleteRequest');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageDeleteRequest(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageDeleteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StorageDeleteResponse(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageDeleteResponse)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageDeleteResponse');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageDeleteResponse(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageDeleteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StorageExistsRequest(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageExistsRequest)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageExistsRequest');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageExistsRequest(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageExistsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StorageExistsResponse(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageExistsResponse)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageExistsResponse');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageExistsResponse(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageExistsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StorageListBlobsRequest(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageListBlobsRequest)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageListBlobsRequest');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageListBlobsRequest(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageListBlobsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StorageListBlobsResponse(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageListBlobsResponse)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageListBlobsResponse');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageListBlobsResponse(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageListBlobsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StoragePreSignUrlRequest(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StoragePreSignUrlRequest)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StoragePreSignUrlRequest');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StoragePreSignUrlRequest(buffer_arg) {
    return proto_storage_v2_storage_pb.StoragePreSignUrlRequest.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StoragePreSignUrlResponse(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StoragePreSignUrlResponse)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StoragePreSignUrlResponse');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StoragePreSignUrlResponse(buffer_arg) {
    return proto_storage_v2_storage_pb.StoragePreSignUrlResponse.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StorageReadRequest(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageReadRequest)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageReadRequest');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageReadRequest(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageReadRequest.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StorageReadResponse(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageReadResponse)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageReadResponse');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageReadResponse(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageReadResponse.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StorageWriteRequest(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageWriteRequest)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageWriteRequest');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageWriteRequest(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageWriteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_storage_v2_StorageWriteResponse(arg) {
    if (!(arg instanceof proto_storage_v2_storage_pb.StorageWriteResponse)) {
        throw new Error('Expected argument of type nitric.proto.storage.v2.StorageWriteResponse');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_storage_v2_StorageWriteResponse(buffer_arg) {
    return proto_storage_v2_storage_pb.StorageWriteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}
// Services for storage and retrieval of blobs in the form of byte arrays, such as text and binary blobs.
var StorageService = exports.StorageService = {
    // Retrieve an item from a bucket
    read: {
        path: '/nitric.proto.storage.v2.Storage/Read',
        requestStream: false,
        responseStream: false,
        requestType: proto_storage_v2_storage_pb.StorageReadRequest,
        responseType: proto_storage_v2_storage_pb.StorageReadResponse,
        requestSerialize: serialize_nitric_proto_storage_v2_StorageReadRequest,
        requestDeserialize: deserialize_nitric_proto_storage_v2_StorageReadRequest,
        responseSerialize: serialize_nitric_proto_storage_v2_StorageReadResponse,
        responseDeserialize: deserialize_nitric_proto_storage_v2_StorageReadResponse,
    },
    // Store an item to a bucket
    write: {
        path: '/nitric.proto.storage.v2.Storage/Write',
        requestStream: false,
        responseStream: false,
        requestType: proto_storage_v2_storage_pb.StorageWriteRequest,
        responseType: proto_storage_v2_storage_pb.StorageWriteResponse,
        requestSerialize: serialize_nitric_proto_storage_v2_StorageWriteRequest,
        requestDeserialize: deserialize_nitric_proto_storage_v2_StorageWriteRequest,
        responseSerialize: serialize_nitric_proto_storage_v2_StorageWriteResponse,
        responseDeserialize: deserialize_nitric_proto_storage_v2_StorageWriteResponse,
    },
    // Delete an item from a bucket
    delete: {
        path: '/nitric.proto.storage.v2.Storage/Delete',
        requestStream: false,
        responseStream: false,
        requestType: proto_storage_v2_storage_pb.StorageDeleteRequest,
        responseType: proto_storage_v2_storage_pb.StorageDeleteResponse,
        requestSerialize: serialize_nitric_proto_storage_v2_StorageDeleteRequest,
        requestDeserialize: deserialize_nitric_proto_storage_v2_StorageDeleteRequest,
        responseSerialize: serialize_nitric_proto_storage_v2_StorageDeleteResponse,
        responseDeserialize: deserialize_nitric_proto_storage_v2_StorageDeleteResponse,
    },
    // Generate a pre-signed URL for direct operations on an item
    preSignUrl: {
        path: '/nitric.proto.storage.v2.Storage/PreSignUrl',
        requestStream: false,
        responseStream: false,
        requestType: proto_storage_v2_storage_pb.StoragePreSignUrlRequest,
        responseType: proto_storage_v2_storage_pb.StoragePreSignUrlResponse,
        requestSerialize: serialize_nitric_proto_storage_v2_StoragePreSignUrlRequest,
        requestDeserialize: deserialize_nitric_proto_storage_v2_StoragePreSignUrlRequest,
        responseSerialize: serialize_nitric_proto_storage_v2_StoragePreSignUrlResponse,
        responseDeserialize: deserialize_nitric_proto_storage_v2_StoragePreSignUrlResponse,
    },
    // List blobs currently in the bucket
    listBlobs: {
        path: '/nitric.proto.storage.v2.Storage/ListBlobs',
        requestStream: false,
        responseStream: false,
        requestType: proto_storage_v2_storage_pb.StorageListBlobsRequest,
        responseType: proto_storage_v2_storage_pb.StorageListBlobsResponse,
        requestSerialize: serialize_nitric_proto_storage_v2_StorageListBlobsRequest,
        requestDeserialize: deserialize_nitric_proto_storage_v2_StorageListBlobsRequest,
        responseSerialize: serialize_nitric_proto_storage_v2_StorageListBlobsResponse,
        responseDeserialize: deserialize_nitric_proto_storage_v2_StorageListBlobsResponse,
    },
    // Determine is an object exists in a bucket
    exists: {
        path: '/nitric.proto.storage.v2.Storage/Exists',
        requestStream: false,
        responseStream: false,
        requestType: proto_storage_v2_storage_pb.StorageExistsRequest,
        responseType: proto_storage_v2_storage_pb.StorageExistsResponse,
        requestSerialize: serialize_nitric_proto_storage_v2_StorageExistsRequest,
        requestDeserialize: deserialize_nitric_proto_storage_v2_StorageExistsRequest,
        responseSerialize: serialize_nitric_proto_storage_v2_StorageExistsResponse,
        responseDeserialize: deserialize_nitric_proto_storage_v2_StorageExistsResponse,
    },
};
exports.StorageClient = grpc.makeGenericClientConstructor(StorageService, 'Storage');
