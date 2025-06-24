locals {
  required_services = [
    // Enable Certificate Manager API
    "certificatemanager.googleapis.com",
    // Enable DNS API
    "dns.googleapis.com",
    // Enable Compute API (Networking/Load Balancing)
    "compute.googleapis.com"
  ]

  root_origins = {
    for k, v in var.nitric.origins : k => v
    if v.path == "/"
  }

  default_origin = length(local.root_origins) > 0 ? keys(local.root_origins)[0] : keys(var.nitric.origins)[0]

  cloud_storage_origins = {
    for k, v in var.nitric.origins : k => v
    if contains(keys(v.resources), "google_storage_bucket")
  }

  cloud_run_origins = {
    for k, v in var.nitric.origins : k => v
    if contains(keys(v.resources), "google_cloud_run_v2_service")
  }

  other_origins = {
    for k, v in var.nitric.origins : k => v
    if !contains(keys(v.resources), "google_storage_bucket") && !contains(keys(v.resources), "google_cloud_run_v2_service")
  }
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

resource "random_string" "cdn_prefix" {
  length  = 8
  special = false
  lower = true
  upper = false
}

# Create Network Endpoint Groups for services
resource "google_compute_region_network_endpoint_group" "service_negs" {
  for_each = local.cloud_run_origins
  provider = google-beta

  name                  = "${each.key}-service-neg"
  region                = var.region
  network_endpoint_type = "SERVERLESS"

  cloud_run {
    service = each.value.id
  }

  depends_on = [ google_project_service.required_services ]
}

# Create Backend Services for services
resource "google_compute_backend_service" "service_backends" {
  for_each = local.cloud_run_origins

  project = var.project_id

  name     = "${provider::corefunc::str_kebab(each.key)}-service-bs"
  protocol = "HTTPS"
  enable_cdn  = false

  backend {
    group = google_compute_region_network_endpoint_group.service_negs[each.key].self_link
  }
}

# Create a Global IP Address for the CDN
resource "google_compute_global_address" "cdn_ip" {
  name = "cdn-ip-${random_string.cdn_prefix.result}"
  project = var.project_id
}

data "google_storage_bucket" "bucket" {
  for_each = local.cloud_storage_origins

  name = each.value.id
}

resource "google_storage_bucket_iam_binding" "website_bucket_iam" {
  for_each = local.cloud_storage_origins

  bucket = data.google_storage_bucket.bucket[each.key].name
  role   = "roles/storage.objectViewer"

  members = [
    "allUsers"
  ]
}

resource "google_compute_backend_bucket" "website_backends" {
  for_each = local.cloud_storage_origins

  project = var.project_id

  name        = "${provider::corefunc::str_kebab(each.key)}-site-bucket"
  bucket_name = data.google_storage_bucket.bucket[each.key].name
  enable_cdn  = true
}

resource "google_compute_region_network_endpoint_group" "external_negs" {
  for_each = local.other_origins
  provider = google-beta

  name                  = "${each.key}-external-neg"
  network_endpoint_type = "INTERNET_FQDN_PORT"
  region = var.region

  connection {
    host = each.value.domain_name
    port = 443
  }

  depends_on = [ google_project_service.required_services ]
}

resource "google_compute_backend_service" "external_backends" {
  for_each = local.other_origins

  project = var.project_id

  name     = "${provider::corefunc::str_kebab(each.key)}-external-bs"
  protocol = "HTTPS"
  enable_cdn  = false

  backend {
    group = google_compute_region_network_endpoint_group.external_negs[each.key].self_link
  }
}

# Create a URL Map for routing requests
resource "google_compute_url_map" "https_url_map" {
  name            = "https-site-url-map-${random_string.cdn_prefix.result}"
  project = var.project_id
  default_service = google_compute_backend_service.service_backends[local.default_origin].self_link

  host_rule {
    hosts        = ["*"]
    path_matcher = "all-paths"
  }

  path_matcher {
    name            = "all-paths"
    default_service = google_compute_backend_service.service_backends[local.default_origin].self_link

    dynamic "path_rule" {
      for_each = local.cloud_run_origins

      content {
        service = google_compute_backend_service.service_backends[path_rule.key].self_link
        paths   = [
          startswith(path_rule.value.path, "/") ? "/${path_rule.value.base_path}${path_rule.value.path}/*" : "/${path_rule.value.base_path}/${path_rule.value.path}/*", // Ensure /${path}/*
          startswith(path_rule.value.path, "/") ? "/${path_rule.value.base_path}${path_rule.value.path}" : "/${path_rule.value.base_path}/${path_rule.value.path}"] // Ensure /${path}
        
        route_action {
          url_rewrite {
            path_prefix_rewrite = "/"
          }
        }

      }
    }

    dynamic "path_rule" {
      for_each = local.cloud_storage_origins

      content {
        service = google_compute_backend_bucket.website_backends[path_rule.key].self_link
        paths   = [
          startswith(path_rule.value.path, "/") ? "/${path_rule.value.base_path}${path_rule.value.path}/*" : "/${path_rule.value.base_path}/${path_rule.value.path}/*", // Ensure /${path}/*
          startswith(path_rule.value.path, "/") ? "/${path_rule.value.base_path}${path_rule.value.path}" : "/${path_rule.value.base_path}/${path_rule.value.path}"] // Ensure /${path}
        route_action {
          url_rewrite {
            path_prefix_rewrite = "/"
          }
        }
      }
    }

    dynamic "path_rule" {
      for_each = local.other_origins

      content {
        service = google_compute_backend_service.external_backends[path_rule.key].self_link
        paths   = [
          startswith(path_rule.value.path, "/") ? "/${path_rule.value.base_path}${path_rule.value.path}/*" : "/${path_rule.value.base_path}/${path_rule.value.path}/*", // Ensure /${path}/*
          startswith(path_rule.value.path, "/") ? "/${path_rule.value.base_path}${path_rule.value.path}" : "/${path_rule.value.base_path}/${path_rule.value.path}"] // Ensure /${path}
        route_action {
          url_rewrite {
            path_prefix_rewrite = "/"
          }
        }
      }
    }
  }
}

# Lookup the Managed Zone for the CDN Domain
data "google_dns_managed_zone" "cdn_zone" {
  name = var.cdn_domain.zone_name
  project = var.project_id

  depends_on = [ google_project_service.required_services ]
}

# Create DNS Records for the CDN
resource "google_dns_record_set" "cdn_dns_record" {
  name         = endswith(var.cdn_domain.domain_name, ".") ? var.cdn_domain.domain_name : "${var.cdn_domain.domain_name}."
  managed_zone = data.google_dns_managed_zone.cdn_zone.name
  type         = "A"
  rrdatas      = [google_compute_global_address.cdn_ip.address]
  ttl = var.cdn_domain.domain_ttl
  project = var.project_id
}

resource "google_dns_record_set" "www_cdn_dns_record" {
  name         = endswith(var.cdn_domain.domain_name, ".") ? "www.${var.cdn_domain.domain_name}" : "www.${var.cdn_domain.domain_name}."
  managed_zone = data.google_dns_managed_zone.cdn_zone.name
  type         = "A"
  rrdatas      = [google_compute_global_address.cdn_ip.address]
  ttl = var.cdn_domain.domain_ttl
  project = var.project_id
}

resource "google_certificate_manager_certificate" "cdn_cert" {
  provider    = google-beta
  name        = "cdn-cert-${random_string.cdn_prefix.result}"
  description = "Nitric stack CDN SSL certificate"
  scope       = "DEFAULT"

  managed {
    domains = [var.cdn_domain.domain_name]
  }

  depends_on = [ google_project_service.required_services ]
}

resource "google_certificate_manager_certificate_map" "cdn_cert_map" {
  provider    = google-beta
  name        = "cert-map-${random_string.cdn_prefix.result}"
  description = "CDN Certificate Map"

  depends_on = [ google_project_service.required_services ]
}

resource "google_certificate_manager_certificate_map_entry" "cdn_cert_map_entry" {
  name        = "cdn-cert-map-entry-${random_string.cdn_prefix.result}"
  description = "CDN Certificate Map Entry"
  map = google_certificate_manager_certificate_map.cdn_cert_map.name 
  certificates = [google_certificate_manager_certificate.cdn_cert.id]
  matcher = "PRIMARY"
  project = var.project_id
}

# Create a Target HTTPS Proxy
resource "google_compute_target_https_proxy" "https_proxy" {
  name             = "https-proxy-${random_string.cdn_prefix.result}"
  certificate_map = "//certificatemanager.googleapis.com/${google_certificate_manager_certificate_map.cdn_cert_map.id}"
  url_map          = google_compute_url_map.https_url_map.self_link
  project = var.project_id
}

# Create a Global Forwarding Rule
resource "google_compute_global_forwarding_rule" "https_forwarding_rule" {
  name        = "https-forwarding-rule-${random_string.cdn_prefix.result}"
  ip_address  = google_compute_global_address.cdn_ip.address
  ip_protocol = "TCP"
  port_range  = "443"
  target      = google_compute_target_https_proxy.https_proxy.self_link
  project = var.project_id
}