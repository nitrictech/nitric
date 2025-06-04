variable "nitric" {
  type = object({
    name     = string
    stack_id = string
    services = list(object({
      name    = string
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

variable "read_role" {

}
