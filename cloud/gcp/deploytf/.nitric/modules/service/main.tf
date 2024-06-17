terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

# Create a GCR repository for the service image
resource "google_container_registry_repository" "repo" {
  name = var.service_name
}

data "google_client_config" "gcp_config" {
}

provider "docker" {
  registry_auth {
    address  = "https://gcr.io"
    username = "oauth2accesstoken"
    password = data.google_client_config.access_token
  }
}

# Tag the provided docker image with the repository url
resource "docker_tag" "tag" {
  source_image = var.image
  target_image = google_container_registry_repository.repo.repository_url
}

# Push the tagged image to the repository
resource "docker_registry_image" "push" {
  name = google_container_registry_repository.repo.repository_url
  triggers = {
    source_image_id = docker_tag.tag.source_image_id
  }
}

# Create a service account for the google cloud run instance
resource "google_service_account" "service_account" {
  account_id   = var.service_name
  display_name = var.service_name
}

# Create a random password for events that will target this service
resource "random_password" "password" {
  length  = 32
  special = true
  keepers = {
    "name" = var.service_name
  }
}

# Create a cloud run service
resource "google_cloud_run_service" "service" {
  name     = var.service_name
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
      service_account_name = google_service_account.service_account.email
      container_concurrency = var.container_concurrency
      timeout_seconds = var.timeout_seconds
      containers {
        env {
          name  = "EVENT_TOKEN"
          value = random_password.password.result
        }
        env {
          name = "SERVICE_ACCOUNT_EMAIL"
          value = google_service_account.service_account.email
        }
        env {
          name = "GCP_REGION"
          value = var.region
        }

        dynamic env {
          for_each = var.environment
          content {
            name  = env.key
            value = env.value
          }
        }
        image = "${google_container_registry_repository.repo.repository_url}@${docker_registry_image.push.sha256_digest}"
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

# Create an invoker service account for the google cloud run instance
resource "google_service_account" "invoker_service_account" {
  account_id   = "${var.service_name}-cr-invoker"
}

# Give the above service account permissions to execute the CloudRun service
resource "google_cloud_run_service_iam_member" "invoker" {
  service = google_cloud_run_service.service.name
  location = google_cloud_run_service.service.location
  role = "roles/run.invoker"
  member = "serviceAccount:${google_service_account.invoker_service_account.email}"
}
