variable "nitric" {
  type = object({
    name     = string
    stack_id = string
    content_path = string
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
