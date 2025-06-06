variable "nitric" {
  type = object({
    name     = string
    stack_id = string
    # A map of path to origin
    origins = map(object({
      path = string
      type = string
      http_endpoint = string
      # TODO: Possibly add an identity to use to authenticate the service
    }))
  })
}
