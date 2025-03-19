terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

locals {
  endpoint_name               = "${var.stack_name}-cdn"
  default_origin_group_name   = "website-origin-group"
  default_origin_name         = "website-origin"
  api_proxy_origin_group_name = "api-proxy-origin-group"
  api_proxy_origin_name       = "api-proxy-origin"

  # sorted api keys
  api_keys = sort(keys(var.apis))
  
  # Boolean flag to indicate if we have any APIs
  has_apis = length(local.api_keys) > 0

  # ensure any path that ends with index.html is treated as a directory
  transformed_paths = sort([
    for path in keys(var.cdn_purge_paths) : (
      endswith(path, "/index.html") ? trimsuffix(path, "index.html") : path
    )
  ])

  changed_path_md5_hashes = join(",", sort(values(var.cdn_purge_paths)))
}

# Create the CDN profile for the website
resource "azurerm_cdn_frontdoor_profile" "cdn_profile" {
  name                = "${var.stack_name}-cdn-profile"
  resource_group_name = var.resource_group_name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_endpoint" "cdn_endpoint" {
  name                     = local.endpoint_name
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.cdn_profile.id
}

resource "azurerm_cdn_frontdoor_origin_group" "default_origin_group" {
  name                     = local.default_origin_group_name
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.cdn_profile.id

  load_balancing {
    additional_latency_in_milliseconds = 200
    sample_size                        = 5
    successful_samples_required        = 3
  }
}


resource "azurerm_cdn_frontdoor_origin" "default_origin" {
  name                          = local.default_origin_name
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.default_origin_group.id

  certificate_name_check_enabled = false

  host_name          = var.storage_account_primary_web_host
  http_port          = 80
  https_port         = 443
  origin_host_header = var.storage_account_primary_web_host

  depends_on = [ azurerm_cdn_frontdoor_origin_group.default_origin_group ]
}

# create default ruleset if apis are present
resource "azurerm_cdn_frontdoor_rule_set" "default_ruleset" {
  count                     = local.has_apis ? 1 : 0
  name                      = "apiruleset"
  cdn_frontdoor_profile_id  = azurerm_cdn_frontdoor_profile.cdn_profile.id
}

# create an origin group for each api
resource "azurerm_cdn_frontdoor_origin_group" "api_origin_group" {
  for_each = { for idx, key in local.api_keys : key => var.apis[key] }
  
  name                      = "${local.api_proxy_origin_group_name}-${each.key}"
  cdn_frontdoor_profile_id  = azurerm_cdn_frontdoor_profile.cdn_profile.id
  load_balancing {
    additional_latency_in_milliseconds = 100 # Reduced latency for API
    sample_size                        = 5 # Increased sample size for better accuracy
    successful_samples_required        = 2 # Reduced successful samples required for faster failover
  }
}

# create an origin for each api
resource "azurerm_cdn_frontdoor_origin" "api_origin" {
  for_each = { for idx, key in local.api_keys : key => var.apis[key] }

  name                           = "${local.api_proxy_origin_name}-${each.key}"
  cdn_frontdoor_origin_group_id  = azurerm_cdn_frontdoor_origin_group.api_origin_group[each.key].id

  certificate_name_check_enabled = false

  host_name          = replace(each.value.gateway_url, "https://", "")
  http_port          = 80
  https_port         = 443
  origin_host_header = replace(each.value.gateway_url, "https://", "")
}

# create a rule for each api
resource "azurerm_cdn_frontdoor_rule" "api_rule" {
  for_each                  = { for idx, key in local.api_keys : key => var.apis[key] }
  name                      = "apirule${each.key}"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.default_ruleset[0].id
  order                     = index(local.api_keys, each.key) + 1

  
  conditions {
    url_path_condition {
      operator         = "BeginsWith"
      negate_condition = false
      match_values     = ["/api/${each.key}"]
      transforms       = ["Lowercase"]
    }
  }

  actions {
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.api_origin_group[each.key].id
      forwarding_protocol = "HttpsOnly"
      cache_behavior = "HonorOrigin"
      query_string_caching_behavior = "IgnoreQueryString"
    }
    
    url_rewrite_action {
      source_pattern   = "/api/${each.key}/"
      destination      = "/"
      preserve_unmatched_path = false
    }
  }
 
  depends_on = [ 
    azurerm_cdn_frontdoor_rule_set.default_ruleset,
    azurerm_cdn_frontdoor_origin.api_origin,
    azurerm_cdn_frontdoor_origin_group.api_origin_group
  ]
}

# Create the CDN route
resource "azurerm_cdn_frontdoor_route" "main_route" {
  name                       = "${var.stack_name}-main-route"
  cdn_frontdoor_endpoint_id =  azurerm_cdn_frontdoor_endpoint.cdn_endpoint.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.default_origin_group.id
  cdn_frontdoor_origin_ids = [
    azurerm_cdn_frontdoor_origin.default_origin.id
  ]
  cdn_frontdoor_rule_set_ids = local.has_apis ? [azurerm_cdn_frontdoor_rule_set.default_ruleset[0].id] : []
  supported_protocols = [ "Https" ]
  https_redirect_enabled = false
  patterns_to_match =  [
    "/*"
  ]

  cache {
    query_string_caching_behavior = "IgnoreQueryString"
    compression_enabled = true
    content_types_to_compress     = [
      "text/html", "text/css", "application/javascript",
      "application/json", "image/svg+xml"
    ]
  }

  depends_on = [ 
    azurerm_cdn_frontdoor_origin.default_origin,
    azurerm_cdn_frontdoor_origin_group.default_origin_group
  ]
}

data "azurerm_subscription" "current" {}

resource "terraform_data" "endpoint_purge" {  
  # This will run on every Terraform apply
  triggers_replace = [
    # Force this to run on every apply when paths are changed
    local.changed_path_md5_hashes
  ]
  
  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = <<EOF
      # Purge the endpoint if local.transformed_paths isn't empty
      if [ ${length(local.transformed_paths)} -gt 0 ]; then
      MSYS_NO_PATHCONV=1 az afd endpoint purge \
        --resource-group ${var.resource_group_name} \
        --profile-name ${sensitive(azurerm_cdn_frontdoor_profile.cdn_profile.name)} \
        --subscription ${sensitive(data.azurerm_subscription.current.subscription_id)} \
        --endpoint-name ${sensitive(local.endpoint_name)} \
        --content-paths ${join(" ", formatlist("\"%s\"", local.transformed_paths))} \
        --no-wait
      fi
    EOF
  }

  depends_on = [azurerm_cdn_frontdoor_endpoint.cdn_endpoint]
}