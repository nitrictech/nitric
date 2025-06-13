locals {
  required_services = [
    # Enable the IAM API
    "iam.googleapis.com",
    # Enable cloud run
    "run.googleapis.com",
    # Enable Compute API (Networking/Load Balancing)
    "compute.googleapis.com",
    # Enable Artifact Registry API and Container Registry API
    "artifactregistry.googleapis.com",
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
  project = var.project_id
  # Leave API enabled on destroy
  disable_on_destroy = false
  disable_dependent_services = false
}

resource "google_artifact_registry_repository" "service-image-repo" {
  project       = var.project_id
  location      = var.region
  repository_id = "${var.nitric.name}-repo"
  description   = "service images for nitric stack ${var.nitric.name}"
  format        = "DOCKER"

  depends_on = [ google_project_service.required_services ]
}

locals {
  artifact_registry_url = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.service-image-repo.name}"
  service_image_url = "${local.artifact_registry_url}/${var.nitric.name}"
}

# Tag the provided docker image with the repository url
resource "docker_tag" "tag" {
  source_image = var.nitric.image_id
  target_image = local.service_image_url
}

data "google_client_config" "default" {
}


# Push the tagged image to the repository
resource "docker_registry_image" "push" {
  name = local.service_image_url
  auth_config {
    address  = "${var.region}-docker.pkg.dev"
    username = "oauth2accesstoken"
    password = data.google_client_config.default.access_token
  }
  triggers = {
    source_image_id = docker_tag.tag.source_image_id
  }
}

locals {
  ids_prefix = "nitric-"
}

# Create a random ID for the service name, so that it confirms to regex restrictions
resource "random_string" "service_account_id" {
  length  = 30 - length(local.ids_prefix)
  special = false
  upper   = false
}

# Create a service account for the google cloud run instance
resource "google_service_account" "service_account" {
  account_id   = "${local.ids_prefix}${random_string.service_account_id.id}"
  project      = var.project_id
  display_name = "${var.nitric.name} service account"
  description  = "Service account which runs the ${var.nitric.name} service"

  depends_on = [ google_project_service.required_services ]
}

# Create a random password for events that will target this service
resource "random_password" "event_token" {
  length  = 32
  special = false
  keepers = {
    "name" = var.nitric.name
  }
}

# Create a cloud run service
resource "google_cloud_run_v2_service" "service" {
  name = replace(var.nitric.name, "_", "-")

  location = var.region
  project  = var.project_id
  launch_stage        = "GA"
  ingress = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"

  deletion_protection = false

  template {
    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }
    
    containers {
      image = "${local.service_image_url}@${docker_registry_image.push.sha256_digest}"
      resources {
        limits = {
          cpu    = var.cpus
          memory = "${var.memory_mb}Mi"
        }
      }

      ports {
        container_port = var.ingress_port
      }

      env {
        name  = "NITRIC_STACK_ID"
        value = var.nitric.stack_id
      }
      env {
        name  = "INGRESS_PORT"
        value = var.ingress_port
      }
      env {
        name  = "EVENT_TOKEN"
        value = random_password.event_token.result
      }
      env {
        name  = "SERVICE_ACCOUNT_EMAIL"
        value = google_service_account.service_account.email
      }
      env {
        name  = "GCP_REGION"
        value = var.region
      }

      dynamic "env" {
        for_each = merge(var.environment, var.nitric.env)
        content {
          name  = env.key
          value = env.value
        }
      }
    }

    service_account = google_service_account.service_account.email
    timeout         = "${var.timeout_seconds}s"
  }

  depends_on = [docker_registry_image.push, google_project_service.required_services]
}

# Create a random ID for the service name, so that it confirms to regex restrictions
resource "random_string" "service_id" {
  length  = 30 - length(local.ids_prefix)
  special = false
  upper   = false
}

# Give all users permissions to execute the CloudRun service (is ingress only)
resource "google_cloud_run_service_iam_member" "invoker" {
  project  = var.project_id
  service  = google_cloud_run_v2_service.service.name
  location = google_cloud_run_v2_service.service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}