terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "3.2.2"
    }

    local = {
      source  = "hashicorp/local"
      version = "2.5.1"
    }

    time = {
      source  = "hashicorp/time"
      version = "0.11.1"
    }
  }
}

provider "null" {}

provider "local" {}

provider "time" {}

resource "null_resource" "cluster_old" {
  triggers = {
    secret = "very very secure"
  }
}

resource "null_resource" "cluster" {
  triggers = {
    secret = "secure"
  }
}

resource "time_offset" "example" {
  offset_days = 7
}

resource "local_file" "foo" {
  content = jsonencode({
    array_unchanged = [
      "foo",
      "bar",
      "baz",
    ]
    array_changed = [
      "foo",
      "bar",
      "baz",
    ]
    array_removed = [
      "foo",
      "bar",
      "baz",
    ]
    array_item_removed = [
      "foo",
      "bar",
      "baz",
    ]
    array_item_added = [
      "foo",
      "bar",
      "baz",
    ]
    array_to_object = [
      "foo",
      "bar",
      "baz",
    ]
    number_unchanged = 10
    number_changed   = 10
    number_removed   = 10
    object_unchanged = {
      key = "value"
    }
    object_changed = {
      key = "value"
    }
    object_removed = {
      key = "value"
    }
    string_unchanged = "foo"
    string_changed   = "bar"
    string_removed   = "removed"
  })
  filename = "${path.module}/foo.bar"
}

resource "local_file" "string2json" {
  content  = "some random string"
  filename = "${path.module}/string2json"
}

resource "local_file" "json2string" {
  content = jsonencode({
    "key": "some random json"
  })
  filename = "${path.module}/json2string"
}

output "string" {
  value = "some text"
}
