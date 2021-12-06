data "gitinfo-repo" "repo" {}

locals {
  commit = data.gitinfo-repo.repo.commit
  branch = data.gitinfo-repo.repo.branch
  tags = data.gitinfo-repo.repo.tags
  latest_tag = data.gitinfo-repo.repo.latest_tag
}

source "null" "test" {
  communicator = "none"
}

build {
  sources = [
    "source.null.test"
  ]

  provisioner "shell-local" {
    inline = [
      "echo commit: ${local.commit}",
      "echo branch: ${local.branch}",
      "echo tags: ${join(",",local.tags)}",
      "echo latest_tag: ${local.latest_tag}",
    ]
  }
}
