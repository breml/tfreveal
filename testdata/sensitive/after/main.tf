locals {
  secret_value = sensitive("new secret value")
  other_value = jsonencode(
    {
      key = "value"
    }
  )
}

resource "null_resource" "sensitive" {
  triggers = {
    secret = local.secret_value
  }
}

output "other_value" {
  value = local.other_value
}
