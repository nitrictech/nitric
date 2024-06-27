terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
    }
  }
}

# Create a GCR repository for the service image
data "google_container_registry_repository" "repo" {
  project = var.project_id
}

locals {
  service_image_url = "${data.google_container_registry_repository.repo.repository_url}/${var.service_name}"
}

# Tag the provided docker image with the repository url
resource "docker_tag" "tag" {
  source_image = var.image
  target_image = local.service_image_url
}

# Push the tagged image to the repository
resource "docker_registry_image" "push" {
  name = local.service_image_url
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
  display_name = "${var.service_name} service account"
  description  = "Service account which runs the ${var.service_name} service"
}

# Create a random password for events that will target this service
resource "random_password" "event_token" {
  length  = 32
  special = false
  keepers = {
    "name" = var.service_name
  }
}

# Create a cloud run service
resource "google_cloud_run_service" "service" {
  name = replace(var.service_name, "_", "-")

  location = var.region
  project  = var.project_id

  template {
    metadata {
      annotations = {
        // TODO: Add configuration here
        "autoscaling.knative.dev/minScale" = "0"
        "autoscaling.knative.dev/maxScale" = "100"
      }
    }
    spec {
      service_account_name  = google_service_account.service_account.email
      container_concurrency = var.container_concurrency
      timeout_seconds       = var.timeout_seconds
      containers {
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
          for_each = var.environment
          content {
            name  = env.key
            value = env.value
          }
        }
        image = "${local.service_image_url}@${docker_registry_image.push.sha256_digest}"
        ports {
          container_port = 9001
        }
        resources {
          limits = {
            # TODO: enable cpu configuration
            # cpu    = "1000m"
            memory = "${var.memory_mb}Mi"
          }
        }
      }
    }
  }

  depends_on = [docker_registry_image.push]
}

# Create a random ID for the service name, so that it confirms to regex restrictions
resource "random_string" "service_id" {
  length  = 30 - length(local.ids_prefix)
  special = false
  upper   = false
}

# Create an invoker service account for the google cloud run instance
resource "google_service_account" "invoker_service_account" {
  project      = var.project_id
  account_id   = "${local.ids_prefix}${random_string.service_id.id}"
  display_name = "${var.service_name} invoker"
  description  = "Service account which allows other resources to invoke the ${var.service_name} service"
}

# Give the above service account permissions to execute the CloudRun service
resource "google_cloud_run_service_iam_member" "invoker" {
  service  = google_cloud_run_service.service.name
  location = google_cloud_run_service.service.location
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.invoker_service_account.email}"
}

resource "google_project_iam_member" "project_member" {
  project = var.project_id
  member  = "serviceAccount:${google_service_account.service_account.email}"
  # p.BaseComputeRole.Name,
  role = var.base_compute_role
}

# Give the above service account permissions to act as itself
resource "google_service_account_iam_member" "account_member" {
  service_account_id = google_service_account.invoker_service_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.invoker_service_account.email}"
}
