# terraform-provider-postmark

## Installation 

First, following the [local setup](https://www.terraform.io/docs/extend/writing-custom-providers.html#local-setup) from Terraform.

After, run:

```
$ mkdir -p ~/.terraform.d/plugins
$ git clone git@github.com:megafono/terraform-provider-postmark.git
$ cd terraform-provider-postmark
$ make install
```

## Usage

Setup the provider

```
provider "postmark" {
  account_key = "xxxxxxxxxxxxxxxxxxx"
}
```

### Creating a server

```
resource "postmark_server" "default" {
  name = "myserver"
}
```
