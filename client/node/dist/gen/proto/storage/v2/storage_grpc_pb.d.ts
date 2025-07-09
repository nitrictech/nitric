export namespace StorageService {
    export namespace read {
        export let path: string;
        export let requestStream: boolean;
        export let responseStream: boolean;
        export let requestType: typeof proto_storage_v2_storage_pb.StorageReadRequest;
        export let responseType: typeof proto_storage_v2_storage_pb.StorageReadResponse;
        export { serialize_nitric_proto_storage_v2_StorageReadRequest as requestSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageReadRequest as requestDeserialize };
        export { serialize_nitric_proto_storage_v2_StorageReadResponse as responseSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageReadResponse as responseDeserialize };
    }
    export namespace write {
        let path_1: string;
        export { path_1 as path };
        let requestStream_1: boolean;
        export { requestStream_1 as requestStream };
        let responseStream_1: boolean;
        export { responseStream_1 as responseStream };
        let requestType_1: typeof proto_storage_v2_storage_pb.StorageWriteRequest;
        export { requestType_1 as requestType };
        let responseType_1: typeof proto_storage_v2_storage_pb.StorageWriteResponse;
        export { responseType_1 as responseType };
        export { serialize_nitric_proto_storage_v2_StorageWriteRequest as requestSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageWriteRequest as requestDeserialize };
        export { serialize_nitric_proto_storage_v2_StorageWriteResponse as responseSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageWriteResponse as responseDeserialize };
    }
    export namespace _delete {
        let path_2: string;
        export { path_2 as path };
        let requestStream_2: boolean;
        export { requestStream_2 as requestStream };
        let responseStream_2: boolean;
        export { responseStream_2 as responseStream };
        let requestType_2: typeof proto_storage_v2_storage_pb.StorageDeleteRequest;
        export { requestType_2 as requestType };
        let responseType_2: typeof proto_storage_v2_storage_pb.StorageDeleteResponse;
        export { responseType_2 as responseType };
        export { serialize_nitric_proto_storage_v2_StorageDeleteRequest as requestSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageDeleteRequest as requestDeserialize };
        export { serialize_nitric_proto_storage_v2_StorageDeleteResponse as responseSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageDeleteResponse as responseDeserialize };
    }
    export { _delete as delete };
    export namespace preSignUrl {
        let path_3: string;
        export { path_3 as path };
        let requestStream_3: boolean;
        export { requestStream_3 as requestStream };
        let responseStream_3: boolean;
        export { responseStream_3 as responseStream };
        let requestType_3: typeof proto_storage_v2_storage_pb.StoragePreSignUrlRequest;
        export { requestType_3 as requestType };
        let responseType_3: typeof proto_storage_v2_storage_pb.StoragePreSignUrlResponse;
        export { responseType_3 as responseType };
        export { serialize_nitric_proto_storage_v2_StoragePreSignUrlRequest as requestSerialize };
        export { deserialize_nitric_proto_storage_v2_StoragePreSignUrlRequest as requestDeserialize };
        export { serialize_nitric_proto_storage_v2_StoragePreSignUrlResponse as responseSerialize };
        export { deserialize_nitric_proto_storage_v2_StoragePreSignUrlResponse as responseDeserialize };
    }
    export namespace listBlobs {
        let path_4: string;
        export { path_4 as path };
        let requestStream_4: boolean;
        export { requestStream_4 as requestStream };
        let responseStream_4: boolean;
        export { responseStream_4 as responseStream };
        let requestType_4: typeof proto_storage_v2_storage_pb.StorageListBlobsRequest;
        export { requestType_4 as requestType };
        let responseType_4: typeof proto_storage_v2_storage_pb.StorageListBlobsResponse;
        export { responseType_4 as responseType };
        export { serialize_nitric_proto_storage_v2_StorageListBlobsRequest as requestSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageListBlobsRequest as requestDeserialize };
        export { serialize_nitric_proto_storage_v2_StorageListBlobsResponse as responseSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageListBlobsResponse as responseDeserialize };
    }
    export namespace exists {
        let path_5: string;
        export { path_5 as path };
        let requestStream_5: boolean;
        export { requestStream_5 as requestStream };
        let responseStream_5: boolean;
        export { responseStream_5 as responseStream };
        let requestType_5: typeof proto_storage_v2_storage_pb.StorageExistsRequest;
        export { requestType_5 as requestType };
        let responseType_5: typeof proto_storage_v2_storage_pb.StorageExistsResponse;
        export { responseType_5 as responseType };
        export { serialize_nitric_proto_storage_v2_StorageExistsRequest as requestSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageExistsRequest as requestDeserialize };
        export { serialize_nitric_proto_storage_v2_StorageExistsResponse as responseSerialize };
        export { deserialize_nitric_proto_storage_v2_StorageExistsResponse as responseDeserialize };
    }
}
export const StorageClient: grpc.ServiceClientConstructor;
import proto_storage_v2_storage_pb = require("../../../proto/storage/v2/storage_pb.js");
declare function serialize_nitric_proto_storage_v2_StorageReadRequest(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageReadRequest(buffer_arg: any): proto_storage_v2_storage_pb.StorageReadRequest;
declare function serialize_nitric_proto_storage_v2_StorageReadResponse(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageReadResponse(buffer_arg: any): proto_storage_v2_storage_pb.StorageReadResponse;
declare function serialize_nitric_proto_storage_v2_StorageWriteRequest(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageWriteRequest(buffer_arg: any): proto_storage_v2_storage_pb.StorageWriteRequest;
declare function serialize_nitric_proto_storage_v2_StorageWriteResponse(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageWriteResponse(buffer_arg: any): proto_storage_v2_storage_pb.StorageWriteResponse;
declare function serialize_nitric_proto_storage_v2_StorageDeleteRequest(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageDeleteRequest(buffer_arg: any): proto_storage_v2_storage_pb.StorageDeleteRequest;
declare function serialize_nitric_proto_storage_v2_StorageDeleteResponse(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageDeleteResponse(buffer_arg: any): proto_storage_v2_storage_pb.StorageDeleteResponse;
declare function serialize_nitric_proto_storage_v2_StoragePreSignUrlRequest(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StoragePreSignUrlRequest(buffer_arg: any): proto_storage_v2_storage_pb.StoragePreSignUrlRequest;
declare function serialize_nitric_proto_storage_v2_StoragePreSignUrlResponse(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StoragePreSignUrlResponse(buffer_arg: any): proto_storage_v2_storage_pb.StoragePreSignUrlResponse;
declare function serialize_nitric_proto_storage_v2_StorageListBlobsRequest(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageListBlobsRequest(buffer_arg: any): proto_storage_v2_storage_pb.StorageListBlobsRequest;
declare function serialize_nitric_proto_storage_v2_StorageListBlobsResponse(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageListBlobsResponse(buffer_arg: any): proto_storage_v2_storage_pb.StorageListBlobsResponse;
declare function serialize_nitric_proto_storage_v2_StorageExistsRequest(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageExistsRequest(buffer_arg: any): proto_storage_v2_storage_pb.StorageExistsRequest;
declare function serialize_nitric_proto_storage_v2_StorageExistsResponse(arg: any): Buffer<ArrayBuffer>;
declare function deserialize_nitric_proto_storage_v2_StorageExistsResponse(buffer_arg: any): proto_storage_v2_storage_pb.StorageExistsResponse;
import grpc = require("@grpc/grpc-js");
export {};
