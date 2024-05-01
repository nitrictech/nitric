# Membrane Operating Modes

The membrane operates in a number of different modes. These modes control how the membrane communicated with it's child process. The mode is configured by setting the system environment variable `MEMBRANE_MODE`. This environment variable will typically be implemented in nitric templates, based on their template type. Available modes are:

* FaaS: `MEMBRANE_MODE="FAAS"`
* HTTP Proxy: `MEMBRANE_MODE="HTTP_PROXY"`