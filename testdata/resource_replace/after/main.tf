terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "3.2.2"
    }
  }
}

provider "null" {}

resource "null_resource" "cluster" {
  triggers = {
    secret = "very very secure"
  }
}
