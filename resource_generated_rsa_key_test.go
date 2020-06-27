package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"testing"
)

func decodePublicKeyFromPEM(str string) (*rsa.PublicKey, error) {
	key, _ := pem.Decode([]byte(str))

	parsedKey, err := x509.ParsePKIXPublicKey(key.Bytes)
	if err != nil {
		return nil, err
	}
	return parsedKey.(*rsa.PublicKey), nil
}

func TestPublicKeyGeneratorPEM(t *testing.T) {
	privateKeyPEM, err := generateKey(4096)
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := privateKeyFromPEM(privateKeyPEM)
	if err != nil {
		t.Fatal(err)
	}

	publicKeyPEM, err := publicKeyPEM(privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := decodePublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		t.Fatal(err)
	}
	if publicKey.E != privateKey.PublicKey.E {
		t.Fatal("public key E is incorrect")
	}
	if publicKey.N.Cmp(privateKey.PublicKey.N) != 0 {
		t.Fatal("public key N is incorrect")
	}

}

func TestPublicKeyGeneratorSSH(t *testing.T) {

	privateKeyPEM, err := generateKey(4096)
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := privateKeyFromPEM(privateKeyPEM)
	if err != nil {
		t.Fatal(err)
	}

	publicKeySSH, err := publicKeySSH(privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKeySSH))
	if err != nil {
		t.Fatal(err)
	}
	sshPublicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(sshPublicKey.Marshal(), publicKey.Marshal()) {
		t.Fatal("public key is incorrect")
	}
}
