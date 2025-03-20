locals {
  api_origin_group_name = "api-origin-group-${var.name}"
  api_origin_name       = "api-origin-${var.name}"
  api_rule_name         = "apirule${var.name}"
}

# create an origin group for each api
resource "azurerm_cdn_frontdoor_origin_group" "api_origin_group" {  
  name                      = local.api_origin_group_name
  cdn_frontdoor_profile_id  = var.cdn_frontdoor_profile_id
  load_balancing {
    additional_latency_in_milliseconds = 100 # Reduced latency for API
    sample_size                        = 5 # Increased sample size for better accuracy
    successful_samples_required        = 2 # Reduced successful samples required for faster failover
  }
}

# create an origin for each api
resource "azurerm_cdn_frontdoor_origin" "api_origin" {
  name                           = local.api_origin_name
  cdn_frontdoor_origin_group_id  = azurerm_cdn_frontdoor_origin_group.api_origin_group.id
  enabled                        = true

  certificate_name_check_enabled = false

  host_name          = replace(var.api_host_name, "https://", "")
  http_port          = 80
  https_port         = 443
  origin_host_header = replace(var.api_host_name, "https://", "")

  depends_on = [azurerm_cdn_frontdoor_origin_group.api_origin_group]
}

# create a rule for each api
resource "azurerm_cdn_frontdoor_rule" "api_rule" {
  name                      = local.api_rule_name
  cdn_frontdoor_rule_set_id = var.cdn_frontdoor_rule_set_id
  order                     = var.rule_order

  
  conditions {
    url_path_condition {
      operator         = "BeginsWith"
      negate_condition = false
      match_values     = ["/api/${var.name}"]
      transforms       = ["Lowercase"]
    }
  }

  actions {
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.api_origin_group.id
      forwarding_protocol = "HttpsOnly"
      cache_behavior = "HonorOrigin"
      query_string_caching_behavior = "UseQueryString"
      compression_enabled = true
    }
    
    url_rewrite_action {
      source_pattern   = "/api/${var.name}/"
      destination      = "/"
      preserve_unmatched_path = true
    }
  }
 
  depends_on = [ 
    azurerm_cdn_frontdoor_origin_group.api_origin_group,
    azurerm_cdn_frontdoor_origin.api_origin
  ]
}