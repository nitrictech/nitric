
output "changed_files" {
  description = "Map of changed paths that need purging"
  value = local.changed_files
}