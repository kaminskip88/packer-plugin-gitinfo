package gitinfo

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

//go:embed test-fixtures/template.pkr.hcl
var testDatasourceHCL2Basic string

func TestAccGitinfoDatasource(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "gitinfo_datasource_basic_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testDatasourceHCL2Basic,
		Type:     "gitinfo-repo-datasource",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := ioutil.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)

			commitLog := "null.test: commit: \\w{40}"
			branchLog := "null.test: branch: \\w+"

			if matched, _ := regexp.MatchString(commitLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected foo value %q", logsString)
			}
			if matched, _ := regexp.MatchString(branchLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected bar value %q", logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
