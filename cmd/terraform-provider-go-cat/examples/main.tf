terraform {
  required_providers {
    gocat = {
      version = "0.0.1"
      source  = "registry.devops.sorcero.com/general/go-cat"
    }
  }
}

provider "gocat" {}

resource "gocat_infra" "hello_world_service" {
  name = "hello-world"
  subsystem = "world"
  deployment_link = "https://en.wikipedia.org/wiki/%22Hello,_World!%22_program"
  cloud = "example"
  cloud_project_id = "example_project"
  type = "webpage"
  commit_sha = "12cadec"
}
