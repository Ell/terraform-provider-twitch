terraform {
  required_providers {
    twitch = {
      source = "ellg/twitch"
    }
  }
}

variable client_id {
  type = string
  sensitive = true
}

variable app_token {
  type = string
  sensitive = true
}

provider "twitch" {
  client_id     = var.client_id
  app_token     = var.app_token
}

resource "twitch_channel" "ellg" {
  title = "twitch terraform provider development | this was set from terraform"
  tags = ["coding", "testing", "terraform"]
}
