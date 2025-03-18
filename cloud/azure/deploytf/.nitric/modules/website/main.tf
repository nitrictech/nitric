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

# Retrieve the list of blobs in the storage account with their md5 hashes converted to hex
data "external" "existing_base64_md5_to_hex" {
  for_each = local.transformed_files
  program = ["bash", "-c", "echo -n \"{\\\"output\\\":\\\"$(az storage blob show --connection-string $1 -c '$web' -n $2 --query properties.contentSettings.contentMd5 --output tsv | base64 -d | xxd -p -c 32)\\\"}\"", "bash", var.storage_account_connection_string, trimprefix(each.key, "/")]
}

locals {
  existing_md5_map = { for key, data in data.external.existing_base64_md5_to_hex : trimprefix(key, "/") => data.result["output"] }

  # Filter out files that have changed by comparing the md5 hashes
  changed_files = {
    for path, file in local.transformed_files : "/${trimprefix(path, "/")}" => file.digests.md5
    if lookup(local.existing_md5_map, trimprefix(path, "/"), null) != null 
    && file.digests.md5 != lookup(local.existing_md5_map, trimprefix(path, "/"), null)
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

  # Make sure we check for existing blobs first
  depends_on = [
    data.external.existing_base64_md5_to_hex,
  ]
}

