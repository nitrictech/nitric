output "stack_id" {
  description = "The randomized Id of the nitric stack"
  value       =  random_string.id.result
}