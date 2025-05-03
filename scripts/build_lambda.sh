#!/bin/bash
set -e

echo "Building Lambda function for ARM64..."

# Navigate to Lambda function code
cd "$(dirname "$0")/../lambda/check_license"

# Clean existing artifacts
rm -f bootstrap function.zip

# Build for Linux ARM64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap .

# Create ZIP package
zip -j function.zip bootstrap

# Copy to Terraform directory
mkdir -p ../../terraform/lambda
cp function.zip ../../terraform/lambda/

echo "Lambda package built successfully: ../../terraform/lambda/function.zip"

