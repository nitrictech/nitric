# Nitric URI Spec

Nitric uris are used to define a reference to a unique reference within the scope of a single nitric stack.

The basic structure of a nitric uri is

`nitric:<resource_type>:<resource_name>?<query...>`

## Additional URI specifications

### Bucket

`nitric:bucket:<bucket_name>?<bucket_query>`

#### Queries

##### event

`write` AND/OR `delete`

##### file

Defined as a string, can be an absolute path to a file on the bucket or a pattern

### Topic

`nitric:topic:<topic_name>`

#### Queries

There are no currently valid queries that can be defined for a topic