terraform {
  cloud {
    organization = "Meta-Sound-Tools"

    workspaces {
      name = "mst-infra"
    }
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.83.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.4.3"
    }
  }
  required_version = "~> 1.11.4"
}

provider "aws" {
  region = var.region
}
