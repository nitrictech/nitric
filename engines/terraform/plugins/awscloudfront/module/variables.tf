variable "nitric" {
  type = object({
    name     = string
    stack_id = string
    # A map of path to origin
    origins = map(object({
      path = string
      type = string
      domain_name = string
      id = string
      raw = map(any)
    }))
  })
}
