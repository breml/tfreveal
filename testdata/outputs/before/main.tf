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

output "nested_object" {
  value = {
    key1 = {
      subkey1 = "1.1"
    }
  }
}
