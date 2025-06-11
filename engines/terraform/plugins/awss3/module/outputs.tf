output "nitric" {
    value = {
        id = aws_s3_bucket.bucket.arn
        domain_name = aws_s3_bucket.bucket.bucket_regional_domain_name
        exports = {
            # Export env variables to be mapped to all services that access this resource
            # TODO: May need a per service mapping as well
            env = {}
            # Export resources to be read into other modules
            resources = {
                "aws_s3_bucket" = aws_s3_bucket.bucket.arn
            }
        }
    }
}