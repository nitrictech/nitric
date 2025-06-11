variable "nitric" {
  type = object({
    name     = string
    stack_id = string
    services = map(object({
      actions = list(string)
      identities = map(object({
        id   = string
        role = any
      }))
    }))
  })
}

variable "existing" {
    type = object({
        project_id = string
        branch_id = optional(string, null)
        role_name = optional(string, null)
        database_name = optional(string, true)
        endpoint_id = optional(string, null)
    })
    default = null
    nullable = true
}
