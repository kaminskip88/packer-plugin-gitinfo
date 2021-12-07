packer {
  required_plugins {
    gitinfo = {
      version = ">= 0.0.1"
      source = "github.com/kaminskip88/gitinfo"
    }
  }
}

data "gitinfo-repo" "repo" {
  path = "."  # default
}

locals {
  commit = data.gitinfo-repo.repo.commit
  branch = data.gitinfo-repo.repo.branch
  tags = data.gitinfo-repo.repo.tags
  latest_tag = data.gitinfo-repo.repo.latest_tag
}
