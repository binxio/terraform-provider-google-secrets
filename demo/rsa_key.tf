

resource "google_secret_manager_secret" "my-rsa-key" {
  secret_id = "my-rsa-key"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_generated_rsa_key" "my-rsa-key" {
  secret            = google_secret_manager_secret.my-rsa-key.id
  size              = 4096
  return_secret     = true
  delete_on_destroy = true

  provider = google-secrets
}

output "private-key" {
  value = "${google_secret_manager_generated_rsa_key.my-rsa-key.value}"
}

output "public-key" {
  value = "${google_secret_manager_generated_rsa_key.my-rsa-key.public_key}"
}
output "public-key-ssh" {
  value = "${google_secret_manager_generated_rsa_key.my-rsa-key.public_key_ssh}"
}
