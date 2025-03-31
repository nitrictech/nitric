
output "uploaded_files" {
  description = "Map of uploaded files with their MD5 hashes"
  value = local.uploaded_files_md5
}

output "storage_account_web_host" {
  value = azurerm_storage_account.storage_website_account.primary_web_host
}