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
resource "azurerm_cdn_profile" "website_profile" {
  name                = "${var.stack_name}-cdn-profile"
  resource_group_name = var.resource_group_name
  location            = var.location
  sku                 = "Standard_Microsoft"
}

# Create API Management instance in consumption tier only if there are APIs
# We are using an api for proxying to the backend APIs
# This is a workaround for the current limitation of Azure CDN with circular dependencies
resource "azurerm_api_management" "apim_cdn" {
  count               = local.has_apis ? 1 : 0
  
  name                = "${var.stack_name}-apim-cdn"
  location            = var.location
  resource_group_name = var.resource_group_name
  publisher_name      = var.publisher_name
  publisher_email     = var.publisher_email
  
  # Consumption tier for cost optimization
  sku_name            = "Consumption_0"
}

# Create a single API for proxying
resource "azurerm_api_management_api" "proxy_api" {
  count               = local.has_apis ? 1 : 0
  
  name                = "proxy"
  resource_group_name = var.resource_group_name
  api_management_name = azurerm_api_management.apim_cdn[0].name
  revision            = "1"
  display_name        = "API Proxy"
  protocols           = ["https"]
  
  subscription_required = false
  
  import {
    content_format = "openapi+json"
    content_value  = jsonencode({
      "openapi": "3.0.1",
      "info": {
        "title": "API Proxy",
        "version": "1.0"
      },
      "paths": {
        "/*": {
          "get": { "responses": { "200": { "description": "Successful response" } } },
          "post": { "responses": { "200": { "description": "Successful response" } } },
          "put": { "responses": { "200": { "description": "Successful response" } } },
          "delete": { "responses": { "200": { "description": "Successful response" } } },
          "patch": { "responses": { "200": { "description": "Successful response" } } },
          "options": { "responses": { "200": { "description": "Successful response" } } }
        }
      }
    })
  }
}

# Policy at API level for routing - only if there are APIs
resource "azurerm_api_management_api_policy" "proxy_policy" {
  count               = local.has_apis ? 1 : 0
  
  api_management_name = azurerm_api_management.apim_cdn[0].name
  resource_group_name = var.resource_group_name
  api_name            = azurerm_api_management_api.proxy_api[0].name
  
  xml_content = <<XML
<policies>
  <inbound>
    <base />
    <choose>
      ${join("", [for api_key in local.api_keys : <<-EOT
        <when condition="@(context.Request.Url.Path.StartsWith("api/${api_key}"))" >
          <set-backend-service base-url="${var.apis[api_key].gateway_url}" />
          <rewrite-uri template="@(context.Request.Url.Path.Substring("api/${api_key}/".Length))" />
        </when>
      EOT
      ])}
      <otherwise>
        <return-response>
          <set-status code="404" reason="Not Found" />
          <set-header name="Content-Type" exists-action="override">
          <value>application/json</value>
          </set-header>
          <set-body>{"error": "API not found"}</set-body> 
        </return-response>
      </otherwise>
    </choose>
  </inbound>
  <backend>
    <base />
  </backend>
  <outbound>
    <base />
  </outbound>
  <on-error>
    <base />
  </on-error>
</policies>
XML
}

# Create the CDN endpoint with conditional API configuration
resource "azapi_resource" "cdn_endpoint" {
  type        = "Microsoft.Cdn/profiles/endpoints@2024-02-01"
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
      
      # Define origins
      origins = concat(
        [
          # Website origin
          {
            name = var.storage_account_name
            properties = {
              hostName = var.storage_account_primary_web_host
              originHostHeader = var.storage_account_primary_web_host
              enabled = true
              httpPort = 80
              httpsPort = 443
            }
          }
        ],
        # Conditionally add API origin
        local.has_apis ? [
          {
            name = local.api_proxy_origin_name
            properties = {
              hostName = replace(azurerm_api_management.apim_cdn[0].gateway_url, "https://", "")
              originHostHeader = replace(azurerm_api_management.apim_cdn[0].gateway_url, "https://", "")
              enabled = true
              httpPort = 80
              httpsPort = 443
            }
          }
        ] : []
      )
      
      # Define origin groups
      originGroups = concat(
        [
          # Website origin group
          {
            name = local.default_origin_group_name
            properties = {
              origins = [{
                id = "${azurerm_cdn_profile.website_profile.id}/endpoints/${local.endpoint_name}/origins/${var.storage_account_name}"
              }]
            }
          }
        ],
        # Conditionally add API origin group
        local.has_apis ? [
          {
            name = local.api_proxy_origin_group_name
            properties = {
              origins = [{
                id = "${azurerm_cdn_profile.website_profile.id}/endpoints/${local.endpoint_name}/origins/${local.api_proxy_origin_name}"
              }]
            }
          }
        ] : []
      )
      
      # Set default origin group
      defaultOriginGroup = {
        id = "${azurerm_cdn_profile.website_profile.id}/endpoints/${local.endpoint_name}/originGroups/${local.default_origin_group_name}"
      }
      
      # Define delivery rules
      deliveryPolicy = {
        description = "Content delivery rules"
        rules = concat(
          # Conditionally add API routing rule
          local.has_apis ? [
            {
              name = "RouteAPIRequests"
              order = 1
              conditions = [
                {
                  name = "UrlPath"
                  parameters = {
                    typeName = "DeliveryRuleUrlPathMatchConditionParameters"
                    operator = "BeginsWith"
                    matchValues = ["/api/"]
                  }
                }
              ],
              actions = [
                {
                  name = "OriginGroupOverride"
                  parameters = {
                    typeName = "DeliveryRuleOriginGroupOverrideActionParameters"
                    originGroup = {
                      id = "${azurerm_cdn_profile.website_profile.id}/endpoints/${local.endpoint_name}/originGroups/${local.api_proxy_origin_group_name}"
                    }
                  }
                }
              ]
            }
          ] : []
        )
      }
    }
  }

  response_export_values = ["properties.hostName"]

  depends_on = [azurerm_cdn_profile.website_profile]
}

resource "terraform_data" "cdn_purge" {  
  # This will run on every Terraform apply
  triggers_replace = [
    # Force this to run on every apply when paths are changed
    local.changed_path_md5_hashes
  ]
  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = <<EOF
      # Purge the CDN endpoint if local.transformed_paths isn't empty
      if [ ${length(local.transformed_paths)} -gt 0 ]; then
      MSYS_NO_PATHCONV=1 az cdn endpoint purge \
        --resource-group ${var.resource_group_name} \
        --profile-name ${azurerm_cdn_profile.website_profile.name} \
        --name ${local.endpoint_name} \
        --content-paths ${join(" ", formatlist("\"%s\"", local.transformed_paths))}
      fi
    EOF
  }

  depends_on = [azapi_resource.cdn_endpoint]
}