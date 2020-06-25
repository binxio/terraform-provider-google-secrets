package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func addSecret(d *schema.ResourceData, config *Config, value string) error {
	name := d.Get("secret").(string)
	var request = secretmanagerpb.AddSecretVersionRequest{
		Parent:  name,
		Payload: &secretmanagerpb.SecretPayload{Data: []byte(value)},
	}
	response, err := config.client.AddSecretVersion(config.ctx, &request)
	if err != nil {
		return err
	}
	d.SetId(response.Name)
	d.Set("value", value)
	return nil
}

func getSecret(d *schema.ResourceData, config *Config) (string, error) {
	name := d.Id()
	response, err := config.client.AccessSecretVersion(config.ctx, &secretmanagerpb.AccessSecretVersionRequest{Name: name})
	if err != nil {
		return "", err
	}

	result := string(response.Payload.Data)
	d.SetId(name)
	if d.Get("return_secret").(bool) {
		d.Set("value", result)
	} else {
		d.Set("value", "")
	}

	return result, nil
}

func deleteSecret(d *schema.ResourceData, config *Config) error {
	if d.Get("delete_on_destroy").(bool) {
		name := d.Id()
		_, err := config.client.DisableSecretVersion(config.ctx, &secretmanagerpb.DisableSecretVersionRequest{Name: name})
		if err != nil {
			return err
		}
		_, err = config.client.DestroySecretVersion(config.ctx, &secretmanagerpb.DestroySecretVersionRequest{Name: name})
		if err != nil {
			return err
		}
	}
	return nil
}
