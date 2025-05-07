resource "aws_amplify_app" "mst-website" {
  name         = "mst-website"
  repository   = "https://github.com/2bitburrito/mst-website"
  access_token = var.mst_website_github_token

  # The default build_spec added by the Amplify Console for React.
  build_spec = <<-EOT
    version: 1
    frontend:
      phases:
        preBuild:
          commands:
            - npm install
        build:
          commands:
            - npm run build
      artifacts:
        baseDirectory: out
        files:
          - "**/*"
      customHeaders:
        - pattern: "**/_next/static/**"
          headers:
            - key: "Cache-Control"
              value: "public, max-age=31536000, immutable"
        - pattern: "**/*.html"
          headers:
            - key: "Cache-Control"
              value: "public, max-age=0, must-revalidate"
        - pattern: "**/*.css"
          headers:
            - key: "Content-Type"
              value: "text/css"
        - pattern: "**/*.js"
          headers:
            - key: "Content-Type"
              value: "application/javascript"
      cache:
        paths:
          - node_modules/**/*
          - .next/cache/**/*
  EOT

  # The default rewrites and redirects added by the Amplify Console.
  custom_rule {
    source = "</^[^.]+$|\\.(?!(css|gif|ico|jpg|js|png|txt|svg|woff|woff2|ttf|map|json)$)([^.]+$)/>"
    status = "200"
    target = "/index.html"
  }

  environment_variables = {
    ENV = "test"
  }
}
