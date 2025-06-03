
# AWS S3 bucket
resource "aws_s3_bucket" "bucket" {
  bucket = "${var.nitric.stack_id}-${var.nitric.name}"
  tags = var.tags
}
