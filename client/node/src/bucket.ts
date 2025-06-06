import { StorageClient } from "./gen/proto/storage/v2/storage_grpc_pb";
import {
  StorageReadRequest,
  StorageWriteRequest,
  StorageDeleteRequest,
  StorageListBlobsRequest,
  StorageExistsRequest,
  StoragePreSignUrlRequest,
} from "./gen/proto/storage/v2/storage_pb";
import { Duration } from "google-protobuf/google/protobuf/duration_pb";

export enum Mode {
  Read = "READ",
  Write = "WRITE",
}

export interface PresignUrlOptions {
  mode: Mode;
  expiry: number; // Duration in seconds
}

export class Bucket {
  private name: string;
  private storageClient: StorageClient;

  constructor(storageClient: StorageClient, bucketName: string) {
    this.name = bucketName;
    this.storageClient = storageClient;
  }

  /**
   * Read a file from the bucket
   */
  async read(key: string): Promise<Buffer> {
    const request = new StorageReadRequest();
    request.setBucketName(this.name);
    request.setKey(key);

    return new Promise((resolve, reject) => {
      this.storageClient.read(request, (error, response) => {
        if (error) {
          reject(
            new Error(
              `Failed to read file from the ${this.name} bucket: ${error.message}`
            )
          );
          return;
        }
        if (!response) {
          reject(
            new Error(
              `Failed to read file from the ${this.name} bucket: No response received`
            )
          );
          return;
        }
        resolve(Buffer.from(response.getBody_asU8()));
      });
    });
  }

  /**
   * Write a file to the bucket
   */
  async write(key: string, data: Buffer): Promise<void> {
    const request = new StorageWriteRequest();
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
  async delete(key: string): Promise<void> {
    const request = new StorageDeleteRequest();
    request.setBucketName(this.name);
    request.setKey(key);

    return new Promise((resolve, reject) => {
      this.storageClient.delete(request, (error) => {
        if (error) {
          reject(
            new Error(`Failed to delete file from bucket: ${error.message}`)
          );
          return;
        }
        resolve();
      });
    });
  }

  /**
   * List files in the bucket with a given prefix
   */
  async list(prefix: string): Promise<string[]> {
    const request = new StorageListBlobsRequest();
    request.setBucketName(this.name);
    request.setPrefix(prefix);

    return new Promise((resolve, reject) => {
      this.storageClient.listBlobs(request, (error, response) => {
        if (error) {
          reject(new Error(`Failed to list files in bucket: ${error.message}`));
          return;
        }
        if (!response) {
          reject(
            new Error(`Failed to list files in bucket: No response received`)
          );
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
  async exists(key: string): Promise<boolean> {
    const request = new StorageExistsRequest();
    request.setBucketName(this.name);
    request.setKey(key);

    return new Promise((resolve, reject) => {
      this.storageClient.exists(request, (error, response) => {
        if (error) {
          reject(
            new Error(
              `Failed to check if file exists in bucket: ${error.message}`
            )
          );
          return;
        }
        if (!response) {
          reject(
            new Error(
              `Failed to check if file exists in bucket: No response received`
            )
          );
          return;
        }
        resolve(response.getExists());
      });
    });
  }

  /**
   * Get a presigned URL for downloading a file from the bucket
   */
  async getDownloadUrl(
    key: string,
    options?: Partial<PresignUrlOptions>
  ): Promise<string> {
    const defaultOptions: PresignUrlOptions = {
      mode: Mode.Read,
      expiry: 300, // 5 minutes in seconds
    };

    const opts = { ...defaultOptions, ...options };
    return this.preSignUrl(key, opts);
  }

  /**
   * Get a presigned URL for uploading a file to the bucket
   */
  async getUploadUrl(
    key: string,
    options?: Partial<PresignUrlOptions>
  ): Promise<string> {
    const defaultOptions: PresignUrlOptions = {
      mode: Mode.Write,
      expiry: 300, // 5 minutes in seconds
    };

    const opts = { ...defaultOptions, ...options };
    return this.preSignUrl(key, opts);
  }

  private async preSignUrl(
    key: string,
    options: PresignUrlOptions
  ): Promise<string> {
    const request = new StoragePreSignUrlRequest();
    request.setBucketName(this.name);
    request.setKey(key);
    request.setOperation(
      options.mode === Mode.Read
        ? StoragePreSignUrlRequest.Operation.READ
        : StoragePreSignUrlRequest.Operation.WRITE
    );

    const duration = new Duration();
    duration.setSeconds(options.expiry);
    request.setExpiry(duration);

    return new Promise((resolve, reject) => {
      this.storageClient.preSignUrl(request, (error, response) => {
        if (error) {
          reject(
            new Error(`Failed to get presigned URL for file: ${error.message}`)
          );
          return;
        }
        if (!response) {
          reject(
            new Error(
              `Failed to get presigned URL for file: No response received`
            )
          );
          return;
        }
        resolve(response.getUrl());
      });
    });
  }
}
