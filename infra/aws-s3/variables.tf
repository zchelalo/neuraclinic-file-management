variable "aws_region" {
  description = "AWS region where the file storage bucket will be created."
  type        = string
  default     = "us-east-1"
}

variable "project" {
  description = "Project name used for tags."
  type        = string
  default     = "neuraclinic"
}

variable "environment" {
  description = "Deployment environment used for tags."
  type        = string
  default     = "dev"
}

variable "bucket_name" {
  description = "Globally unique S3 bucket name."
  type        = string
}

variable "force_destroy" {
  description = "Allow Terraform to delete the bucket even when objects exist."
  type        = bool
  default     = false
}

variable "cors_allowed_origins" {
  description = "Origins allowed to use presigned upload and download URLs."
  type        = list(string)
  default     = ["*"]
}

variable "tags" {
  description = "Additional tags for all resources."
  type        = map(string)
  default     = {}
}
