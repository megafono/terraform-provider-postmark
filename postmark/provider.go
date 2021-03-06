package postmark

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("POSTMARK_ACCOUNT_KEY", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"postmark_server": resourceServer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccountKey: d.Get("account_key").(string),
	}

	log.Println("[INFO] Initializing postmark client")

	client := config.Client()

	return client, nil
}
