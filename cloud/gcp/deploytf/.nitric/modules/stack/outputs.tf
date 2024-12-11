output "stack_id" {
  value = random_id.stack_id.hex
  description = "A unique id for this deployment"
}

output "base_compute_role" {
  value = google_project_iam_custom_role.base_role.id
  description = "The base compute role to use for the service"
}

output "iam_roles" {
  value = module.iam_roles
}

output "container_registry_uri" {
  value = "${var.location}-docker.pkg.dev/${data.google_project.project.project_id}/${google_artifact_registry_repository.service-image-repo.name}"
  description = "The name of the container registry repository"
}

# The KMS key ring
output "cmek_key_ring" {
  value = length(google_kms_key_ring.cmek_key_ring) > 0 ? google_kms_key_ring.cmek_key_ring[0].id : null
  description = "The name of the KMS key ring"
}

output "cmek_key" {
  value = length(google_kms_crypto_key.cmek_key) > 0 ? google_kms_crypto_key.cmek_key[0].id : null
  description = "The name of the KMS key"
}

output "kms_key_iam_binding" {
  value = length(google_kms_crypto_key.cmek_key) > 0 ? google_kms_crypto_key_iam_binding.cmek_key_binding[0] : null
  description = "The IAM binding for the KMS key"
}

output "firestore_database_id" {
  value = google_firestore_database.database[0] != null ? google_firestore_database.database[0].name : "(default)"
  description = "Firestore database for stack"
}