package cmd

/*import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/loft-sh/devspace/cmd/flags"
	cloudpkg "github.com/loft-sh/devspace/pkg/devspace/cloud"
	cloudconfig "github.com/loft-sh/devspace/pkg/devspace/cloud/config"
	cloudlatest "github.com/loft-sh/devspace/pkg/devspace/cloud/config/versions/latest"
	"github.com/loft-sh/devspace/pkg/devspace/config/loader"
	"github.com/loft-sh/devspace/pkg/devspace/config/constants"
	"github.com/loft-sh/devspace/pkg/devspace/config/localcache"
	"github.com/loft-sh/devspace/pkg/devspace/config/versions/latest"
	"github.com/loft-sh/devspace/pkg/util/fsutil"
	"github.com/loft-sh/devspace/pkg/util/kubeconfig"
	"github.com/loft-sh/devspace/pkg/util/log"
	"github.com/loft-sh/devspace/pkg/util/survey"
	"k8s.io/client-go/tools/clientcmd"

	"gopkg.in/yaml.v3"
	"gotest.tools/assert"
)

type runTestCase struct {
	name string

	fakeConfig       *latest.Config
	fakeKubeConfig   clientcmd.ClientConfig
	files            map[string]interface{}
	graphQLResponses []interface{}
	providerList     []*cloudlatest.Provider
	answers          []string
	args             []string

	globalFlags flags.GlobalFlags

	expectedErr string
}

func TestRun(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}

	wdBackup, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Error changing working directory: %v", err)
	}
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		//Delete temp folder
		err = os.Chdir(wdBackup)
		if err != nil {
			t.Fatalf("Error changing dir back: %v", err)
		}
		err = os.RemoveAll(dir)
		if err != nil {
			t.Fatalf("Error removing dir: %v", err)
		}
	}()

	testCases := []runTestCase{
		runTestCase{
			name: "empty command",
			files: map[string]interface{}{
				constants.DefaultConfigPath: &latest.Config{
					Version: latest.Version,
				},
			},
			args:        []string{""},
			expectedErr: "execute command: Couldn't find command '' in devspace.yaml",
		},
		runTestCase{
			name: "shell error",
			files: map[string]interface{}{
				constants.DefaultConfigPath: &latest.Config{
					Version: latest.Version,
					Commands: []*latest.CommandConfig{
						&latest.CommandConfig{
							Name:    "exit",
							Command: "exit",
						},
					},
				},
			},
			args:        []string{"exit", "1"},
			expectedErr: "exit code 1",
		},
		runTestCase{
			name: "command error",
			files: map[string]interface{}{
				constants.DefaultConfigPath: &latest.Config{
					Version: latest.Version,
					Commands: []*latest.CommandConfig{
						&latest.CommandConfig{
							Name:    "fail",
							Command: "ThisCommandDoesntExist",
						},
					},
				},
			},
			args:        []string{"fail"},
			expectedErr: "exit code 127",
		},
		runTestCase{
			name: "successful run",
			files: map[string]interface{}{
				constants.DefaultConfigPath: &latest.Config{
					Version: latest.Version,
					Commands: []*latest.CommandConfig{
						&latest.CommandConfig{
							Name:    "exit",
							Command: "exit",
						},
					},
				},
			},
			args: []string{"exit", "0"},
		},
	}

	log.SetInstance(&log.DiscardLogger{PanicOnExit: true})

	for _, testCase := range testCases {
		testRun(t, testCase)
	}
}

func testRun(t *testing.T, testCase runTestCase) {
	defer func() {
		for path := range testCase.files {
			removeTask := strings.Split(path, "/")[0]
			err := os.RemoveAll(removeTask)
			assert.NilError(t, err, "Error cleaning up folder in testCase %s", testCase.name)
		}
		err := os.RemoveAll(log.Logdir)
		assert.NilError(t, err, "Error cleaning up folder in testCase %s", testCase.name)
	}()

	cloudpkg.DefaultGraphqlClient = &customGraphqlClient{
		responses: testCase.graphQLResponses,
	}

	for _, answer := range testCase.answers {
		survey.SetNextAnswer(answer)
	}

	providerConfig, err := cloudconfig.Load()
	assert.NilError(t, err, "Error getting provider config in testCase %s", testCase.name)
	providerConfig.Providers = testCase.providerList

	loader.SetFakeConfig(testCase.fakeConfig)
	loader.ResetConfig()
	generated.ResetConfig()
	kubeconfig.SetFakeConfig(testCase.fakeKubeConfig)

	for path, content := range testCase.files {
		asYAML, err := yaml.Marshal(content)
		assert.NilError(t, err, "Error parsing config to yaml in testCase %s", testCase.name)
		err = fsutil.WriteToFile(asYAML, path)
		assert.NilError(t, err, "Error writing file in testCase %s", testCase.name)
	}

	err = (&RunCmd{
		GlobalFlags: &testCase.globalFlags,
	}).RunRun(nil, testCase.args)

	if testCase.expectedErr == "" {
		assert.NilError(t, err, "Unexpected error in testCase %s.", testCase.name)
	} else {
		assert.Error(t, err, testCase.expectedErr, "Wrong or no error in testCase %s.", testCase.name)
	}
}*/
