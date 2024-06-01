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

variable "db_username" {
  description = "Username for the recipe app api database"
  type        = string
  default     = "awasm"
}

variable "db_password" {
  description = "Password for the Terraform database"
  type        = string
  sensitive   = true
}
