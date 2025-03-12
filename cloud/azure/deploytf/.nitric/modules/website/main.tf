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


resource "azurerm_storage_blob" "website_files" {
  for_each = local.transformed_files

  name                   = trimprefix(each.key, "/")

  storage_account_name   = var.storage_account_name
  storage_container_name = "$web"
  type                   = "Block"
  source                 = each.value.source_path
  content_type           = each.value.content_type

  # required to detect file changes in Terraform 
  content_md5            = each.value.digests.md5
}

