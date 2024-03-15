resource "null_resource" "cluster" {
  triggers = {
    secret = "very very secure"
    // secret = "foobar"
  }
}

resource "local_file" "foo" {
  content = jsonencode({
    foo = "bar"
  })
  file_permission = "0666"
  filename        = "${path.module}/foo.bar"
}
