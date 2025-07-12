# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    pelican = {
      source = "hashicorp.com/edu/pelican"
    }
  }
}

provider "pelican" {}

data "pelican_coffees" "example" {}
