package main

import (
	"fmt"
	"os"
	repoData "packer-plugin-gitinfo/datasource/repo"
	version "packer-plugin-gitinfo/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource("repo", new(repoData.Datasource))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
