variable "region" {
  type    = string
  default = "ap-southeast-1"
}

variable "prefix" {
  type    = string
  default = "dev-awasm"
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
  default = "dev-awasm-cd-001"
}
