output "bucket_arn" {
  description = "The ARN of the deployed bucket"
  value       =  aws_s3_bucket.bucket.arn
}
