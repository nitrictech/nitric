variable "principals" {
    description = "principals (roles) to apply the policies to"
    type = map(string)
}

variable "actions" {
    description = "actions to allow"
    type = set(string)
}

variable "resources" {
    description = "resources to apply the policies to"
    type = set(string)
}