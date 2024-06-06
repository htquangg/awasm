variable "region" {
  type    = string
  default = "ap-southeast-1"
}

variable "prefix" {
  type    = string
  default = "awasm"
}

variable "project" {
  type    = string
  default = "awasm"
}

variable "tf_state_bucket" {
  type    = string
  default = "dev-awasm-tf-state"
}

variable "tf_state_lock_table" {
  type    = string
  default = "dev-awasm-tf-lock"
}

variable "iam_user" {
  type    = string
  default = "awasm-cd-001"
}

variable "ecr_name" {
  type    = string
  default = "stg-awasm"
}

variable "ecr_app_image" {
  type        = string
  description = "Path to the ECR repo with the API image"
}

variable "db_username" {
  description = "Username for the app api database"
  type        = string
  default     = "awasm"
}

variable "db_password" {
  description = "Password for the Terraform database"
  type        = string
  sensitive   = true
}

variable "db_schema" {
  type    = string
  default = "stg_local_awasm_001"
}
