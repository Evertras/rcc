variable "name" {
  description = "The name of the lambda."
  type        = string
}

variable "binary_name" {
  description = "The path to the built Go binary of the lambda in the repository's ./bin directory."
  type        = string
}

variable "prefix" {
  description = "The prefix to use for naming"
  type        = string
}

variable "environment_vars" {
  description = "Environment variables to apply to the lambda"
  type        = map(string)
  default     = {}
}

variable "policies" {
  description = "Additional policies to attach"
  type        = set(string)
  default     = []
}

locals {
  prefix = "${var.prefix}-${var.name}"
}
