provider "local" {}

variable "db_password" {
  type    = string
  default = "example-password"
}

resource "local_file" "release_marker" {
  filename = "${path.module}/release.txt"
  content  = "risky-service demo deployment"
}
