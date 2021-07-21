# Nitric Membrane Plugins

The Nitric Plugins represent pluggable services that back the unified nitric service APIs

## Providers

Providers are a collection of plugin factories, these are compiled into [Go plugins](https://golang.org/pkg/plugin/) to be loaded by the pluggable memrane.

Example providers:
 * [AWS](../providers/aws/README.md)
 * [GCP](../providers/gcp/README.md)

## Gateway Plugin

Gateway plugins are used to abstract implementation detail for application input/output. This usually involves presenting some kind of external interface (for example a HTTP server), normalizing incoming data from this interface to pass to the hosted application (over HTTP) and then denormalizing output to be returned back to the hosted environment. 

Official Gateway Plugins:

* [AWS Lambda](./gateway/lambda/README.md)
* [AWS HTTP](./gateway/ecs/README.md) (typically used with ECS)
* [GCP HTTP](./gateway/cloudrun/README.md) (typically used with Cloud Run)
* [Dev HTTP](./gateway/dev/README.md) (used for local development)

## Eventing Plugin

Eventing plugins are used for communication between services, it exposes a simple publish model that allows [NitricEvent(s)]() to be pushed onto topics. Topic subscriptions are normally defined at deploy time and configured through the [Nitric Stack]() definition file.

Official Eventing Plugins:
* [AWS SNS](./eventing/sns/README.md)
* [GCP Pubsub](./eventing/pubsub/README.md)
* [Dev Eventing](./eventing/dev/README.md)

## Document Plugin

Documents plugins are used to provide a queryable key value store interface that allows users to store simple object data under a given collection and key.

Official Document Plugins:
* [AWS Dynamodb](./document/dynamodb/README.md)
* [GCP Firestore](./document/firestore/README.md)
* [Dev Documents](./document/boltdb/README.md)

## Queue Plugin

Queue plugins provide a simple push/pop interface allowing users to asynchronously process batch operations.

Official Queue Plugins:
* [AWS SQS](./queue/sqs/README.md)
* [GCP Pubsub](./queue/pubsub/README.md)
* [Dev Queue](./queue/dev/README.md)

## Storage Plugins

Storage plugins provide access to provide blob stores, allowing storage of files using a simple bucket/key interface for storage and retrieval.

Official Storage Plugins:
* [AWS S3](./storage/s3/README.md)
* [GCP Cloud Storage](./storage/storage/README.md)
* [Dev Storage](./storage/dev/README.md)



