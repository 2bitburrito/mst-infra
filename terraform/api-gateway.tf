# ! NOT YET DEPLOYED
module "api_gateway" {
  source = "terraform-aws-modules/apigateway-v2/aws"

  name          = "mst-api-gateway"
  description   = "Main API gateway for mst backend"
  protocol_type = "HTTP"

  cors_configuration = {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }

  # Custom domain - Awaiting purchase of mst.com
  domain_name = "api.metasoundtools.com"

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

  # Authorizer(s)
  authorizers = {
    "api" = {
      authorizer_type = "REQUEST"
      name            = "api-auth"
    }
  }

  # Routes & Integration(s)
  routes = {
    "POST/check-license" = {
      integration = {
        uri                    = module.check_license_lambda.lambda_function_invoke_arn
        payload_format_version = "2.0"
        timeout_milliseconds   = 12000
      }
    }
  }

  tags = {
    Environment = "prod"
    Terraform   = "true"
  }
}
