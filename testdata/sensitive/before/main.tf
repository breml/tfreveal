locals {
  secret_value = sensitive("some secret value")
}

resource "null_resource" "sensitive" {
  triggers = {
    secret = local.secret_value
  }
}
