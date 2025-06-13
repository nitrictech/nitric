output "nitric" {
    value = {
        id = google_cloud_run_v2_service.service.name
        http_endpoint = google_cloud_run_v2_service.service.uri
    }
}

