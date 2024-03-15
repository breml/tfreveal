output "string" {
  value = "some text"
}

output "number" {
  value = 10
}

output "object" {
  value = {
    key1 = "1"
    key3 = "3"
    key4 = "4 remove"
  }
}

output "array" {
  value = [
    "foo",
    "bar",
  ]
}

output "jsonString" {
  value = jsonencode({
    key1 = "1"
    key3 = "3"
    key4 = "4 remove"
  })
}
