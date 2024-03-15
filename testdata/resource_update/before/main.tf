terraform {
  required_providers {
    time = {
      source  = "hashicorp/time"
      version = "0.11.1"
    }
  }
}

provider "time" {}

resource "time_offset" "example" {
  offset_days = 7
}
