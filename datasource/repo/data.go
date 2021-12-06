//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package gitinfo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	Path string `mapstructure:"path"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	Commit    string   `mapstructure:"commit"`
	Branch    string   `mapstructure:"branch"`
	Tags      []string `mapstructure:"tags"`
	LatestTag string   `mapstructure:"latest_tag"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}
	if d.config.Path == "" {
		d.config.Path = "."
	}
	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	output := DatasourceOutput{}
	emptyOutput := hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec())
	r, err := git.PlainOpenWithOptions(d.config.Path,
		&git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return emptyOutput, err
	}
	ref, err := r.Head()
	if err != nil {
		return emptyOutput, err
	}
	output.Commit = ref.Hash().String()
	output.Branch = ref.Name().Short()

	tagrefs, err := r.Tags()
	if err != nil {
		return emptyOutput, err
	}
	var tags []string
	var latestTag string
	var latestTagCommit *object.Commit
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		tags = append(tags, t.Name().Short())
		revision := plumbing.Revision(t.Name().String())
		tagCommitHash, err := r.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := r.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTag = t.Name().Short()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTag = t.Name().Short()
		}
		return nil
	})
	if err != nil {
		return emptyOutput, err
	}
	output.Tags = tags
	output.LatestTag = latestTag
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
