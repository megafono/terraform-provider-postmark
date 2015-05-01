package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/megafono/terraform-provider-postmark/postmark"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: postmark.Provider,
	})
}
