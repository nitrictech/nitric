output "nitric" {
    value = {
        id = google_storage_bucket.bucket.id
        domain_name = google_storage_bucket.bucket.url
        exports = {
            resources = {
                "google_storage_bucket" = google_storage_bucket.bucket.id
            }
        }
    }
}