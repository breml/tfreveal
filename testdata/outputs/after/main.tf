output "string" {
  value = "some text 1"
}

output "number_equal" {
  value = 5
}

output "number_change" {
  value = 12
}

output "number_add" {
  value = 30
}

output "object" {
  value = {
    key1 = "1"
    key2 = "2 added"
    key3 = "3 changed"
  }
}

output "array" {
  value = [
    "foo",
    "more",
    "bar",
  ]
}

output "jsonString" {
  value = jsonencode({
    key1 = "1"
    key2 = "2 added"
    key3 = "3 changed"
  })
}
