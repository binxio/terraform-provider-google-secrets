## google\_secret\_manage\_generated\_password
A generated secret version resource.

This will generate a secret and store the value directly in the Google Secret manager secret, to avoid having the secret to be specified in the terraform file.

### Example basic usage
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

### Argument reference
The following arguments are supported:

secret - (Required) Secret Manager secret resource.
length - (Optional) the length of the secret to generate, default = 32.
alphabet - (Optional) the characters to generate the secret from, default = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789".
required - (Optional) array of required characters in the secret, specifying the minimum number of characters and the alphabet for each.
logical\_version - (Optional) an opaque string to force the secret to be regenerated.
return\_secret - (Optional) indicating whether the generated secret should be return as plain text `value`, default false
delete\_on\_destroy - (Optional) indicating whether the version should be deleted when the resource is destroyed,  default true

### Attribute reference
In addition to the arguments listed above, the following computed attributes are exported:

id - The name of the SecretVersion. Format: projects/{{project}}/secrets/{{secret\_id}}/versions/{{version}}
value - The generate value, if return\_secret is true, otherwise "".
