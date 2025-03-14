resource "random_string" "id" {
  length  = 8
  special = false
  upper   = false
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

# AWS S3 bucket
resource "aws_s3_bucket" "bucket" {
  bucket = "website-bucket-${random_string.id.result}"

  tags = {
    "x-nitric-${var.stack_id}-name" = "website-bucket-${random_string.id.result}"
    "x-nitric-${var.stack_id}-type" = "bucket"
  }
}
