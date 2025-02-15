package hooks

import (
	"context"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/loft-sh/devspace/cmd"
	"github.com/loft-sh/devspace/cmd/flags"
	"github.com/loft-sh/devspace/e2e/framework"
	"github.com/loft-sh/devspace/e2e/kube"
	"github.com/loft-sh/devspace/pkg/util/factory"
	"github.com/onsi/ginkgo"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ = DevSpaceDescribe("hooks", func() {
	initialDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// create a new factory
	var (
		f          factory.Factory
		kubeClient *kube.KubeHelper
	)

	ginkgo.BeforeEach(func() {
		f = framework.NewDefaultFactory()

		kubeClient, err = kube.NewKubeHelper()
		framework.ExpectNoError(err)
	})

	ginkgo.It("should execute error hooks", func() {
		tempDir, err := framework.CopyToTempDir("tests/hooks/testdata/error")
		framework.ExpectNoError(err)
		defer framework.CleanupTempDir(initialDir, tempDir)

		ns, err := kubeClient.CreateNamespace("hooks")
		framework.ExpectNoError(err)
		defer func() {
			err := kubeClient.DeleteNamespace(ns)
			framework.ExpectNoError(err)
		}()

		buildCmd := &cmd.BuildCmd{
			GlobalFlags: &flags.GlobalFlags{
				NoWarn:    true,
				Namespace: ns,
			},
			SkipPush: false,
		}
		err = buildCmd.Run(f)
		framework.ExpectError(err)

		// check if files are there
		out, err := ioutil.ReadFile("before1.txt")
		framework.ExpectNoError(err)
		framework.ExpectEqual(string(out), "before1")
		framework.ExpectNoError(os.Remove("before1.txt"))
		out, err = ioutil.ReadFile("error1.txt")
		framework.ExpectNoError(err)
		framework.ExpectEqual(string(out), "error1")
		framework.ExpectNoError(os.Remove("error1.txt"))
		_, err = os.Stat("after1.txt")
		framework.ExpectError(err)
		_, err = os.Stat("before3.txt")
		framework.ExpectError(err)

		// now execute devspace dev and fail on deploy
		devCmd := &cmd.DevCmd{
			GlobalFlags: &flags.GlobalFlags{
				NoWarn:    true,
				Namespace: ns,
			},
			SkipPush: true,
		}
		err = devCmd.Run(f)
		framework.ExpectError(err)

		// check if files are correctly created
		out, err = ioutil.ReadFile("before1.txt")
		framework.ExpectNoError(err)
		framework.ExpectEqual(string(out), "before1")
		out, err = ioutil.ReadFile("after1.txt")
		framework.ExpectNoError(err)
		framework.ExpectEqual(string(out), "after1")
		out, err = ioutil.ReadFile("before2.txt")
		framework.ExpectNoError(err)
		framework.ExpectEqual(string(out), "before2")
		out, err = ioutil.ReadFile("error2.txt")
		framework.ExpectNoError(err)
		framework.ExpectEqual(string(out), "error2")
		out, err = ioutil.ReadFile("before3.txt")
		framework.ExpectNoError(err)
		framework.ExpectEqual(string(out), "before3")
		out, err = ioutil.ReadFile("error3.txt")
		framework.ExpectNoError(err)
		framework.ExpectEqual(string(out), "error3")
		out, err = ioutil.ReadFile("after3.txt")
		framework.ExpectNoError(err)
		framework.ExpectEqual(string(out), "after3")
		_, err = os.Stat("after2.txt")
		framework.ExpectError(err)
	})

	ginkgo.It("should execute hook once", func() {
		tempDir, err := framework.CopyToTempDir("tests/hooks/testdata/once")
		framework.ExpectNoError(err)
		defer framework.CleanupTempDir(initialDir, tempDir)

		ns, err := kubeClient.CreateNamespace("hooks")
		framework.ExpectNoError(err)
		defer func() {
			err := kubeClient.DeleteNamespace(ns)
			framework.ExpectNoError(err)
		}()

		// waitGroup for both commands
		waitGroup := sync.WaitGroup{}

		// create first dev command
		cancelCtx1, cancel1 := context.WithCancel(context.Background())
		defer cancel1()
		devCmd1 := &cmd.DevCmd{
			GlobalFlags: &flags.GlobalFlags{
				NoWarn:    true,
				Namespace: ns,
			},
			Ctx: cancelCtx1,
		}

		// start the command
		waitGroup.Add(1)
		go func() {
			defer ginkgo.GinkgoRecover()
			defer waitGroup.Done()
			err = devCmd1.Run(f)
			framework.ExpectNoError(err)
		}()

		// Read the 'once' hook output
		onceOutput1 := ""
		err = wait.PollImmediate(time.Second, time.Minute*2, func() (done bool, err error) {
			onceOutput1, err = kubeClient.ExecByImageSelector("node", ns, []string{"cat", "/app/once.out"})
			if err != nil {
				return false, nil
			}

			return onceOutput1 != "", nil
		})
		framework.ExpectNoError(err)

		// Read the 'each' hook output
		eachOutput1 := ""
		err = wait.PollImmediate(time.Second, time.Minute*2, func() (done bool, err error) {
			eachOutput1, err = kubeClient.ExecByImageSelector("node", ns, []string{"cat", "/app/each.out"})
			if err != nil {
				return false, nil
			}

			return eachOutput1 != "", nil
		})
		framework.ExpectNoError(err)

		// stop first command
		cancel1()

		// create second dev command
		cancelCtx2, cancel2 := context.WithCancel(context.Background())
		defer cancel2()
		devCmd2 := &cmd.DevCmd{
			GlobalFlags: &flags.GlobalFlags{
				NoWarn:    true,
				Namespace: ns,
			},
			Ctx: cancelCtx2,
		}

		// start the command
		waitGroup.Add(1)
		go func() {
			defer ginkgo.GinkgoRecover()
			defer waitGroup.Done()
			err = devCmd2.Run(f)

			framework.ExpectNoError(err)
		}()

		// Wait for 'each' hook output to change
		eachOutput2 := ""
		err = wait.PollImmediate(time.Second, time.Minute*2, func() (done bool, err error) {
			eachOutput2, err = kubeClient.ExecByImageSelector("node", ns, []string{"cat", "/app/each.out"})
			if err != nil {
				return false, nil
			}

			return eachOutput2 != eachOutput1, nil
		})
		framework.ExpectNoError(err)

		// Read the 'once' hook output again
		onceOutput2 := ""
		err = wait.PollImmediate(time.Second, time.Minute*2, func() (done bool, err error) {
			onceOutput2, err = kubeClient.ExecByImageSelector("node", ns, []string{"cat", "/app/once.out"})
			if err != nil {
				return false, nil
			}

			return onceOutput2 != "", nil
		})
		framework.ExpectNoError(err)

		// stop second command
		cancel2()

		// Verify that the 'once' hook did not run again
		framework.ExpectEqual(onceOutput1, onceOutput2)

		// wait for the command to finish
		waitGroup.Wait()
	})
})
