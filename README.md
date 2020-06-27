# google secrets provider
The google secrets provider aims to provide secrets by generating them and storing them directly in the
  Google Secret Manager, to avoid having the specify the secrets in the terraform file or have them appear
  in plain text in the state file.
  
Currently there are two resources:
- [google\_secret\_manage\_generated\_password](#google_secret_manage_generated_password) - for passwords
- [google\_secret\_manage\_generated\_rsa_key](#google_secret_manage_generated_rsa_key) - for RSA keys

## google\_secret\_manage\_generated\_password
A generated secret version resource.

This will generate a secret and store the value directly in the Google Secret manager secret, 
  to avoid having the secret to be specified in the terraform file.

## Example basic usage
```
resource "google_secret_manager_secret" "mysql_user_password" {
  secret_id = "mysql-user-password"
}

resource "google_generated_password" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id

  length = 20
  alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  required = [
    {
       count = 1
       alphabet = "012356789"
    }
    {
       count = 2
       alphabet = "@!#$%^&*()_+-=:;<>,./?"
    }
  ]
  logical_version = "v1"
 return_secret = true
  delete_on_destroy = true
}
```

## Argument reference
The following arguments are supported:

```
secret - (Required) Secret Manager secret resource.
length - (Optional) the length of the secret to generate, default = 32.
alphabet - (Optional) the characters to generate the secret from, default = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789".
required - (Optional) array of required characters in the secret, specifying the minimum number of characters and the alphabet for each.
logical_version - (Optional) an opaque string to force the secret to be regenerated.
return_secret - (Optional) indicating whether the generated secret should be return as plain text `value`, default false
delete_on_destroy - (Optional) indicating whether the version should be deleted when the resource is destroyed,  default true
```

### Attribute reference
In addition to the arguments listed above, the following computed attributes are exported:
```

id - The name of the SecretVersion. Format: projects/{{project}}/secrets/{{secret\_id}}/versions/{{version}}
value - The generate value, if return\_secret is true, otherwise "".
```


## google\_secret\_manager\_generated\_rsa_key

This will generate a RSA Key and storeit directly in the Google Secret manager secret, to avoid having the secret to 
  be specified in the terraform file.

## Example basic usage
```


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
```

## Argument reference
The following arguments are supported:

```
size - (Optional) number of bits in the key, default 4096.
logical_version - (Optional) an opaque string to force a new key  to be regenerated.
return_secret - (Optional) indicating whether the generated key should be return as plain text `value`, default false
delete_on_destroy - (Optional) indicating whether the version should be deleted when the resource is destroyed,  default true
```

### Attribute reference
In addition to the arguments listed above, the following computed attributes are exported:
```

id - The name of the SecretVersion. Format: projects/{{project}}/secrets/{{secret\_id}}/versions/{{version}}
value - The generate private key in PEM format, if return\_secret is true, otherwise "".
public_key - the public key of the generated private key in PEM format.
public_key_ssh - the public key of the generated private key in SSH format.
```
