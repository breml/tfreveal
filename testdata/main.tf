resource "null_resource" "cluster" {
  triggers = {
    secret = "very very secure 1"
    secret = "foobar"
  }
}
