# Create the GCP bucket
resource "google_storage_bucket" "website_bucket" {
  name     = "${var.stack_id}-${var.website_name}"
  location = var.region

  website {
    main_page_suffix = var.index_document
    not_found_page   = var.error_document
  }
}

# Set public access permissions for the bucket
resource "google_storage_bucket_iam_binding" "website_bucket_iam" {
  bucket = google_storage_bucket.website_bucket.name
  role   = "roles/storage.objectViewer"

  members = [
    "allUsers"
  ]
}
module "template_files" {
  source  = "hashicorp/dir/template"
  version = "1.0.2"

  base_dir = var.local_directory
}

# Upload files from the local directory to the bucket
resource "google_storage_bucket_object" "website_files" {
  for_each = module.template_files.files

  name         = trimprefix(each.key, "/")
  bucket       = google_storage_bucket.website_bucket.name
  source       = each.value.source_path
  content_type = each.value.content_type
}

