provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Environment = terraform.workspace
      Project     = "awasm"
    }
  }
}

locals {
  prefix = "${var.prefix}-${terraform.workspace}"
}
