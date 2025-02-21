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

resource "jenkins_job" "demo" {
  name = "fromterraform"
  file = "./changed-job.xml"
}
