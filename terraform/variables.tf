variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-1"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

