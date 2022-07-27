terraform {
  required_providers {
    gocat = {
      version = "0.0.12"
      source  = "sorcero/go-cat"
    }
  }
}

provider "gocat" {}

resource "gocat_infra" "hello_world_service" {
  for_each = toset(["person", "human", "world"])
  name = "hello-${each.value}"
  subsystem = "world"
  deployment_links = ["https://en.wikipedia.org/wiki/%22Hello,_World!%22_program", "srev.in"]
  cloud = "example"
  cloud_project_id = "example_project"
  type = "webpage"
  commit_sha = "12cadec"
}

data "gocat_infra" "hello_world_service" {
  depends_on = [gocat_infra.hello_world_service]

  id = "example/example_project/world/hello-world"
}

