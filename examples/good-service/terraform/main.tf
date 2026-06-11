terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}

provider "local" {}

variable "environment" {
  type    = string
  default = "demo"
}

resource "local_file" "release_marker" {
  filename = "${path.module}/release-${var.environment}.txt"
  content  = "good-service demo deployment"
}
