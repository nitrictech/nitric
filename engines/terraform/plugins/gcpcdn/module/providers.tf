terraform {
  required_providers {
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 6.17.0"
    }
  }
}

provider "google-beta" {
  project = var.project_id
}