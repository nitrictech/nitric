output "bucket_name" {
  description = "The name of the bucket."
  value       = google_storage_bucket.website_bucket.name
}

output "index_document" {
  description = "The index document for the bucket."
  value       = one(google_storage_bucket.website_bucket.website).main_page_suffix
}

output "error_document" {
  description = "The error document for the bucket."
  value       = one(google_storage_bucket.website_bucket.website).not_found_page
}

output "file_md5s" {
  description = "md5 values of the uploaded websites files"
  value = local.uploaded_files_md5
}

output "local_directory" {
  description = "The local directory to sync with the bucket."
  value       = var.local_directory
}