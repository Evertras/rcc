module "lambda_http" {
  source = "./modules/lambda"

  name        = "http"
  binary_name = "lambda"
  prefix      = local.prefix
}
