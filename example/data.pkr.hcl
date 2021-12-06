packer {
  required_plugins {
    gitinfo = {
      version = ">= 0.0.1"
      source = "github.com/kaminskip88/gitinfo"
    }
  }
}

data "gitinfo-repo" "repo" {}

locals {
  commit = data.gitinfo-repo.repos.commit
  branch = data.gitinfo-repo.repos.branch
}