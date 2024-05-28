
# Deploy an aws dynamodb table
resource "aws_dynamodb_table" "table" {
  name         = var.kvstore_name
  attribute {
    name = "_pk"
    type = "S"
  }
  attribute {
    name = "_sk"
    type = "S"
  }
  hash_key  = "_pk"
  range_key = "_sk"
  billing_mode = "PAY_PER_REQUEST"
  tags = {
    "x-nitric-${var.stack_id}-name" = var.kvstore_name
    "x-nitric-${var.stack_id}-type" = "kvstore"
  }
}
