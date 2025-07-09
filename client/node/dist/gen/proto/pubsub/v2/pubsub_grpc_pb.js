// GENERATED CODE -- DO NOT EDIT!
'use strict';
var grpc = require('@grpc/grpc-js');
var proto_pubsub_v2_pubsub_pb = require('../../../proto/pubsub/v2/pubsub_pb.js');
function serialize_nitric_proto_pubsub_v1_PubsubPublishRequest(arg) {
    if (!(arg instanceof proto_pubsub_v2_pubsub_pb.PubsubPublishRequest)) {
        throw new Error('Expected argument of type nitric.proto.pubsub.v1.PubsubPublishRequest');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_pubsub_v1_PubsubPublishRequest(buffer_arg) {
    return proto_pubsub_v2_pubsub_pb.PubsubPublishRequest.deserializeBinary(new Uint8Array(buffer_arg));
}
function serialize_nitric_proto_pubsub_v1_PubsubPublishResponse(arg) {
    if (!(arg instanceof proto_pubsub_v2_pubsub_pb.PubsubPublishResponse)) {
        throw new Error('Expected argument of type nitric.proto.pubsub.v1.PubsubPublishResponse');
    }
    return Buffer.from(arg.serializeBinary());
}
function deserialize_nitric_proto_pubsub_v1_PubsubPublishResponse(buffer_arg) {
    return proto_pubsub_v2_pubsub_pb.PubsubPublishResponse.deserializeBinary(new Uint8Array(buffer_arg));
}
// Service for publishing asynchronous messages
var PubsubService = exports.PubsubService = {
    // Publishes a message to a given topic
    publish: {
        path: '/nitric.proto.pubsub.v1.Pubsub/Publish',
        requestStream: false,
        responseStream: false,
        requestType: proto_pubsub_v2_pubsub_pb.PubsubPublishRequest,
        responseType: proto_pubsub_v2_pubsub_pb.PubsubPublishResponse,
        requestSerialize: serialize_nitric_proto_pubsub_v1_PubsubPublishRequest,
        requestDeserialize: deserialize_nitric_proto_pubsub_v1_PubsubPublishRequest,
        responseSerialize: serialize_nitric_proto_pubsub_v1_PubsubPublishResponse,
        responseDeserialize: deserialize_nitric_proto_pubsub_v1_PubsubPublishResponse,
    },
};
exports.PubsubClient = grpc.makeGenericClientConstructor(PubsubService, 'Pubsub');
