# Nitric Membrane SDK

The Nitric SDK is a collection of golang interfaces that are used to abstract services provided by Nitric.

## Providers

Providers are a collection of plugin factories, these are compiled into [Go plugins](https://golang.org/pkg/plugin/) to be loaded by the pluggable memrane.

Example providers:
 * [AWS](../aws/README.md)
 * [GCP](../gcp/README.md)

## Gateway Plugin

Gateway plugins are used to abstract implementation detail for application input/output. This usually involves presenting some kind of external interface (for example a HTTP server), normalizing incoming data from this interface to pass to the hosted application (over HTTP) and then denormalizing output to be returned back to the hosted environment. 

Official Gateway Plugins:

* [AWS Lambda](../aws/gateway/lambda/README.md)
* [AWS HTTP](../aws/gateway/http/README.md) (typically used with ECS)
* [GCP HTTP](../gcp/gateway/http/README.md) (typically used with Cloud Run)
* [Dev HTTP](../dev/gateway/README.md) (used for local development)

## Eventing Plugin

Eventing plugins are used for communication between services, it exposes a simple publish model that allows [NitricEvent(s)]() to be pushed onto topics. Topic subscriptions are normally defined at deploy time and configured through the [Nitric Stack]() definition file.

Official Eventing Plugins:
* [AWS SNS](../aws/eventing/sns/README.md)
* [GCP Pubsub](../gcp/eventing/pubsub/README.md)
* [Dev Eventing](../dev/eventing/README.md)

## Documents Plugin

Documents plugins are used to provide a simple document store interface, that allows users to store simple object data under a given collection and key value.

Official Documents Plugins:
* [AWS Dynamodb](../aws/documents/dynamodb/README.md)
* [GCP Firestore](../gcp/documents/firestore/README.md)
* [Dev Documents](../dev/documents/README.md)

## Queue Plugin

Queue plugins provide a simple push/pop interface allowing users to asynchronously process batch operations.

Official Queue Plugins:
* [AWS SQS](../aws/queue/sqs/README.md)
* [GCP Pubsub](../gcp/queue/pubsub/README.md)
* [Dev Queue](../dev/queue/README.md)

## Storage Plugin

Storage plugins provide access to provide blob stores, allowing storage of files using a simple bucket/key interface for storage and retrieval.

Official Storage Plugins:
* [AWS S3](../aws/storage/s3/README.md)
* [GCP Cloud Storage](../gcp/storage/storage/README.md)
* [Dev Storage](../dev/storage/README.md)


## Auth Plugin

Auth plugins allow for the creation and management of users and tenants as well as login/verification of users and their tokens.

Official Auth Plugins:
* [AWS Cognito](../aws/auth/cognito/README.md)
* [GCP Identity Platform](../gcp/auth/identityplatform/README.md)
* [Dev Auth](../dev/auth/README.md)




