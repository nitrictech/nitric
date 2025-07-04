from typing import Optional, List
from sdk.gen.storage.v2 import storage_pb2, storage_pb2_grpc
from google.protobuf.duration_pb2 import Duration


class Mode:
    READ = "READ"
    WRITE = "WRITE"


class PresignUrlOptions:
    def __init__(self, mode: str = Mode.READ, expiry: int = 300):
        self.mode = mode
        self.expiry = expiry


class Bucket:
    def __init__(self, storage_client: storage_pb2_grpc.StorageStub, bucket_name: str):
        self.name = bucket_name
        self.storage_client = storage_client

    async def read(self, key: str) -> bytes:
        request = storage_pb2.StorageReadRequest(bucket_name=self.name, key=key)
        response = await self.storage_client.Read(request)
        return response.body

    async def write(self, key: str, data: bytes) -> None:
        request = storage_pb2.StorageWriteRequest(
            bucket_name=self.name,
            key=key,
            body=data,
        )
        await self.storage_client.Write(request)

    async def delete(self, key: str) -> None:
        request = storage_pb2.StorageDeleteRequest(bucket_name=self.name, key=key)
        await self.storage_client.Delete(request)

    async def list(self, prefix: str = "") -> List[str]:
        request = storage_pb2.StorageListBlobsRequest(bucket_name=self.name, prefix=prefix)
        response = await self.storage_client.ListBlobs(request)
        return [blob.key for blob in response.blobs]

    async def exists(self, key: str) -> bool:
        request = storage_pb2.StorageExistsRequest(bucket_name=self.name, key=key)
        response = await self.storage_client.Exists(request)
        return response.exists

    async def get_download_url(self, key: str, options: Optional[PresignUrlOptions] = None) -> str:
        opts = options or PresignUrlOptions(mode=Mode.READ, expiry=300)
        return await self._pre_sign_url(key, opts)

    async def get_upload_url(self, key: str, options: Optional[PresignUrlOptions] = None) -> str:
        opts = options or PresignUrlOptions(mode=Mode.WRITE, expiry=300)
        return await self._pre_sign_url(key, opts)

    async def _pre_sign_url(self, key: str, options: PresignUrlOptions) -> str:
        operation = (
            storage_pb2.StoragePreSignUrlRequest.READ
            if options.mode == Mode.READ
            else storage_pb2.StoragePreSignUrlRequest.WRITE
        )

        request = storage_pb2.StoragePreSignUrlRequest(
            bucket_name=self.name,
            key=key,
            operation=operation,
            expiry=Duration(seconds=options.expiry),
        )

        response = await self.storage_client.PreSignUrl(request)
        return response.url
