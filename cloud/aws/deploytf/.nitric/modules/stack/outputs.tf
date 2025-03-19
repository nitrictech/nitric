output "stack_id" {
  description = "The randomized Id of the nitric stack"
  value       =  random_string.id.result
}

output "website_bucket_name" {
  value = aws_s3_bucket.bucket.bucket
  description = "The name of the bucket"
}

output "website_bucket_id" {
  description = "The ID for the website bucket"
  value = aws_s3_bucket.bucket.id
}

output "website_bucket_arn" {
  description = "The ARN for the website bucket"
  value = aws_s3_bucket.bucket.arn
}

output "website_bucket_domain_name" {
  description = "The domain name for the website bucket"
  value = aws_s3_bucket.bucket.bucket_regional_domain_name
}