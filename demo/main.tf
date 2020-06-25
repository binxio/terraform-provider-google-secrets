resource "google_secret_manager_secret" "mysql-root-password" {
  secret_id = "mysql-root-password"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_generated_password" "mysql-root-password" {
  secret = google_secret_manager_secret.mysql-root-password.id
  logical_version = "v8"


  required {
      count = 1
      alphabet = "~!@#$%^&*()_+-="
  }

  required {
    count = 1
    alphabet = "1234567890"
  }

  required {
    count = 1
    alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
  }

  required {
    count = 2
    alphabet = "abcdefghijklmnopqrstuvwxyz"
  }

  return_secret = true
  delete_on_destroy = false
  provider = google-secrets
}

provider "google" {
  project = var.project
}

provider "google-secrets" {
  project = var.project
}

variable "project" {
  type = string
  default = "speeltuin-mvanholsteijn"
}
output "secret" {
  value = "${google_secret_manager_generated_password.mysql-root-password.id} = ${google_secret_manager_generated_password.mysql-root-password.value}"
}
