module "template_files" {
  source  = "hashicorp/dir/template"
  version = "1.0.2"
  
  base_dir = var.local_directory
}

# AWS S3 bucket
resource "aws_s3_bucket" "website_bucket" {
  bucket = "website-bucket-${var.name}-${var.stack_id}"
}

resource "aws_s3_object" "object" {
  for_each = module.template_files.files
  bucket = aws_s3_bucket.website_bucket.id

  key    = each.key
  source = each.value.source_path
  content_type = each.value.content_type

  # required to detect file changes in Terraform 
  etag = filemd5(each.value.source_path)
}
