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

resource "null_resource" "cluster_new" {
  triggers = {
    secret = "very very secure"
  }
}

resource "null_resource" "cluster" {
  triggers = {
    secret = "still secure"
  }
}

resource "time_offset" "example" {
  offset_days = 8
}

resource "local_file" "foo" {
  content = jsonencode({
    array_unchanged = [
      "foo",
      "bar",
      "baz",
    ]
    array_changed = [
      "foo2",
      "bar",
      "baz2",
    ]
    array_new = [
      "foo",
      "bar",
      "baz",
    ]
    array_item_removed = [
      "foo",
      "baz",
    ]
    array_item_added = [
      "foo",
      "bar",
      "baz",
      "biz",
    ]
    array_to_object = {
      0 = "foo"
      1 = "bar"
      2 = "baz"
    }
    number_unchanged = 10
    number_changed   = 14
    number_new       = 14
    object_unchanged = {
      key = "value"
    }
    object_changed = {
      key = "new value"
    }
    object_new = {
      key = "value"
    }
    string_unchanged = "foo"
    string_changed   = "bar changed"
    string_new       = "new"
  })
  filename        = "${path.module}/foo.bar"
  file_permission = "0660"
}

resource "local_file" "string2json" {
  content = jsonencode({
    "key": "some random json"
  })
  filename = "${path.module}/string2json"
}

resource "local_file" "json2string" {
  content  = "some random string"
  filename = "${path.module}/json2string"
}

output "string" {
  value = "some text 1"
}
