
# Generate a random id for the bucket
resource "random_id" "bucket_id" {
  byte_length = 8

  keepers = {
    # Generate a new id each time we switch to a new AMI id
    bucket_name = var.bucket_name
  }
}

# AWS S3 bucket
resource "aws_s3_bucket" "bucket" {
  bucket = "${var.bucket_name}-${random_id.bucket_id.hex}"
  acl    = "private"

  tags = {
    "x-nitric-${var.stack_id}-name" = var.bucket_name
    "x-nitric-${var.stack_id}-type" = "bucket"
  }
}