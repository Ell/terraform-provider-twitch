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

data "twitch_channel" "ellg" {
    user_name = "ellg"
}

output "ellg_channel_output" {
  value = data.twitch_channel.ellg
}
