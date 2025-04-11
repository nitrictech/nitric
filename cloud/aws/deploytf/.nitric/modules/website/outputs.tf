output "changed_files" {
  description = "a list of all changed files"
  value = [for obj in aws_s3_object.object : obj.etag]
}

output "website_bucket_domain" {
  description = "the domain name of the website bucket"
  value = aws_s3_bucket.website_bucket.bucket_regional_domain_name
}

output "website_arn" {
  description = "the arn of the website bucket"
  value = aws_s3_bucket.website_bucket.arn
}

output "website_id" {
  description = "the ID of the website bucket"
  value = aws_s3_bucket.website_bucket.id
}