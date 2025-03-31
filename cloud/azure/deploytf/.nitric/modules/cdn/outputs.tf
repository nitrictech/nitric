output "cdn_url" {
  description = "The URL of the endpoint"
  value = "https://${azurerm_cdn_frontdoor_endpoint.cdn_endpoint.host_name}"
}

output "cdn_frontdoor_profile_id" {
  description = "The ID of the CDN profile"
  value = azurerm_cdn_frontdoor_profile.cdn_profile.id
}

output "cdn_frontdoor_api_rule_set_id" {
  description = "The ID of the API rewrite rule set"
  value = one(azurerm_cdn_frontdoor_rule_set.api_ruleset) != null ? one(azurerm_cdn_frontdoor_rule_set.api_ruleset).id : null
}

output "cdn_frontdoor_default_rule_set_id" {
  description = "The ID of the default rule set"
  value = azurerm_cdn_frontdoor_rule_set.default_ruleset.id
}