export namespace PubsubService {
    namespace publish {
        export let path: string;
        export let requestStream: boolean;
        export let responseStream: boolean;
        export let requestType: typeof proto_pubsub_v2_pubsub_pb.PubsubPublishRequest;
        export let responseType: typeof proto_pubsub_v2_pubsub_pb.PubsubPublishResponse;
        export { serialize_nitric_proto_pubsub_v1_PubsubPublishRequest as requestSerialize };
        export { deserialize_nitric_proto_pubsub_v1_PubsubPublishRequest as requestDeserialize };
        export { serialize_nitric_proto_pubsub_v1_PubsubPublishResponse as responseSerialize };
        export { deserialize_nitric_proto_pubsub_v1_PubsubPublishResponse as responseDeserialize };
    }
}
export const PubsubClient: grpc.ServiceClientConstructor;
import proto_pubsub_v2_pubsub_pb = require("../../../proto/pubsub/v2/pubsub_pb.js");
declare function serialize_nitric_proto_pubsub_v1_PubsubPublishRequest(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_pubsub_v1_PubsubPublishRequest(buffer_arg: any): proto_pubsub_v2_pubsub_pb.PubsubPublishRequest;
declare function serialize_nitric_proto_pubsub_v1_PubsubPublishResponse(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_pubsub_v1_PubsubPublishResponse(buffer_arg: any): proto_pubsub_v2_pubsub_pb.PubsubPublishResponse;
import grpc = require("@grpc/grpc-js");
export {};
