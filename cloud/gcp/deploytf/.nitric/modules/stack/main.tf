# Generate a random id for the bucket
resource "random_id" "stack_id" {
  byte_length = 4

  prefix = "${var.stack_name}-"

  depends_on = [ google_kms_crypto_key_iam_binding.cmek_key_binding[0] ]
}

module "iam_roles" {
  source = "../roles"
}

locals {
  full_stack_id = "${random_id.stack_id.hex}"
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
    # Enable Artifact Registry API and Container Registry API
    "containerregistry.googleapis.com",
    "artifactregistry.googleapis.com",
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
    "serviceusage.googleapis.com",
    # Enable KMS service
    "cloudkms.googleapis.com",
  ]
}

# Enable the required services
resource "google_project_service" "required_services" {
  for_each = toset(local.required_services)

  service = each.key
  # Leave API enabled on destroy
  disable_on_destroy         = false
  disable_dependent_services = false
}

# Get the GCP project number
data "google_project" "project" {
}

resource "google_project_iam_member" "pubsub_token_creator" {
  project    = data.google_project.project.project_id
  role       = "roles/iam.serviceAccountTokenCreator"
  member     = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
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
  role_id     = "${replace(random_id.stack_id.hex, "-", "_")}_svc_base_role"
  title       = "${random_id.stack_id.hex} service base role"
  permissions = local.base_compute_permissions
}

# Deploy a artifact registry repository
resource "google_artifact_registry_repository" "service-image-repo" {
  location      = var.location
  repository_id = "${local.full_stack_id}-services"
  description   = "service images for nitric stack ${var.stack_name}"
  kms_key_name  = var.cmek_enabled ? google_kms_crypto_key.cmek_key[0].id : null
  format        = "DOCKER"
  depends_on = [ google_kms_crypto_key_iam_binding.cmek_key_binding[0] ]
}

resource "random_id" "random_kms_id" {
  byte_length = 4

  prefix = "${var.stack_name}-"
}

# Deploy a KMS keyring and key if cmek enabled
# TODO: May want multiple keys for different services
resource "google_kms_key_ring" "cmek_key_ring" {
  count    = var.cmek_enabled ? 1 : 0
  location = var.location
  name     = "${random_id.random_kms_id.hex}-key-ring"
}

resource "google_kms_crypto_key" "cmek_key" {
  count    = var.cmek_enabled ? 1 : 0
  name     = "${random_id.random_kms_id.hex}-key-ring"
  key_ring = google_kms_key_ring.cmek_key_ring[0].id

  # lifecycle {
  #   prevent_destroy = true
  # }
}

resource "google_project_service_identity" "secret_manager_sa" {
  count    = var.cmek_enabled ? 1 : 0
  provider = google-beta

  project = data.google_project.project.project_id
  service = "secretmanager.googleapis.com"
}

locals {
  kms_reader_service_accounts = [
    // Artifact registry service account
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-artifactregistry.iam.gserviceaccount.com",
    // Pubsub service account
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com",
    // Cloud run service account
    "serviceAccount:service-${data.google_project.project.number}@gs-project-accounts.iam.gserviceaccount.com",
    // Cloud scheduler service account
    "serviceAccount:service-${data.google_project.project.number}@serverless-robot-prod.iam.gserviceaccount.com",
    // Cloud scheduler service account
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  ]
}

# Allow the project service account access to the kms key ring
resource "google_kms_crypto_key_iam_binding" "cmek_key_binding" {
  count         = var.cmek_enabled ? 1 : 0
  crypto_key_id = google_kms_crypto_key.cmek_key[0].id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members       = toset(local.kms_reader_service_accounts)
  depends_on    = [google_project_service.required_services, google_project_service_identity.secret_manager_sa[0]]
}
