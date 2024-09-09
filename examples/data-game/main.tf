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

data "twitch_game" "programming" {
  name = "Software and Game Development"
}

output "programming_game_output" {
  value = data.twitch_game
}
