package main

import (
	"cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	//"google.golang.org/api/option"
	"log"
	"regexp"
)

// Config is the configuration structure used to instantiate the Google
// provider.
type Config struct {
	Credentials         string
	AccessToken         string
	Project             string
	Region              string
	Zone                string
	Scopes              []string
	UserProjectOverride bool

	terraformVersion  string
	tokenSource       oauth2.TokenSource
	SecretManagerPath string
	client            *secretmanager.Client
	ctx               context.Context
}

var defaultClientScopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
}

func (c *Config) LoadAndValidate() error {
	c.ctx = context.Background()

	if len(c.Scopes) == 0 {
		c.Scopes = defaultClientScopes
	}

	tokenSource, err := c.getTokenSource(c.Scopes)
	if err != nil {
		return err
	}
	c.tokenSource = tokenSource

	log.Printf("[INFO] Instantiating SecretManager Client")
	c.client, err = secretmanager.NewClient(c.ctx, option.WithTokenSource(c.tokenSource))
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) getTokenSource(clientScopes []string) (oauth2.TokenSource, error) {
	if c.AccessToken != "" {
		contents, _, err := pathorcontents.Read(c.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("Error loading access token: %s", err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		token := &oauth2.Token{AccessToken: contents}
		return oauth2.StaticTokenSource(token), nil
	}

	if c.Credentials != "" {
		contents, _, err := pathorcontents.Read(c.Credentials)
		if err != nil {
			return nil, fmt.Errorf("Error loading credentials: %s", err)
		}

		creds, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(contents), clientScopes...)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse credentials from '%s': %s", contents, err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return creds.TokenSource, nil
	}

	log.Printf("[INFO] Authenticating using DefaultClient...")
	log.Printf("[INFO]   -- Scopes: %s", clientScopes)
	return googleoauth.DefaultTokenSource(context.Background(), clientScopes...)
}

var SecretManagerDefaultBasePath = "https://secretmanager.googleapis.com/v1/"

// Remove the `/{{version}}/` from a base path, replacing it with `/`
func removeBasePathVersion(url string) string {
	return regexp.MustCompile(`/[^/]+/$`).ReplaceAllString(url, "/")
}

func ConfigureBasePaths(c *Config) {
	c.SecretManagerPath = SecretManagerDefaultBasePath
}
