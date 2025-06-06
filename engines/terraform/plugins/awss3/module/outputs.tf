output "nitric" {
    value = {
        id = aws_s3_bucket.bucket.arn
        domain_name = aws_s3_bucket.bucket.bucket_regional_domain_name
        raw_type = "aws_s3_bucket"
        raw = aws_s3_bucket.bucket
    }
}