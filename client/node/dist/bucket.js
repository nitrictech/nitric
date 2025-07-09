"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Bucket = exports.Mode = void 0;
const storage_pb_1 = require("./gen/proto/storage/v2/storage_pb");
const duration_pb_1 = require("google-protobuf/google/protobuf/duration_pb");
var Mode;
(function (Mode) {
    Mode["Read"] = "READ";
    Mode["Write"] = "WRITE";
})(Mode || (exports.Mode = Mode = {}));
class Bucket {
    constructor(storageClient, bucketName) {
        this.name = bucketName;
        this.storageClient = storageClient;
    }
    /**
     * Read a file from the bucket
     */
    async read(key) {
        const request = new storage_pb_1.StorageReadRequest();
        request.setBucketName(this.name);
        request.setKey(key);
        return new Promise((resolve, reject) => {
            this.storageClient.read(request, (error, response) => {
                if (error) {
                    reject(new Error(`Failed to read file from the ${this.name} bucket: ${error.message}`));
                    return;
                }
                if (!response) {
                    reject(new Error(`Failed to read file from the ${this.name} bucket: No response received`));
                    return;
                }
                resolve(Buffer.from(response.getBody_asU8()));
            });
        });
    }
    /**
     * Write a file to the bucket
     */
    async write(key, data) {
        const request = new storage_pb_1.StorageWriteRequest();
        request.setBucketName(this.name);
        request.setKey(key);
        request.setBody(data);
        return new Promise((resolve, reject) => {
            this.storageClient.write(request, (error) => {
                if (error) {
                    reject(new Error(`Failed to write file to bucket: ${error.message}`));
                    return;
                }
                resolve();
            });
        });
    }
    /**
     * Delete a file from the bucket
     */
    async delete(key) {
        const request = new storage_pb_1.StorageDeleteRequest();
        request.setBucketName(this.name);
        request.setKey(key);
        return new Promise((resolve, reject) => {
            this.storageClient.delete(request, (error) => {
                if (error) {
                    reject(new Error(`Failed to delete file from bucket: ${error.message}`));
                    return;
                }
                resolve();
            });
        });
    }
    /**
     * List files in the bucket with a given prefix
     */
    async list(prefix) {
        const request = new storage_pb_1.StorageListBlobsRequest();
        request.setBucketName(this.name);
        request.setPrefix(prefix);
        return new Promise((resolve, reject) => {
            this.storageClient.listBlobs(request, (error, response) => {
                if (error) {
                    reject(new Error(`Failed to list files in bucket: ${error.message}`));
                    return;
                }
                if (!response) {
                    reject(new Error(`Failed to list files in bucket: No response received`));
                    return;
                }
                const blobs = response.getBlobsList();
                resolve(blobs.map((blob) => blob.getKey()));
            });
        });
    }
    /**
     * Check if a file exists in the bucket
     */
    async exists(key) {
        const request = new storage_pb_1.StorageExistsRequest();
        request.setBucketName(this.name);
        request.setKey(key);
        return new Promise((resolve, reject) => {
            this.storageClient.exists(request, (error, response) => {
                if (error) {
                    reject(new Error(`Failed to check if file exists in bucket: ${error.message}`));
                    return;
                }
                if (!response) {
                    reject(new Error(`Failed to check if file exists in bucket: No response received`));
                    return;
                }
                resolve(response.getExists());
            });
        });
    }
    /**
     * Get a presigned URL for downloading a file from the bucket
     */
    async getDownloadUrl(key, options) {
        const defaultOptions = {
            mode: Mode.Read,
            expiry: 300, // 5 minutes in seconds
        };
        const opts = { ...defaultOptions, ...options };
        return this.preSignUrl(key, opts);
    }
    /**
     * Get a presigned URL for uploading a file to the bucket
     */
    async getUploadUrl(key, options) {
        const defaultOptions = {
            mode: Mode.Write,
            expiry: 300, // 5 minutes in seconds
        };
        const opts = { ...defaultOptions, ...options };
        return this.preSignUrl(key, opts);
    }
    async preSignUrl(key, options) {
        const request = new storage_pb_1.StoragePreSignUrlRequest();
        request.setBucketName(this.name);
        request.setKey(key);
        request.setOperation(options.mode === Mode.Read
            ? storage_pb_1.StoragePreSignUrlRequest.Operation.READ
            : storage_pb_1.StoragePreSignUrlRequest.Operation.WRITE);
        const duration = new duration_pb_1.Duration();
        duration.setSeconds(options.expiry);
        request.setExpiry(duration);
        return new Promise((resolve, reject) => {
            this.storageClient.preSignUrl(request, (error, response) => {
                if (error) {
                    reject(new Error(`Failed to get presigned URL for file: ${error.message}`));
                    return;
                }
                if (!response) {
                    reject(new Error(`Failed to get presigned URL for file: No response received`));
                    return;
                }
                resolve(response.getUrl());
            });
        });
    }
}
exports.Bucket = Bucket;
