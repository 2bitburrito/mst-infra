resource "aws_lambda_function" "backend_api" {
  function_name = "backend_api"
  description   = "Main entrypoint for restful(ish) Api service for MST"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  # TODO: Change this to the correct name after writing the lambda:
  filename         = "${path.module}/lambda/function.zip"
  source_code_hash = filebase64sha256("${path.module}/lambda/function.zip")

  vpc_config {
    subnet_ids         = var.mst_db_vpc_subnets
    security_group_ids = [aws_security_group.mst_db_sg.id]
  }

  environment {
    variables = {
      DB_NAME      = "mst_db"
      DB_HOST      = var.db_endpoint
      DB_URL_WRITE = "mst-aurora-db.cluster-cvq42ycqkt4f.us-west-1.rds.amazonaws.com"
      DB_URL_READ  = "mst-aurora-db.cluster-ro-cvq42ycqkt4f.us-west-1.rds.amazonaws.com"
      DB_USER      = "mst_admin"
      DB_PORT      = 5432
      DB_PASSWORD  = var.db_password
    }
  }

  role = aws_iam_role.lambda_role.arn
}

# Create IAM role for the Lambda function
resource "aws_iam_role" "lambda_role" {
  name = "lambda_backend_api_role"

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
resource "aws_iam_role_policy_attachment" "backend_api_network_attachment" {
  role       = aws_iam_role.cognito_lambda.name
  policy_arn = aws_iam_policy.lambda_network_policy.arn
}
