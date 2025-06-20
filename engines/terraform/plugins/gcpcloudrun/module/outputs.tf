output "nitric" {
    value = {
        id = google_cloud_run_v2_service.service.name
        domain_name = google_cloud_run_v2_service.service.uri
        exports = {
            resources = {
                "google_cloud_run_v2_service" = google_cloud_run_v2_service.service.name
            }
        }
    }
}

