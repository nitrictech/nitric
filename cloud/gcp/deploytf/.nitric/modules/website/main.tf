# Create the GCP bucket
resource "google_storage_bucket" "website_bucket" {
  name     = var.bucket_name
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

locals {
  # Apply the base path logic for key transformation
  transformed_files = {
    for path, file in module.template_files.files : (
      var.base_path == "/" ? 
        path : 
        "${trimsuffix(var.base_path, "/")}/${path}"
    ) => file
  }
}


# Upload files from the local directory to the bucket
resource "google_storage_bucket_object" "website_files" {
  for_each = local.transformed_files

  name        = trimprefix(each.key, "/")
  bucket      = google_storage_bucket.website_bucket.name
  source                 = each.value.source_path
  content_type           = each.value.content_type
}

