resource "random_string" "id" {
  length  = 8
  special = false
  upper   = false
}

resource "aws_resourcegroups_group" "group" {
  name = "${var.project_name}-${var.stack_name}-${random_string.id.result}"

  resource_query {
    query = <<JSON
{

    "ResourceTypeFilters":["AWS::AllSupported"],
	"TagFilters":[{"Key":"x-nitric-${random_string.id.result}-name"}]
}
JSON
  }
}
