locals {
  endpoint_name               = "${var.stack_name}-cdn"
  default_origin_group_name   = "${var.stack_name}-default-origin-group"
  default_origin_name         = "${var.stack_name}-default-origin"

  changed_path_md5_hashes     = join("", sort(values(var.uploaded_files)))
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
  enabled            = true

  host_name          = var.primary_web_host
  http_port          = 80
  https_port         = 443
  origin_host_header = var.primary_web_host

  depends_on = [ azurerm_cdn_frontdoor_origin_group.default_origin_group ]
}

# create default ruleset for default rules
resource "azurerm_cdn_frontdoor_rule_set" "default_ruleset" {
  name                      = "defaultruleset"
  cdn_frontdoor_profile_id  = azurerm_cdn_frontdoor_profile.cdn_profile.id
}

resource "azurerm_cdn_frontdoor_rule" "redirect_slash_rule" {
  name                      = "redirectslash"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.default_ruleset.id
  order                     = 1

  
  conditions {
    url_path_condition {
      operator         = "RegEx"
      match_values     = [".*\\/$"]
    }
  }

  actions {
    url_redirect_action {
      redirect_type        = "Found"
      redirect_protocol    = "MatchRequest"
      destination_path     = "/{url_path:0:-1}"
      destination_hostname = ""
    }
  }
 
  depends_on = [ azurerm_cdn_frontdoor_rule_set.default_ruleset ]
}

# create default ruleset if apis are present
resource "azurerm_cdn_frontdoor_rule_set" "api_ruleset" {
  count                     = var.enable_api_rewrites ? 1 : 0
  name                      = "apiruleset"
  cdn_frontdoor_profile_id  = azurerm_cdn_frontdoor_profile.cdn_profile.id
}

# Create the CDN route
resource "azurerm_cdn_frontdoor_route" "main_route" {
  name                       = "${var.stack_name}-main-route"
  cdn_frontdoor_endpoint_id =  azurerm_cdn_frontdoor_endpoint.cdn_endpoint.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.default_origin_group.id
  cdn_frontdoor_origin_ids = [
    azurerm_cdn_frontdoor_origin.default_origin.id
  ]
  cdn_frontdoor_rule_set_ids = var.enable_api_rewrites ? [azurerm_cdn_frontdoor_rule_set.default_ruleset.id, azurerm_cdn_frontdoor_rule_set.api_ruleset[0].id] : [azurerm_cdn_frontdoor_rule_set.default_ruleset.id]
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
    azurerm_cdn_frontdoor_profile.cdn_profile,
    azurerm_cdn_frontdoor_endpoint.cdn_endpoint,
    azurerm_cdn_frontdoor_origin_group.default_origin_group,
    azurerm_cdn_frontdoor_origin.default_origin,
    azurerm_cdn_frontdoor_rule_set.default_ruleset,
    azurerm_cdn_frontdoor_rule_set.api_ruleset
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
      MSYS_NO_PATHCONV=1 az afd endpoint purge \
        --resource-group ${var.resource_group_name} \
        --profile-name ${sensitive(azurerm_cdn_frontdoor_profile.cdn_profile.name)} \
        --subscription ${sensitive(data.azurerm_subscription.current.subscription_id)} \
        --endpoint-name ${sensitive(local.endpoint_name)} \
        --content-paths '/*' \
        --no-wait
    EOF
  }

  depends_on = [azurerm_cdn_frontdoor_endpoint.cdn_endpoint]
}