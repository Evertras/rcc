locals {
  prefix = terraform.workspace == "default" ? "evertras-rcc" : "evertras-rcc-${terraform.workspace}"
}
