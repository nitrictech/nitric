<p align="center">
  <a href="https://nitric.io">
    <img src="assets/nitric-logo.svg" width="120" alt="Nitric Logo"/>
  </a>
</p>

<h3 align="center">Effortless backends with infrastructure from code</h3>

<p align="center">
  <img alt="GitHub release (latest SemVer)" src="https://img.shields.io/github/v/release/nitrictech/nitric?style=for-the-badge">
  <img alt="GitHub" src="https://img.shields.io/github/license/nitrictech/nitric?style=for-the-badge">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/nitrictech/nitric/test.yaml?branch=develop&style=for-the-badge">
  <img alt="Codecov" src="https://img.shields.io/codecov/c/github/nitrictech/nitric?style=for-the-badge">
  <a href="https://nitric.io/chat"><img alt="Discord" src="https://img.shields.io/discord/955259353043173427?label=discord&style=for-the-badge"></a>
</p>

[Nitric](https://nitric.io) is a multi-language framework, with concise inline infrastructure from code. Modern applications should be robust, productive and a joy to build. Nitric solves common problems building for modern platforms:

- [Easy infrastructure](https://nitric.io/docs/get-started/foundations/why-nitric) from code
- Build for [any host without coupling](https://nitric.io/docs/providers)
- [Run locally](https://nitric.io/docs/get-started/foundations/projects/local-development#local-dashboard)
- [IAM for humans](https://nitric.io/docs/get-started/foundations/infrastructure/security)
- Common resources like [databases](https://nitric.io/docs/sql), [queues/topics](https://nitric.io/docs/messaging), [APIs](https://nitric.io/docs/apis), [key-value](https://nitric.io/docs/keyvalue), [buckets](https://nitric.io/docs/storage) and more
- [Change services, IaC tools or cloud providers](https://nitric.io/docs/providers) without changing code

We also know abstraction should mean building on existing layers, not hiding them. Nitric includes powerful escape hatches for when things get custom.

## Supported Languages

<p>
  <a href="https://github.com/nitrictech/node-sdk"><img src="https://skillicons.dev/icons?i=js"/></a>
  <a href="https://github.com/nitrictech/node-sdk"><img src="https://skillicons.dev/icons?i=ts"/></a>
  <a href="https://github.com/nitrictech/python-sdk"><img src="https://skillicons.dev/icons?i=py"/></a>
  <a href="https://github.com/nitrictech/go-sdk"><img src="https://skillicons.dev/icons?i=go"/></a>
  <a href="https://github.com/nitrictech/dart-sdk"><img src="https://skillicons.dev/icons?i=dart"/></a>
</p>

## Supported Clouds

<p>
  <a href="./cloud/aws"><img src="https://skillicons.dev/icons?i=aws"/></a>
  <a href="./cloud/gcp"><img src="https://skillicons.dev/icons?i=gcp"/></a>
  <a href="./cloud/azure"><img src="https://skillicons.dev/icons?i=azure"/></a>
</p>

> These are supported out of the box, but you can also build [custom providers](https://nitric.io/docs/providers/custom/create) as well

## 🧑‍💻 Get started

💿 **Install Nitric:**

**macOS**:

```
brew install nitrictech/tap/nitric
```

**Linux**:

```
curl -L "https://nitric.io/install?version=latest" | bash
```

**Windows**:

```
scoop bucket add nitric https://github.com/nitrictech/scoop-bucket.git
scoop install nitric
```

🚀 **Start building your first app**:

```
nitric new
```

🕹 **See our example apps**: [Example Apps Repo](https://github.com/nitrictech/examples).

📚 **Prefer a walkthrough?** Read through our [guides](https://nitric.io/docs/guides).

👋 **Any questions?** Join our developer community on [Discord](https://nitric.io/chat).

⭐ **Give us a star** to help support our work!

## 🍿 Visual Learner?

To get up to speed quickly, take a look at our [quick intro to Nitric](https://www.youtube.com/watch?v=Hljs7Ei9SIs).

<a href="https://www.youtube.com/watch?v=Hljs7Ei9SIs">
  <img width="600px" src="https://img.youtube.com/vi/Hljs7Ei9SIs/maxresdefault.jpg"/>
</a>

## 🤷 So.. how does it work?

Nitric focuses on what you want to achieve as the developer:

_What workflow do you need to be productive?_

_What system design are you trying to achieve?_.

All you need to do is write your application code and your infrastructure requirements are inferred. Nitric then orchestrates and configures the deployment of your application, no need to manually write your Terraform or other IaC code. By abstracting these infrastructure requirements, it removes the need to write boilerplate and means your single application is portable across clouds including, AWS, GCP, and Azure.

And, it's all **open-source**

## 📝 Example: Note Taking

Creating production-ready services and resources is simple, with less than 10 lines to deploy an API endpoint and a bucket with all the IAM permissions automatically configured.

```javascript
import { api, bucket } from "@nitric/sdk";

const main = api("main");
const notes = bucket("notes").allow("read", "write");

main.post("/notes/:title", async (ctx) => {
  const { title } = ctx.req.params;
  await notes.file(title).write(ctx.req.text());
});
```

This is the only code needed to deploy a working application to any cloud provider using [`nitric up`](https://nitric.io/docs/get-started/foundations/deployment). Nitric can deploy this application using automatically generated [Pulumi](https://nitric.io/docs/providers/pulumi), [Terraform](https://nitric.io/docs/providers/terraform) or [any other automation tools](https://nitric.io/docs/providers/custom/create) of your choice.

## Why use Nitric?

1. **Developer-Centric Workflow** Nitric lets you design your application architecture, independent of the deployment automation tool or target platform. With highly declarative in-app infrastructure requirements.

2. **Making Implicit Requirements Explicit** If your app needs storage, a database, or a message queue, Nitric ensures these resources are properly set up and integrated into your app, removing the friction of manual configuration.

3. **Cloud-Agnostic and Portable** Nitric decouples your application from the underlying cloud infrastructure. Whether you're using AWS, Azure, GCP, or Kubernetes, Nitric allows you to map your application's requirements to the appropriate services across platforms.

4. **Automated Infrastructure, Best Practices Included** One of the most error-prone aspects of cloud development is managing permissions, configurations, and security policies. Nitric automates this, making security best practices—like least privilege access and proper service configurations easy.

5. **Focus on Application Logic** Nitric's approach allows you to focus on building your application, instead of the scaffolding required to run it in the cloud. By removing the manual steps from the IaC process, Nitric eliminates significant boilerplate and reduces the runtime checking needed to handle configuration errors.

6. **Plugin-Based Architecture** Nitric's plugin-based architecture allows you to use the deployment plugins we provide, which use Pulumi or Terraform for deployment, or write your own. This flexibility allows you to use the tools you're comfortable with, while still benefiting from Nitric's infrastructure automation and cloud-agnostic approach.

## Want more?

Nitric has full documentation at [nitric.io/docs](https://nitric.io/docs), including concepts, reference documentation for various languages and many [tutorials/guides](https://nitric.io/docs/guides).

- Ask questions in [GitHub discussions](https://github.com/nitrictech/nitric/discussions)

- Join us on [Discord](https://nitric.io/chat)

- Find us on [X](https://x.com/nitric_io)

- Or send us an [email](mailto:maintainers@nitric.io)

## Contributing

We greatly appreciate contributions, consider starting with the [contributions guide](./CONTRIBUTING.md) or [development guide](./DEVELOPERS.md), and a chat on [Discord](https://nitric.io/chat) or [GitHub](https://github.com/nitrictech/nitric/discussions).
