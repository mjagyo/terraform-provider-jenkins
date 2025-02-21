<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Jenkins Provider for Terraform

This is the [Jenkins](https://www.jenkins.io/) provider for [Terraform](https://www.terraform.io/).

This provider allows you to manage [Secrets](https://www.jenkins.io/doc/developer/security/secrets/) and [Jobs](https://www.jenkins.io/doc/book/using/working-with-projects/) in your Jenkins cluster using Terraform.


## Contents

* [Requirements](#requirements)
* [Getting Started](#getting-started)
* [Contributing to the provider](#contributing)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) v0.12.x
-	[Go](https://golang.org/doc/install) v1.18.x (to build the provider plugin)

## Getting Started

This is a small project of how to manage Jenkins Secrets (Credentials) and Jobs.
information.

```hcl
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
  token    = "xxxx"
}

resource "jenkins_secret" "demo" {
  secret_type = "auth_pair"
  credential = {
    id          = "identityterraform"
    scope       = "GLOBAL"
    description = "A foobar credentials"
    username    = "foo"
    password    = "bar"
  }
}

resource "jenkins_job" "demo" {
  name = "new_job"
  file = "./job-config.xml"
}
```

job-config.xml
```xml 
<flow-definition plugin="workflow-job@1505.vea_4b_20a_4a_495">
	<keepDependencies>false</keepDependencies>
	<properties/>
	<triggers/>
	<disabled>false</disabled>
</flow-definition>
```


## Contributing
Coming up...
