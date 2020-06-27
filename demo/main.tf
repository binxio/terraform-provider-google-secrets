
provider "google" {
  project = var.project
}

provider "google-secrets" {
  project = var.project
}

variable "project" {
  type    = string
}
