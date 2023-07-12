locals {
  prefix = terraform.workspace == "default" ? "evertras-rcc-ecs" : "evertras-rcc-ecs-${terraform.workspace}"
  port   = 8341
}

variable "app_count" {
  description = "How many instances to run"
  type        = number
  default     = 2
}
