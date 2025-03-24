# Create Network Endpoint Groups for API Gateways
resource "google_compute_region_network_endpoint_group" "api_gateway_negs" {
  provider = google-beta
  for_each = var.api_gateways

  name                  = "${each.key}-apigw-neg"
  region                = each.value.region
  network_endpoint_type = "SERVERLESS"

  serverless_deployment {
    platform = "apigateway.googleapis.com"
    resource = each.value.gateway_id
  }
}

# Create Backend Services for API Gateways
resource "google_compute_backend_service" "api_gateway_backends" {
  for_each = var.api_gateways

  name     = "${each.key}-apigw-bs"
  protocol = "HTTPS"

  backend {
    group = google_compute_region_network_endpoint_group.api_gateway_negs[each.key].self_link
  }
}

# Create Backend Buckets for Websites
resource "google_compute_backend_bucket" "website_backends" {
  for_each = var.website_buckets

  name        = "${each.key}-site-bucket"
  bucket_name = each.value.name
  enable_cdn  = true
}

# Create a Global IP Address for the CDN
resource "google_compute_global_address" "cdn_ip" {
  name = "${var.stack_id}-cdn-ip"
}

# Create a URL Map for routing requests
resource "google_compute_url_map" "https_url_map" {
  name            = "${var.stack_id}-https-site-url-map"
  default_service = google_compute_backend_bucket.website_backends["default"].self_link

  host_rule {
    hosts        = ["*"]
    path_matcher = "all-paths"
  }

  path_matcher {
    name            = "all-paths"
    default_service = google_compute_backend_bucket.website_backends["default"].self_link

    dynamic "path_rule" {
      for_each = var.api_gateways

      content {
        service = google_compute_backend_service.api_gateway_backends[path_rule.key].self_link
        paths   = ["/apis/${path_rule.key}/*"]
        route_action {
          url_rewrite {
            path_prefix_rewrite = "/"
          }
        }

      }
    }

    dynamic "path_rule" {
      for_each = tomap({
        for key, value in var.website_buckets :
        key => value if key != "default"
      })

      content {
        service = google_compute_backend_bucket.website_backends[path_rule.key].self_link
        paths   = [startswith(path_rule.value.base_path, "/") ? "${path_rule.value.base_path}/*" : "/${path_rule.value.base_path}/*"]
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
}

# Create DNS Records for the CDN
resource "google_dns_record_set" "cdn_dns_record" {
  name         = endswith(var.cdn_domain.domain_name, ".") ? var.cdn_domain.domain_name : "${var.cdn_domain.domain_name}."
  managed_zone = data.google_dns_managed_zone.cdn_zone.name
  type         = "A"
  rrdatas      = [google_compute_global_address.cdn_ip.address]
}

resource "google_dns_record_set" "www_cdn_dns_record" {
  name         = endswith(var.cdn_domain.domain_name, ".") ? "www.${var.cdn_domain.domain_name}" : "www.${var.cdn_domain.domain_name}."
  managed_zone = data.google_dns_managed_zone.cdn_zone.name
  type         = "A"
  rrdatas      = [google_compute_global_address.cdn_ip.address]
}

resource "google_certificate_manager_certificate" "cdn_cert" {
  provider    = google-beta
  name        = "${var.stack_id}-cdn-cert"
  description = "Nitric stack CDN SSL certificate"
  scope       = "DEFAULT"

  managed {
    domains = [var.cdn_domain.domain_name]
  }
}

resource "google_certificate_manager_certificate_map" "cdn_cert_map" {
  provider    = google-beta
  name        = "cert-map"
  description = "CDN Certificate Map"
}

resource "google_certificate_manager_certificate_map_entry" "cdn_cert_map_entry" {
  name        = "cdn-cert-map-entry"
  description = "CDN Certificate Map Entry"
  map = google_certificate_manager_certificate_map.cdn_cert_map.name 
  certificates = [google_certificate_manager_certificate.cdn_cert.id]
  matcher = "PRIMARY"
}

# Create a Target HTTPS Proxy
resource "google_compute_target_https_proxy" "https_proxy" {
  name             = "https-proxy"
  # ssl_certificates = [google_compute_managed_ssl_certificate.cdn_cert.self_link]
  certificate_map = "//certificatemanager.googleapis.com/${google_certificate_manager_certificate_map.cdn_cert_map.id}"
  url_map          = google_compute_url_map.https_url_map.self_link
}

# Create a Global Forwarding Rule
resource "google_compute_global_forwarding_rule" "https_forwarding_rule" {
  name        = "https-forwarding-rule"
  ip_address  = google_compute_global_address.cdn_ip.address
  ip_protocol = "TCP"
  port_range  = "443"
  target      = google_compute_target_https_proxy.https_proxy.self_link
}

# Invalidate the CDN cache if files have changed
resource "null_resource" "invalidate_cache" {
  provisioner "local-exec" {
    command = "gcloud compute url-maps invalidate-cdn-cache ${google_compute_url_map.https_url_map.name} --path '/*' --project ${var.project_id}"
  }

  triggers = {
    url_map = google_compute_url_map.https_url_map.self_link
  }

  depends_on = [ google_compute_global_forwarding_rule.https_forwarding_rule ]
}