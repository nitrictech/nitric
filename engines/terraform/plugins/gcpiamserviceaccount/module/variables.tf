variable "nitric" {
  type = object({
    name = string
    stack_id = string
  })
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