resource "aws_dynamodb_table" "rcc" {
  name      = "RCC"
  hash_key  = "Key"

  billing_mode = "PAY_PER_REQUEST"

  # Only define the key attributes here
  attribute {
    name = "Key"
    type = "S"
  }
}
