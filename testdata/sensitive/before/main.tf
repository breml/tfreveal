locals {
  secret_value   = sensitive("some secret value")
  sensitive_json = sensitive(jsonencode({ key = "value" }))
  json_with_nested_sensitive = jsonencode({
    key = sensitive("value")
  })
}

resource "null_resource" "sensitive" {
  triggers = {
    secret = local.secret_value
  }
}

output "null_resource" {
  value     = nonsensitive(null_resource.sensitive)
  sensitive = true
}

output "sensitive" {
  value = nonsensitive(null_resource.sensitive.triggers.secret)
}

output "sensitive_json" {
  value = nonsensitive(local.sensitive_json)
}

output "json_with_nested_sensitive" {
  value = nonsensitive(local.json_with_nested_sensitive)
}
