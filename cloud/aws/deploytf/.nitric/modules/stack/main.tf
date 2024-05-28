resource "random_string" "id" {
  length  = 8
  special = false
  upper   = false
}

// Create a new AWS provider
provider "aws" {
  region = var.region
}

resource "aws_resourcegroups_group" "group" {
  name = "nitric-resource-group-${random_string.id.result}"

  resource_query {
    query = <<JSON
{

    "ResourceTypeFilters":["AWS::AllSupported"],
	"TagFilters":[{"Key":"x-nitric-name-${random_string.id.result}"}]
}
JSON
  }
}