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

variable "tags" {
  type    = map(string)
  default = {}
}
