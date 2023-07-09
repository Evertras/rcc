module "lambda_http" {
  source = "./modules/lambda"

  name        = "http"
  binary_name = "lambda"
  prefix      = local.prefix

  environment_vars = {
    EVERTRAS_RCC_DYNAMODB_TABLE_NAME = aws_dynamodb_table.coverage.name
  }
}

resource "aws_iam_policy" "lambda_readwrite_dynamodb" {
  name        = "${local.prefix}-lambda-data-db-access"
  description = "Allows DynamoDB read/write access for writing measurements"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "dynamodb:PutItem",
        ]
        Effect = "Allow"
        Resource = [
          aws_dynamodb_table.coverage.arn,
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_access_data_dynamodb_attach" {
  role       = module.lambda_http.role_name
  policy_arn = aws_iam_policy.lambda_readwrite_dynamodb.arn
}
