data "archive_file" "lambda_code" {
  type = "zip"

  source_file = "${path.module}/../../../bin/${var.binary_name}"
  output_path = "${path.module}/.archives/lambda-code-${var.name}.zip"
}

resource "aws_lambda_function" "lambda" {
  function_name = local.prefix

  filename = data.archive_file.lambda_code.output_path

  source_code_hash = filebase64sha256(data.archive_file.lambda_code.output_path)

  runtime = "go1.x"
  handler = var.binary_name

  environment {
    variables = merge(
      {
        "DEPLOY_ENVIRONMENT" = terraform.workspace == "default" ? "prod" : terraform.workspace,
      },
      var.environment_vars
    )
  }

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_iam_role" "lambda_exec" {
  name = "${local.prefix}-exec"

  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "lambda.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF
}

resource "aws_iam_policy" "log_policy" {
  name = "${local.prefix}-log-policy"

  description = "Policy to write logs"

  policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": [
       "logs:CreateLogStream",
       "logs:PutLogEvents"
     ],
     "Resource": "arn:aws:logs:*:*:*",
     "Effect": "Allow"
   }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_log_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.log_policy.arn
}
