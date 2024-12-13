output "stack_id" {
  description = "The randomized Id of the nitric stack"
  value       =  random_string.id.result
}

output "resource_group_arn" {
  description = "The resource group for the stack"
  value       = aws_resourcegroups_group.group.arn
}
