locals {
  secret_value   = sensitive("new secret value")
  sensitive_json = sensitive(jsonencode({ key = "value2" }))
  json_with_nested_sensitive = jsonencode({
    key = sensitive("value2")
  })
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

output "null_resource" {
  value     = nonsensitive(null_resource.sensitive)
  sensitive = true
}

output "other_value" {
  value = local.other_value
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
