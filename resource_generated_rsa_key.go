package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"golang.org/x/crypto/ssh"
	"log"
)

func resourceGeneratedRSAKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeneratedRSAKeyCreate,
		Read:   resourceGeneratedRSAKeyRead,
		Update: resourceGeneratedRSAKeyUpdate,
		Delete: resourceGeneratedRSAKeyDelete,
		Schema: map[string]*schema.Schema{
			"secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "to store the private key in",
			},
			"size": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Default:     4096,
				ForceNew:    true,
				Description: "of the RSA key, defaults to 4096",
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Sensitive:   true,
				Description: "the private key in PEM format",
			},
			"public_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "the public key in PEM format",
			},
			"public_key_ssh": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "the public key in SSH format",
			},
			"logical_version": &schema.Schema{
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Default:     "v1",
				ForceNew:    true,
				Description: "to force an update of the key",
			},
			"return_secret": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     false,
				Description: "in `value` if you need it",
			},
			"delete_on_destroy": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     true,
				Description: "to force deletion of the secret version",
			},
		},
	}
}

func generateKey(keySize int) (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return "", fmt.Errorf("failed to generated RSA key of size %d, %s", keySize, err)
	}

	err = privateKey.Validate()
	if err != nil {
		return "", fmt.Errorf("failed to validate generated RSA key, %s", err)
	}

	privateDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateDER,
	}

	return string(pem.EncodeToMemory(&privBlock)), nil
}

func publicKeyPEM(key rsa.PublicKey) (string, error) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&key)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key, %s", err)
	}

	var pemBlock = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	return string(pem.EncodeToMemory(pemBlock)), nil
}

func publicKeySSH(key rsa.PublicKey) (string, error) {
	publicRsaKey, err := ssh.NewPublicKey(&key)
	if err != nil {
		return "", fmt.Errorf("failed to create ssh public key, %s", err)
	}

	return string(ssh.MarshalAuthorizedKey(publicRsaKey)), nil
}

func privateKeyFromPEM(privateKey string) (*rsa.PrivateKey, error) {
	key, rest := pem.Decode([]byte(privateKey))
	if len(rest) > 0 {
		log.Printf("ignoring %d trailing bytes after private key pem encoding\n", len(rest))
	}

	if key.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("key is not a private key, but a %s", key.Type)
	}

	var parsedKey interface{}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key.Bytes)
	if err != nil {
		var pkcs1Error error
		parsedKey, pkcs1Error = x509.ParsePKCS1PrivateKey(key.Bytes)
		if pkcs1Error != nil {
			return nil, fmt.Errorf("failed to parse private key from PKCS1 format, %s", err)
		}
	}
	result := parsedKey.(*rsa.PrivateKey)
	if result == nil {
		return nil, fmt.Errorf("expected an rsa private key, a got %T\n", parsedKey)
	}

	return result, nil
}

func addPublicKeys(d *schema.ResourceData, privateKey string) error {
	var err error
	key, err := privateKeyFromPEM(privateKey)
	if err != nil {
		return err
	}

	public_key, err := publicKeyPEM(key.PublicKey)
	if err == nil {
		d.Set("public_key", public_key)
	}
	public_key, err = publicKeySSH(key.PublicKey)
	if err == nil {
		d.Set("public_key_ssh", public_key)
	}
	return err
}

func resourceGeneratedRSAKeyCreate(d *schema.ResourceData, meta interface{}) error {
	privateKey, err := generateKey(d.Get("size").(int))
	if err != nil {
		return err
	}
	err = addSecret(d, meta.(*Config), privateKey)
	if err != nil {
		return err
	}
	return addPublicKeys(d, privateKey)
}

func resourceGeneratedRSAKeyRead(d *schema.ResourceData, meta interface{}) error {
	privateKey, err := getSecret(d, (meta).(*Config))
	if err != nil {
		return err
	}
	return addPublicKeys(d, privateKey)
}

func resourceGeneratedRSAKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	// only triggered when return_secret or delete_on_destroy is changed
	privateKey, err := getSecret(d, (meta).(*Config))
	if err != nil {
		return err
	}
	return addPublicKeys(d, privateKey)
}

func resourceGeneratedRSAKeyDelete(d *schema.ResourceData, meta interface{}) error {
	return deleteSecret(d, meta.(*Config))
}
