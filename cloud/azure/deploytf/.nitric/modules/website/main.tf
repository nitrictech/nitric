module "template_files" {
  source  = "hashicorp/dir/template"
  version = "1.0.2"
  
  base_dir = var.local_directory
}

locals {
  normalized_base_path = var.base_path == "/" ? "root" : replace(var.base_path, "/", "")

  uploaded_files_md5 = {
    for path, file in module.template_files.files : (
      var.base_path == "/" ? 
        path : 
        "${trimsuffix(var.base_path, "/")}/${path}"
    ) => file.digests.md5
  }
}


resource "azurerm_storage_account" "storage_website_account" {
  name                = replace("${var.stack_name}st${local.normalized_base_path}", "-", "")
  resource_group_name = var.resource_group_name
  location            = var.location
  account_tier        = "Standard"
  access_tier         = "Hot"
 
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
}

# Create a storage container for the website
resource "azurerm_storage_account_static_website" "storage_website" {
  storage_account_id = azurerm_storage_account.storage_website_account.id
  index_document     = var.index_document
  error_404_document = var.error_document

  depends_on = [azurerm_storage_account.storage_website_account]
}

resource "azurerm_storage_blob" "website_files" {
  for_each = module.template_files.files

  name                   = trimprefix(each.key, "/")

  storage_account_name   = azurerm_storage_account.storage_website_account.name
  storage_container_name = "$web"
  type                   = "Block"
  source                 = each.value.source_path
  content_type           = each.value.content_type

  # required to detect file changes in Terraform 
  content_md5            = each.value.digests.md5

  depends_on = [ 
    azurerm_storage_account_static_website.storage_website, 
    azurerm_storage_account.storage_website_account 
  ]
}

