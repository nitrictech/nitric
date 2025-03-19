module "template_files" {
  source  = "hashicorp/dir/template"
  version = "1.0.2"
  
  base_dir = var.local_directory
}

locals {
  # Apply the base path logic for key transformation
  transformed_files = {
    for path, file in module.template_files.files : (
      var.base_path == "/" ? 
        path : 
        "${trimsuffix(var.base_path, "/")}/${path}"
    ) => file
  }
}

resource "aws_s3_object" "object" {
  for_each = local.transformed_files
  bucket = var.website_bucket_id

  key    = trimprefix(each.key, "/")
  source = each.value.source_path
  content_type = each.value.content_type

  # required to detect file changes in Terraform 
  etag = each.value.digests.md5
}
