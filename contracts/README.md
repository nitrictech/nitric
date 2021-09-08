contracts

## Proto Contracts

### Naming Conventions

#### Request/Response Messages

All request and Response message type names are derived from the service and rpc (method) they relate to. The format is:

`Service + Method + Direction`

> __Example:__ In the `User` Service, for the `Create` rpc the message type names are `UserCreateRequest` and `UserCreateResponse`.

All outbound messages have the `Request` suffix, while inbound messages have `Response`.
