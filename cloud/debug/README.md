# Debug Provider

The debug provider doesn't deploy resources to cloud environments. Its purpose is to
assist with validating and troubleshooting the resource graph generated by Nitric for a
specific project.

## Usage

The debug provider isn't distributed as an official Nitric provider, so it must be built and installed from source.

```bash
# build and install the debug provider binary
make install
```

Next, create or update a project `stack` file. The file can reference the locally installed provider.

`nitric-debug.yaml`
```yaml
name: debug-stack
provider: debug/spec@0.0.1
output: ./debug-output.json
```

> Note: the `output` property specifies where to write the debug details on disk.

Use the `nitric up` command to 'deploy' this new debug stack.

```
nitric up -s debug
```

When the up command completes the full stack details (resource graph) should be output to the specified file.