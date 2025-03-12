terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

locals {
  endpoint_name = "${var.stack_name}-cdn"
  default_origin_group_name = "website-origin-group"
}

data "azurerm_client_config" "current" {}

# Create the CDN profile for the website
resource "azurerm_cdn_profile" "website_profile" {
  name = "website-profile"
  resource_group_name = var.resource_group_name
  location = var.location
  sku = "Standard_Microsoft"
}


# Create the CDN endpoint with full configuration using azapi, due to limitations in the azurerm provider
# This includes all origins, origin groups, and delivery rule, see https://github.com/hashicorp/terraform-provider-azurerm/issues/10771
resource "azapi_resource" "cdn_endpoint" {
  type        = "Microsoft.Cdn/profiles/endpoints@2024-09-01"
  name        = local.endpoint_name
  parent_id   = azurerm_cdn_profile.website_profile.id
  location    = var.location

  body = {
    properties = {
      isHttpAllowed          = false
      isHttpsAllowed         = true
      isCompressionEnabled   = true
      contentTypesToCompress = [
        "text/html", "text/css", "application/javascript",
        "application/json", "image/svg+xml", "font/woff", "font/woff2"
      ]
      
      # Define all origins in one place
      origins = concat(
        # Storage origin (duplicated from Terraform resource)
        [{
          name = var.storage_account_name
          properties = {
            hostName = var.storage_account_primary_web_host
            originHostHeader = var.storage_account_primary_web_host
            enabled = true
            httpPort = 80
            httpsPort = 443
          }
        }],
        # API origins
        [for key in sort(keys(var.apis)) : {
          name = "api-origin-${key}"
          properties = {
            hostName = replace(replace(var.apis[key].gateway_url, "https://", ""), "/", "")
            originHostHeader = replace(replace(var.apis[key].gateway_url, "https://", ""), "/", "")
            httpPort = 80
            httpsPort = 443
            enabled = true
          }
        }]
      )
      
      # Define all origin groups
      originGroups = concat(
        # Website origin group
        [{
          name = "website-origin-group"
          properties = {
            origins = [{
              id = "${azurerm_cdn_profile.website_profile.id}/endpoints/${local.endpoint_name}/origins/${var.storage_account_name}"
            }]
          }
        }],
        # API origin groups
        [for key in sort(keys(var.apis)) : {
          name = "api-origin-group-${key}"
          properties = {
            origins = [{
              id = "${azurerm_cdn_profile.website_profile.id}/endpoints/${local.endpoint_name}/origins/api-origin-${key}"
            }]
          }
        }]
      )
      
      # Set default origin group
      defaultOriginGroup = {
        id = "${azurerm_cdn_profile.website_profile.id}/endpoints/${local.endpoint_name}/originGroups/website-origin-group"
      }
      
      # Define delivery rules
      deliveryPolicy = {
        description = "API routing policy"
        rules = [
          for idx, key in { for i, k in sort(keys(var.apis)) : i => k } : {
            name  = "forward_${key}"
            order = idx + 1
            conditions = [{
              name = "UrlPath"
              parameters = {
                typeName = "DeliveryRuleUrlPathMatchConditionParameters"
                operator = "BeginsWith"
                matchValues = ["/api/${key}"]
                transforms = ["Lowercase"]
                negateCondition = false
              }
            }]
            actions = [
              {
                name = "OriginGroupOverride"
                parameters = {
                  typeName = "DeliveryRuleOriginGroupOverrideActionParameters"
                  originGroup = {
                    id = "${azurerm_cdn_profile.website_profile.id}/endpoints/${local.endpoint_name}/originGroups/api-origin-group-${key}"
                  }
                }
              },
              {
                name = "UrlRewrite"
                parameters = {
                  typeName = "DeliveryRuleUrlRewriteActionParameters"
                  sourcePattern = "/api/${key}/"
                  destination = "/"
                }
              }
            ]
          }
        ]
      }
    }
  }

  response_export_values = ["properties.hostName"]

  depends_on = [azurerm_cdn_profile.website_profile]
}

