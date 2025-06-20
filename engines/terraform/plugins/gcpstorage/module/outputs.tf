output "nitric" {
    value = {
        id = google_storage_bucket.bucket.id
        exports = {
            resources = {
                "google_storage_bucket" = google_storage_bucket.bucket.id
            }
        }
    }
}