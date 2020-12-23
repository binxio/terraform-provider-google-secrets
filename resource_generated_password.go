package main

import (
	c "crypto/rand"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

func resourceGeneratedPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceGeneratedPasswordCreate,
		Read:   resourceGeneratedPasswordRead,
		Update: resourceGeneratedPasswordUpdate,
		Delete: resourceGeneratedPasswordDelete,
		Schema: map[string]*schema.Schema{
			"secret": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"length": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				Default:  32,
				ForceNew: true,
			},
			"alphabet": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
				ForceNew: true,
			},
			"required": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  32,
							ForceNew: true,
						},
						"alphabet": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"logical_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "v1",
				ForceNew: true,
			},
			"return_secret": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"delete_on_destroy": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"value": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  false,
				Required:  false,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func generatePassword(d *schema.ResourceData) string {
	var password strings.Builder

	seed, _ := c.Int(c.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	length := d.Get("length").(int)
	alphabet := d.Get("alphabet").(string)

	required := d.Get("required").([]interface{})
	for _, r := range required {
		count := r.(map[string]interface{})["count"].(int)
		alphabet := r.(map[string]interface{})["alphabet"].(string)
		for i := 0; i < count; i++ {
			random := rand.Intn(len(alphabet))
			password.WriteString(string(alphabet[random]))
		}
		length = length - count
	}

	for i := 0; i < length; i++ {
		random := rand.Intn(len(alphabet))
		password.WriteString(string(alphabet[random]))
	}

	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func resourceGeneratedPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	return addSecret(d, meta.(*Config), generatePassword(d))
}

func resourceGeneratedPasswordRead(d *schema.ResourceData, meta interface{}) error {
	_, err := getSecret(d, (meta).(*Config))
	return err
}

func resourceGeneratedPasswordUpdate(d *schema.ResourceData, meta interface{}) error {
	// only triggered when return_secret or delete_on_destroy is changed
	_, err := getSecret(d, (meta).(*Config))
	return err
}

func resourceGeneratedPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	return deleteSecret(d, meta.(*Config))
}
