# Nitric Repository Development Guidelines

## Getting Started

Get started with Nitric development by forking the repository and cloning it to your local machine.

```bash
git clone https://github.com/<your-github-username>/nitric.git
```

### Installation

```bash
go mod tidy
```

Requirements:

- Git
- Golang (1.22)
- Make
- Docker
- Google Protocol Buffers Compiler (protoc)

To install dependencies use:

```bash
make install-tools
```

### Building

```bash
make binaries
```

To install the binaries into your nitric home use:

```bash
make install
```

In your stack file you can then use the version by specifying `0.0.1`, i.e.

```yaml
provider: nitric/aws@0.0.1
```

### Testing

Run unit tests using:

```bash
make test
```

Run integration tests using:

```bash
make test-integration
```

It can be useful to run a provider in a 'service only' mode to test the runtime or deployment components, where the cloud APIs are available but you don't need/want to start a child process to handle incoming request. This can be achieved by setting the MIN_WORKERS variable to `0`:

```bash
export MIN_WORKERS=0; ./cloud/aws/bin/deploy-aws
```

## Custom Providers

There are two main ways to create a custom provider:

1. Create a new provider: This is the most flexible option, but also the most complex. You can [create a new provider from scratch](https://nitric.io/docs/providers/custom/create).
2. Extend an existing provider: This is a good option if you want to leverage the existing provider's deployment automation and only need to make specific changes, such as use your own Terraform modules or deploy Nitric resources to a different cloud service. You [can extend an existing provider](https://nitric.io/docs/providers/custom/extend) to add your own configuration options or change the deployment process.

## Community Channels

- Ask questions in [GitHub discussions](https://github.com/nitrictech/nitric/discussions)
- Join us on [Discord](https://nitric.io/chat)
- Find us on [Twitter](https://twitter.com/nitric_io)
- Send us an [email](mailto:maintainers@nitric.io)

# Documentation Contribution

If you find a mistake or are interested in contributing to our documentation you can fork the [documentation repo](https://github.com/nitrictech/docs), clone to your local machine and then open a pull request. If you found a problem but don't have the time to make the changes directly, then a [opening up an issue](https://github.com/nitrictech/docs/issues/new/choose) is much appreciated.

```bash
git clone https://github.com/<your-github-username>/docs.git
```

The docs repo is organised with all the documentation being under `/docs` and the images under `/public/images`.

### Formatting

All docs files are written using markdown, with some custom components written for the rendering ([shown below](#components)).

When you open a pull request, tests will run in the GitHub actions that will spellcheck, check for broken links, and make sure all the formatting is correct. These scripts can be run locally using the following commands:

```bash
npm run test:spellchecker

npm run format:fix

npm run cypress
```

If there is a word that is flagged by the spellchecker but is actually valid you can update the `dictionary.txt` file.

### Components

There are a few components that have been written specifically for the Nitric documentation. These require that the markdown is written in `.mdx` format. These components are listed below.

If you would like to highlight something important you can use the `<Note>` component. This is normally reserved for common pitfalls, caveats, or warnings that a reader should take not of.

```
<Note>Content that you want to be highlighted</Note>
```

If you would like to have a title above your code examples, you can add metadata to the code markdown. This is generally reserved for adding file names to the tops of examples.

````
```js title:example/document.js

```
````

For tabs to appear with the code examples with switchable languages, you can use the `<CodeSwitcher>` component.

````
<CodeSwitcher>

```js title:example/document.js

```

```python title:example/document.py

```

</CodeSwitcher>
````

### Frontmatter

There is metadata that is required for every docs page. This metadata is stored in the frontmatter at the top of the page. For a default docs page it looks like so:

```
---
description: 'Basic description for this page of example documentation'
---
```

The guides require a bit more metadata as there is filtering that needs to happen on the [guides page](https://nitric.io/docs/guides).

```
---
description: 'Description for this guide'
tags:
  - A tag that can be used for filtering the guide
languages:
  - javascript
  - typescript
  - python
  - dart
  - go
image: /docs/images/guides/example-guide/banner.png
image_alt: 'The featured image for the Example Guide'
featured:
  image: /docs/images/guides/example-guide/featured.png
  image_alt:
published_at: 2024-12-21
updated_at: 2024-12-21
---
```
