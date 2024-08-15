terraform {
  required_providers {
    twitch = {
      source = "ellg/twitch"
    }
  }
}

provider "twitch" {
  client_id     = "test"
  app_token     = "test"
}