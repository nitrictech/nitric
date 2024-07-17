<p align="center">
  <a href="https://nitric.io">
    <img src="docs/assets/nitric-logo.svg" width="120" alt="Nitric Logo"/>
  </a>
</p>

<h3 align="center">build cloud aware applications</h3>

<p align="center">
  <img alt="GitHub release (latest SemVer)" src="https://img.shields.io/github/v/release/nitrictech/nitric?style=for-the-badge">
  <img alt="GitHub" src="https://img.shields.io/github/license/nitrictech/nitric?style=for-the-badge">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/nitrictech/nitric/test.yaml?branch=develop&style=for-the-badge">
  <img alt="Codecov" src="https://img.shields.io/codecov/c/github/nitrictech/nitric?style=for-the-badge">
  <a href="https://discord.gg/Webemece5C"><img alt="Discord" src="https://img.shields.io/discord/955259353043173427?label=discord&style=for-the-badge"></a>
</p>

## About Nitric

[Nitric](https://nitric.io) is a multi-language framework, with concise inline infrastructure from code. Modern applications should be robust, productive and a joy to build. Nitric solves common problems building for modern platforms:

- [Easy infrastructure](https://nitric.io/docs/concepts/introduction#infrastructure-from-code-if-c) from code
- Build for [any host without coupling](https://nitric.io/docs/reference/providers)
- [Run locally](https://nitric.io/docs/getting-started/local-dashboard)
- [IAM for humans](https://nitric.io/docs/concepts/access-control)
- Common resources like [databases](https://nitric.io/docs/sql), [queues/topics](https://nitric.io/docs/messaging), [APIs](https://nitric.io/docs/apis), [key-value](https://nitric.io/docs/keyvalue), [buckets](https://nitric.io/docs/storage) and more
- [Change services, IaC tools or cloud providers](https://nitric.io/docs/reference/providers) without changing code

We also knows abstraction should mean building on existing layers, not hiding them. Nitric includes powerful escape hatches for when things get custom.

## Supported Languages

<p>
  <a href="https://github.com/nitrictech/node-sdk"><img src="https://skillicons.dev/icons?i=js"/></a>
  <a href="https://github.com/nitrictech/node-sdk"><img src="https://skillicons.dev/icons?i=ts"/></a>
  <a href="https://github.com/nitrictech/python-sdk"><img src="https://skillicons.dev/icons?i=py"/></a>
  <a href="https://github.com/nitrictech/go-sdk"><img src="https://skillicons.dev/icons?i=go"/></a>
  <a href="https://github.com/nitrictech/dotnet-sdk"><img src="https://skillicons.dev/icons?i=cs"/></a>
  <a href="https://github.com/nitrictech/jvm-sdk"><img src="https://skillicons.dev/icons?i=java"/></a>
  <a href="https://github.com/nitrictech/dart-sdk"><img src="https://skillicons.dev/icons?i=dart"/></a>
</p>

## Supported Clouds

<p>
  <a href="./cloud/aws"><img src="https://skillicons.dev/icons?i=aws"/></a>
  <a href="./cloud/gcp"><img src="https://skillicons.dev/icons?i=gcp"/></a>
  <a href="./cloud/azure"><img src="https://skillicons.dev/icons?i=azure"/></a>
</p>

> These are supported out of the box, but you can also build [custom providers](https://nitric.io/docs/reference/providers/custom/building-custom-provider) as well

## Example

Creating an API, a bucket with access permissions and writing files to that bucket via a serverless function.

```javascript
// JavaScript Example
import { api, bucket } from "@nitric/sdk";

const main = api("main");
const notes = bucket("notes").allow("read", "write");

main.post("/notes/:title", async (ctx) => {
  const { title } = ctx.req.params;
  await notes.file(title).write(ctx.req.text());
});
```

This is the only code needed to deploy a working application to any cloud provider you choose, using [`nitric up`](https://nitric.io/docs/getting-started/deployment). Nitric can deploy this application using automatically generated [Pulumi](https://nitric.io/docs/reference/providers/pulumi), [Terraform](https://nitric.io/docs/reference/providers/terraform) or [any other automation tools](https://nitric.io/docs/reference/providers/custom/building-custom-provider) you choose.

## Documentation

Nitric has full documentation at [nitric.io/docs](https://nitric.io/docs), including concepts, reference documentation for various languages and many tutorials/guides.

## Get in touch

- Ask questions in [GitHub discussions](https://github.com/nitrictech/nitric/discussions)

- Join us on [Discord](https://nitric.io/chat)

- Find us on [Twitter](https://twitter.com/nitric_io)

- Send us an [email](mailto:maintainers@nitric.io)

## Contributing

We greatly appreciate contributions, consider starting with the [contributions guide](./CONTRIBUTING.md) and a chat on [Discord](https://nitric.io/chat) or [GitHub](https://github.com/nitrictech/nitric/discussions).
