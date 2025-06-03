variable "nitric" {
    type = object({
        name = string
        stack_id = string
    })
}

variable "tags" {
    type = map(string)
    default = {}
}