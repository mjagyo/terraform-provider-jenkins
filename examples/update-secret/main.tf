terraform {
  required_providers {
    jenkins = {
      source = "hashicorp.com/edu/jenkins"
    }
  }
}

provider "jenkins" {
  host     = "http://localhost:8080"
  username = "admin"
  token    = "11b79b85aaaef0653b94b4903986906680"
}

resource "jenkins_secret" "demo" {
  secret_type = "auth_pair"
  credential = {
    id = "usernamefromtf"
    scope = "GLOBAL"
    description = "xxaaa"
    username = "cxzcxz"
    password = "bar"
  }
}
