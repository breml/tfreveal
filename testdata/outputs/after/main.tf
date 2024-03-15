output "string" {
  value = "some text 1"
}

output "number" {
  value = 12
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
