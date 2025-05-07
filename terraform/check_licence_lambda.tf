
resource "aws_lambda_function" "check_license" {
  function_name = "check_license"
  description   = "Handler for checking pings from client apps to check for other machines using the same license"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = "${path.module}/lambda/function.zip"
  source_code_hash = filebase64sha256("${path.module}/lambda/function.zip")

  environment {
    variables = {
      DB_NAME      = "mst_db"
      DB_HOST      = var.db_endpoint
      DB_URL_WRITE = "mst-aurora-db.cluster-cvq42ycqkt4f.us-west-1.rds.amazonaws.com"
      DB_URL_READ  = "mst-aurora-db.cluster-ro-cvq42ycqkt4f.us-west-1.rds.amazonaws.com"
      DB_USER      = "mst_admin"
      DB_PORT      = 5432
      DB-PASSWORD  = var.db_password
    }
  }

  role = aws_iam_role.lambda_role.arn
}

# Create IAM role for the Lambda function
resource "aws_iam_role" "lambda_role" {
  name = "lambda_check_license_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}
