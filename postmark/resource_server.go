package postmark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postmarkSdk "github.com/keighl/postmark"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"color": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"api_tokens": &schema.Schema{
				Type:      schema.TypeList,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				MinItems:  1,
				Elem:      &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	d.Set("name", name)

	client := m.(*postmarkSdk.Client)

	log.Printf("[DEBUG] Postmark Server create: %s", name)

	server, err := createServer(postmarkSdk.Server{
		Name: name,
	}, client)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Postmark Server creation: %s", err)

	d.SetId(strconv.FormatInt(server.ID, 10))

	return nil
}

func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*postmarkSdk.Client)

	server, err := client.GetServer(d.Id())

	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(server.ID, 10))
	d.Set("name", server.Name)
	d.Set("color", server.Color)
	d.Set("api_tokens", flattenStringList(server.ApiTokens))

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

// SDK

func createServer(server postmarkSdk.Server, client *postmarkSdk.Client) (postmarkSdk.Server, error) {
	res := postmarkSdk.Server{}

	err := doRequest(client, parameters{
		Method:    "POST",
		Path:      "servers",
		TokenType: "account",
		Payload:   server,
	}, &res)

	return res, err
}

type parameters struct {
	Method    string
	Path      string
	Payload   interface{}
	TokenType string
}

func doRequest(client *postmarkSdk.Client, opts parameters, dst interface{}) error {
	url := fmt.Sprintf("%s/%s", client.BaseURL, opts.Path)

	req, err := http.NewRequest(opts.Method, url, nil)
	if err != nil {
		return err
	}

	if opts.Payload != nil {
		payloadData, err := json.Marshal(opts.Payload)
		if err != nil {
			return err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(payloadData))
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	switch opts.TokenType {
	case "account":
		req.Header.Add("X-Postmark-Account-Token", client.AccountToken)

	default:
		req.Header.Add("X-Postmark-Server-Token", client.ServerToken)
	}

	log.Printf("[REQUEST] %s %s %s\n", req.RemoteAddr, req.Method, req.URL)

	res, err := client.HTTPClient.Do(req)

	log.Printf("[RESPONSE] %s\n", res.Status)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	log.Printf("[RESPONSE] Body: %s\n", body)

	if res.StatusCode != 200 {
		return errors.New(string(body[:]))
	}

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, dst)
	return err
}
