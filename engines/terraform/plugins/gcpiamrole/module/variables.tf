variable "nitric" {
  type = object({
    name = string
  })
}

variable "trusted_services" {
  type    = list(string)
  default = ["cloudrun.googleapis.com"]
}

variable "trusted_actions" {
  type    = list(string)
  default = [
    "monitoring.timeSeries.create",
    "resourcemanager.projects.get",
  ]
}

variable "project_id" {
  type = string
}