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
  default     = "postgres://mst_admin:Q70AqiE8KOfRIHxmqmN4@tf-20250502141116491500000001.cvq42ycqkt4f.us-west-1.rds.amazonaws.com:5432/mst_db"
}

variable "mst_website_github_token" {
  description = "Github token for mst-website"
  type        = string
  sensitive   = true
}
