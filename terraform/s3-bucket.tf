module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~>4.7.0"

  bucket = "mst-binaries-bucket"
  acl    = "private"

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  versioning = {
    enabled = true
  }
}

