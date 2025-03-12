# TODO purge cdn with changed files
# output "uploaded_files" {
#   value = {
#     for file, blob in azurerm_storage_blob.website_files :
#     file => blob.url
#     if blob.content_md5 != filemd5("${var.local_directory}/${file}") # Filter changed files
#   }
# }