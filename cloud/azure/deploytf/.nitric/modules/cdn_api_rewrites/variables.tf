variable "name" {
  description = "The name of the api"
  type        = string
}

variable "api_host_name" {
  description = "The host name of the api"
  type        = string
}

variable "cdn_frontdoor_profile_id" {
  description = "The id of the cdn frontdoor profile to use for the cdn"
  type        = string
}

variable "cdn_frontdoor_rule_set_id" {
  description = "The id of the default cdn frontdoor rule set to use for the cdn"
  type        = string
}

variable "rule_order" {
  description = "The order of the rule to use for the cdn"
  type        = number
  default     = 1
}