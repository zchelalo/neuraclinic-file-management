output "bucket_name" {
  description = "S3 bucket name used by the file-management service."
  value       = aws_s3_bucket.files.bucket
}

output "bucket_arn" {
  description = "S3 bucket ARN."
  value       = aws_s3_bucket.files.arn
}

output "region" {
  description = "AWS region where the bucket exists."
  value       = var.aws_region
}
