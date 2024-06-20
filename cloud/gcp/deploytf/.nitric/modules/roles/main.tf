# Generate a set of known custom IAM roles
# That translate to nitric permissions
# For a given project this would only need to be done once for all nitric stacks deployed to that project

# Generate a random id the nitric roles
resource "random_id" "role_id" {
  byte_length = 4
}

# Permissions required for compute units to operate
resource "google_project_iam_custom_role" "base_compute_role" {
  role_id     = "NitricBaseCompute_${random_id.role_id.hex}"
  title       = "Nitric Base Compute"
  description = "Custom role for base nitric compute permissions"
  permissions = [
    "storage.buckets.list",
    "storage.buckets.get",
    "cloudtasks.queues.get",
    "cloudtasks.tasks.create",
    "cloudtrace.traces.patch",
    "monitoring.timeSeries.create",
    "iam.serviceAccounts.signBlob",
    "pubsub.topics.list",
    "pubsub.topics.get",
    "pubsub.snapshots.list",
    "pubsub.subscriptions.get",
    "resourcemanager.projects.get",
    "apigateway.gateways.list",
    "secretmanager.secrets.list",
  ]
}

# Permissions required for reading from a bucket
resource "google_project_iam_custom_role" "bucket_reader_role" {
  role_id     = "NitricBucketReader_${random_id.role_id.hex}"
  title       = "Nitric Bucket Reader"
  description = "Custom role that only allows reading from a bucket"
  permissions = ["storage.objects.get", "storage.objects.list"]
}

# Permissions required to write to a bucket
resource "google_project_iam_custom_role" "bucket_writer_role" {
  role_id     = "NitricBucketWriter_${random_id.role_id.hex}"
  title       = "Nitric Bucket Writer"
  description = "Custom role that only allows writing to a bucket"
  permissions = ["storage.objects.create", "storage.objects.delete"]
}

# Permissions required to delete an item from a bucket
resource "google_project_iam_custom_role" "bucket_deleter_role" {
  role_id     = "NitricBucketDeleter_${random_id.role_id.hex}"
  title       = "Nitric Bucket Deleter"
  description = "Custom role that only allows deleting from a bucket"
  permissions = ["storage.objects.delete"]
}

# Permissions required to access a secret
resource "google_project_iam_custom_role" "secret_access_role" {
  role_id     = "SecretAccessRole_${random_id.role_id.hex}"  
  title       = "Secret Access Role"
  permissions = ["resourcemanager.projects.get",
		"secretmanager.locations.get",
		"secretmanager.locations.list",
		"secretmanager.secrets.get",
		"secretmanager.secrets.getIamPolicy",
		"secretmanager.versions.get",
		"secretmanager.versions.access",
		"secretmanager.versions.list",
    ]
}

# Permissions required to put a secret
resource "google_project_iam_custom_role" "secret_put_role" {
  role_id     = "SecretPutRole_${random_id.role_id.hex}"  
  title       = "Secret Put Role"
  permissions = ["resourcemanager.projects.get",
		"secretmanager.versions.add",
		"secretmanager.versions.enable",
		"secretmanager.versions.destroy",
		"secretmanager.versions.disable",
		"secretmanager.versions.get",
		"secretmanager.versions.access",
		"secretmanager.versions.list",
    ]
}

# Permissions required to delete a kv
resource "google_project_iam_custom_role" "kv_deleter_role" {
  role_id     = "KVDeleteRole_${random_id.role_id.hex}"  
  title       = "KV Delete Role"
  permissions = ["resourcemanager.projects.get",
		"appengine.applications.get",
		"datastore.databases.get",
		"datastore.indexes.get",
		"datastore.namespaces.get",
		"datastore.entities.delete",
    ]
}

# Permissions required to read a kv
resource "google_project_iam_custom_role" "kv_reader_role" {
  role_id     = "KVReadRole_${random_id.role_id.hex}"  
  title       = "KV Read Role"
  permissions = ["resourcemanager.projects.get",
		"appengine.applications.get",
		"datastore.databases.get",
		"datastore.entities.get",
		"datastore.indexes.get",
		"datastore.namespaces.get",
		"datastore.entities.list",
    ]
}

# Permissions required to write a kv
resource "google_project_iam_custom_role" "kv_writer_role" {
  role_id     = "KVWriteRole_${random_id.role_id.hex}"  
  title       = "KV Write Role"
  permissions = ["resourcemanager.projects.get",
		"appengine.applications.get",
		"datastore.indexes.list",
		"datastore.namespaces.list",
		"datastore.entities.create",
		"datastore.entities.update",
    ]
}

# Permissions required to enqueue to a queue
resource "google_project_iam_custom_role" "queue_enqueue_role" {
  role_id     = "QueueEnqueue_${random_id.role_id.hex}"  
  title       = "Queue Enqueue"
  permissions = ["pubsub.topics.get",
		"pubsub.topics.publish",
    ]
}

# Permissions required to dequeue from a queue
resource "google_project_iam_custom_role" "queue_dequeue_role" {
  role_id     = "QueueDequeue_${random_id.role_id.hex}"  
  title       = "Queue Dequeue"
  permissions = ["pubsub.topics.get",
		"pubsub.topics.attachSubscription",
		"pubsub.snapshots.seek",
		"pubsub.subscriptions.consume",
    ]
}