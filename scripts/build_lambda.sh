#!/bin/bash
set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <function_name>"
  exit 1
fi

echo "Building Lambda function for ARM64..."

# Navigate to Lambda function code
cd "$(dirname "$0")/../lambda/${1:-}" || exit 1

# Clean existing artifacts
rm -f bootstrap "$1".zip

# Build for Linux ARM64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap .

# Create ZIP package
zip -j "$1".zip bootstrap

# Copy to Terraform directory
mkdir -p ../../terraform/lambda
cp "$1".zip ../../terraform/lambda/

echo "Lambda package built successfully: ../../terraform/lambda/$1.zip"

if [ "$2" == "-deploy" ]; then
  echo "Lambda uploading to AWS: ../../terraform/lambda/$1.zip"
  echo aws CLI args: \"--function-name "$1" --zip-file fileb://../../terraform/lambda/"$1".zip\"
  aws lambda update-function-code --function-name "$1" --zip-file fileb://../../terraform/lambda/"$1".zip
fi
