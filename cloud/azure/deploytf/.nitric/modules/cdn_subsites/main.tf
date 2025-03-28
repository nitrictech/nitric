locals {
  subsite_origin_group_name = "${var.stack_name}-${var.name}-origin-group"
  subsite_origin_name       = "${var.stack_name}-${var.name}-origin"
  subsite_rule_name         = "subsiterule${var.name}"
}

# create an origin group for each subsite
resource "azurerm_cdn_frontdoor_origin_group" "subsite_origin_group" {  
  name                      = local.subsite_origin_group_name
  cdn_frontdoor_profile_id  = var.cdn_frontdoor_profile_id

  # https://learn.microsoft.com/en-us/azure/frontdoor/origin?pivots=front-door-standard-premium#load-balancing-settings
  load_balancing {
    additional_latency_in_milliseconds = 200 # Lower latency tolerance for faster failover
    sample_size                        = 5 # More samples for better decision-making
    successful_samples_required        = 3 # Keep at 3 to maintain reliability
  }
}

# create an origin for each subsite
resource "azurerm_cdn_frontdoor_origin" "subsite_origin" {
  name                           = local.subsite_origin_name
  cdn_frontdoor_origin_group_id  = azurerm_cdn_frontdoor_origin_group.subsite_origin_group.id
  enabled                        = true

  certificate_name_check_enabled = false

  host_name          = replace(var.primary_web_host, "https://", "")
  http_port          = 80
  https_port         = 443

  origin_host_header = replace(var.primary_web_host, "https://", "")

  depends_on = [azurerm_cdn_frontdoor_origin_group.subsite_origin_group]
}

# create a rule for each subsite
resource "azurerm_cdn_frontdoor_rule" "subsite_rule" {
  name                      = local.subsite_rule_name
  cdn_frontdoor_rule_set_id = var.cdn_default_frontdoor_rule_set_id
  order                     = var.rule_order

  conditions {
    url_path_condition {
      operator         = "RegEx"
      negate_condition = false
      match_values     = ["${trimprefix(var.base_path, "/")}(/.*)?$"]
      transforms       = ["Lowercase"]
    }
  }

  actions {
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.subsite_origin_group.id
      forwarding_protocol = "HttpsOnly"
      cache_behavior = "HonorOrigin"
      query_string_caching_behavior = "UseQueryString"
      compression_enabled = true
    }
    
    url_rewrite_action {
      source_pattern   = var.base_path
      destination      = "/"
      preserve_unmatched_path = true
    }
  }
 
  depends_on = [ 
    azurerm_cdn_frontdoor_origin_group.subsite_origin_group,
    azurerm_cdn_frontdoor_origin.subsite_origin
  ]
}