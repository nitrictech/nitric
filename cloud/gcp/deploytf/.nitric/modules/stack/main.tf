# Generate a random id for the bucket
resource "random_id" "stack_id" {
  byte_length = 4

  prefix = "${var.stack_name}-"
}

module "iam_roles" {
  source = "../roles"
}

locals {
  required_services = [
    # Enable the IAM API
    "iam.googleapis.com",
    # Enable cloud run
    "run.googleapis.com",
    # Enable pubsub
    "pubsub.googleapis.com",
    # Enable cloud scheduler
    "cloudscheduler.googleapis.com",
    # Enable cloud scheduler
    "storage.googleapis.com",
    # Enable Compute API (Networking/Load Balancing)
    "compute.googleapis.com",
    # Enable Container Registry API
    "containerregistry.googleapis.com",
    # Enable firestore API
    "firestore.googleapis.com",
    # Enable ApiGateway API
    "apigateway.googleapis.com",
    # Enable SecretManager API
    "secretmanager.googleapis.com",
    # Enable Cloud Tasks API
    "cloudtasks.googleapis.com",
    # Enable monitoring API
    "monitoring.googleapis.com",
    # Enable service usage API
    "serviceusage.googleapis.com"
  ]
}

# Enable the required services
resource "google_project_service" "required_services" {
  for_each = toset(local.required_services)

  service = each.key
  # Leave API enabled on destroy
  disable_on_destroy = false
  disable_dependent_services = false
}

# Get the GCP project number
data "google_project" "project" {
}

resource "google_project_iam_member" "pubsub_token_creator" {
  project = data.google_project.project.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  depends_on = [google_project_service.required_services]
}

locals {
  base_compute_permissions = [
    "storage.buckets.list",
    "storage.buckets.get",
    "cloudtasks.queues.get",
    "cloudtasks.tasks.create",
    "cloudtrace.traces.patch",
    "monitoring.timeSeries.create",
    // permission for blob signing
    // this is safe as only permissions this account has are delegated
    "iam.serviceAccounts.signBlob",
    // Basic list permissions
    "pubsub.topics.list",
    "pubsub.topics.get",
    "pubsub.snapshots.list",
    "pubsub.subscriptions.get",
    "resourcemanager.projects.get",
    "secretmanager.secrets.list",
    "apigateway.gateways.list",

    // telemetry
    "monitoring.metricDescriptors.create",
    "monitoring.metricDescriptors.get",
    "monitoring.metricDescriptors.list",
    "monitoring.monitoredResourceDescriptors.get",
    "monitoring.monitoredResourceDescriptors.list",
    "monitoring.timeSeries.create",
  ]
}

resource "google_project_iam_custom_role" "base_role" {
  role_id      = "${random_id.stack_id.id}_svc_base_role"
  title        = "${random_id.stack_id.id} service base role"
  permissions = local.base_compute_permissions
}