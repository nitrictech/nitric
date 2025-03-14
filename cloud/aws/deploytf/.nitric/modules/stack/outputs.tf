output "stack_id" {
  description = "The randomized Id of the nitric stack"
  value       =  random_string.id.result
}

output "website_bucket_name" {
  value = aws_s3_bucket.bucket.bucket
  description = "The name of the bucket"
}
