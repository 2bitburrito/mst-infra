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

variable "db_endpoint" {
  description = "Url endpoint for db"
  type        = string
  default     = "mst-aurora-db.cluster-ro-cvq42ycqkt4f.us-west-1.rds.amazonaws.com"
}
