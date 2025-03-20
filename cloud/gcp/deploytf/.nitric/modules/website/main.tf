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

# Upload files from the local directory to the bucket
resource "google_storage_bucket_object" "website_files" {
  for_each = fileset(var.local_directory, "**") # Recursively get all files in the directory

  name        = each.value
  bucket      = google_storage_bucket.website_bucket.name
  source      = "${var.local_directory}/${each.value}"
  content_type = lookup(
    {
      ".html" = "text/html",
      ".css"  = "text/css",
      ".js"   = "application/javascript",
      ".png"  = "image/png",
      ".jpg"  = "image/jpeg",
      ".jpeg" = "image/jpeg",
      ".gif"  = "image/gif",
      ".svg"  = "image/svg+xml"
    },
    regex("\\.[^.]+$", each.value), # Match file extension
    "application/octet-stream"     # Default content type
  )
}