
resource "aws_lambda_function" "cognito_reciever" {
  function_name = "cognito_reciever"
  description   = "Handler to receive sign up hooks from cognito and add users to DB"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = "${path.module}/lambda/cognito_receiver.zip"
  source_code_hash = filebase64sha256("${path.module}/lambda/cognito_receiver.zip")

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
resource "aws_iam_role" "cognito_lambda" {
  name = "cognito_lambda_role"



  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = ["sts:AssumeRole"]
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "cognito_lambda_logging_attachment" {
  role       = aws_iam_role.cognito_lambda.name
  policy_arn = aws_iam_policy.lambda_logging_policy.arn
}
