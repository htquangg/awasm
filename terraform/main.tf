provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Environment = terraform.workspace
      Project     = var.project
    }
  }
}

locals {
  prefix = "${terraform.workspace}-${var.prefix}"
}

data "aws_region" "current" {}
