# Claude Development Guidelines

## Build Rules

### CLI Changes

When making changes to files in the `cli/` directory, always run the build command after the changes are made to ensure the CLI is properly compiled:

```bash
cd cli && make
```

This should be done after any modifications to:

- Go source files in `cli/`
- Configuration files affecting the CLI build
- Any other files that impact the CLI functionality

The build ensures that changes are compiled and the binary at `cli/bin/nitric` is up to date.
