# Configuring the Membrane

In most cases these will already be set as environment variables withing a nitric templates Dockerfile. This documentation is targeted at developers building their own nitric application templates and those who require the membrane to be run on a bare-metal server.

## Options

| Environment Variable | Description | Default |
| --- | --- | --- |
| MEMBRANE_MODE | Sets the operating mode of the membrane, see [here](./operating-modes.md) for available options | `FAAS` | 
| SERVICE_ADDRESS | Sets the address that the membrane APIs should be bound to is configured as single string `host:port` | `127.0.0.1:50051` | 
| CHILD_ADDRESS | Sets the address that the child process will be listening on, for requests from the membrane | `127.0.0.1:8080` |
| INVOKE | Sets the command for the child process that the membrane will execute to begin the child process server | `none` |
| TOLERATE_MISSING_SERVICES | Enables/Disables the membranes ability to run with an incomplete set of plugins | `false` |
