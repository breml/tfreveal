locals {
  secret_value = sensitive("new secret value")
}

resource "null_resource" "sensitive" {
  triggers = {
    secret = local.secret_value
  }
}
