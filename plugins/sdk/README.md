# Nitric Membrane SDK

The Nitric SDK is a collection of golang interfaces that are used to abstract services provided by Nitric.

## Gateway Plugin

Gateway plugins are used to abstract implementation detail for application input/output. This usually involves presenting some kind of external interface (for example a HTTP server), normalizing incoming data from this interface to pass to the hosted application (over HTTP) and then denormalizing output to be returned back to the hosted environment. 

Below are some concrete examples:

* [AWS Lambda](../aws/gateway/lambda/README.md)







