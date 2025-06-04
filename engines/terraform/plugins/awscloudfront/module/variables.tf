variable "nitric" {
  type = object({
    name     = string
    stack_id = string
    services = map(object({
      http_endpoint = string
      identity = map(object({
        id   = string
        role = any
      }))
    }))
  })
}
