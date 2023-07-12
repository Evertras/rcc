locals {
  prefix = terraform.workspace == "default" ? "evertras-rcc-ecs" : "evertras-rcc-ecs-${terraform.workspace}"
  port   = 8341
  zones  = 2
}

variable "app_count" {
  description = "How many instances to run"
  type        = number
  default     = 1
}

variable "rcc_version" {
  description = "The version tag to use"
  type        = string
  default     = "v0.1.0"
}
