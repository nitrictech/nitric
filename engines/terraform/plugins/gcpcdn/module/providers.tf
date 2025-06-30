terraform {
  required_providers {
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 6.17.0"
    }

    corefunc = {
      source  = "northwood-labs/corefunc"
      version = "~> 1.4"
    }
  }
}

provider "google-beta" {
  project = var.project_id
}