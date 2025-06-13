locals {
  required_services = [
    // Enable Storage API
    "storage.googleapis.com",
    // Enable IAM API
    "iam.googleapis.com"
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

# Generate a random id for the bucket
resource "random_id" "bucket_id" {
  byte_length = 8

  keepers = {
    # Generate a new id each time we switch to a new AMI id
    bucket_name = var.nitric.name
  }
}

# Google Storage bucket
resource "google_storage_bucket" "bucket" {
  name          = "${var.nitric.name}-${random_id.bucket_id.hex}"
  location      = var.region
  project       = var.project_id
  storage_class = var.storage_class
  labels = {
    "x-nitric-name" = var.nitric.name
  }

  depends_on = [ google_project_service.required_services ]
}

# locals {
#   read_actions = ["storage.objects.get", "storage.objects.list"]
#   write_actions = ["storage.objects.create", "storage.objects.delete"]
#   delete_actions = ["storage.objects.delete"]
# }

# resource "google_project_iam_custom_role" "bucket_access_role" {
#   role_id     = "NitricBucketAccess_${random_id.bucket_id.hex}"
#   title       = "Nitric Bucket Access"
#   description = "Custom role that allows access to a bucket"
#   permissions = distinct(concat(
#       contains(each.value.actions, "read") ? local.read_actions : [],
#       contains(each.value.actions, "write") ? local.write_actions : [],
#       contains(each.value.actions, "delete") ? local.delete_actions : []
#     )
#   )
#
#   depends_on = [ google_project_service.required_services ]
# }

# resource "google_storage_bucket_iam_member" "iam_access" {
#   for_each = var.nitric.origins

#   bucket   = google_storage_bucket.bucket.name
#   role     = google_project_iam_custom_role.bucket_access_role.id
#   member   = "serviceAccount:${each.value.service_account}"
# }
