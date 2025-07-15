import { StorageClient } from "./gen/proto/storage/v2/storage_grpc_pb";
export declare enum Mode {
    Read = "READ",
    Write = "WRITE"
}
export interface PresignUrlOptions {
    mode: Mode;
    expiry: number;
}
export declare class Bucket {
    private name;
    private storageClient;
    constructor(storageClient: StorageClient, bucketName: string);
    /**
     * Read a file from the bucket
     */
    read(key: string): Promise<Buffer>;
    /**
     * Write a file to the bucket
     */
    write(key: string, data: Buffer): Promise<void>;
    /**
     * Delete a file from the bucket
     */
    delete(key: string): Promise<void>;
    /**
     * List files in the bucket with a given prefix
     */
    list(prefix: string): Promise<string[]>;
    /**
     * Check if a file exists in the bucket
     */
    exists(key: string): Promise<boolean>;
    /**
     * Get a presigned URL for downloading a file from the bucket
     */
    getDownloadUrl(key: string, options?: Partial<PresignUrlOptions>): Promise<string>;
    /**
     * Get a presigned URL for uploading a file to the bucket
     */
    getUploadUrl(key: string, options?: Partial<PresignUrlOptions>): Promise<string>;
    private preSignUrl;
}
