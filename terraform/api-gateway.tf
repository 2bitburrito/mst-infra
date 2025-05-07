module "api_gateway" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "~>5.0"

  name          = "mst-api-gateway"
  description   = "Main API gateway for mst backend"
  protocol_type = "HTTP"

  cors_configuration = {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }
  create_domain_name = false

  # Access logs
  stage_access_log_settings = {
    create_log_group            = true
    log_group_retention_in_days = 7
    format = jsonencode({
      context = {
        domainName              = "$context.domainName"
        integrationErrorMessage = "$context.integrationErrorMessage"
        protocol                = "$context.protocol"
        requestId               = "$context.requestId"
        requestTime             = "$context.requestTime"
        responseLength          = "$context.responseLength"
        routeKey                = "$context.routeKey"
        stage                   = "$context.stage"
        status                  = "$context.status"
        error = {
          message      = "$context.error.message"
          responseType = "$context.error.responseType"
        }
        identity = {
          sourceIP = "$context.identity.sourceIp"
        }
        integration = {
          error             = "$context.integration.error"
          integrationStatus = "$context.integration.integrationStatus"
        }
      }
    })
  }

  #? Change this to cognito when out of dev
  #   # Authorizer(s)
  #   authorizers = {
  #     "api" = {
  #       authorizer_type  = "REQUEST"
  #       name             = "api-auth"
  #       authorizer_uri   = aws_iam_role.lambda_role.arn
  #       identity_sources = ["$request.header.Authorization"]
  #     }
  #   }

  # Routes & Integration(s)

  routes = {
    "POST /check_license" = {
      integration = {
        uri                    = aws_lambda_function.check_license.invoke_arn
        payload_format_version = "2.0"
        timeout_milliseconds   = 1200
      }
    }
  }

  tags = {
    Environment = "prod"
    Terraform   = "true"
  }
}

resource "aws_lambda_permission" "api_gw_auth" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.check_license.function_name
  principal     = "apigateway.amazonaws.com"
}
output "api_gateway_endpoint" {
  value = module.api_gateway.api_endpoint
}
