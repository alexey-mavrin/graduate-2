package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Suite test")
}

var _ = BeforeSuite(func() {
	runCmd("go build -o server ../cmd/server/main.go", nil)
})

func cleanUp() {
	By("removing server storage and client cache")
	os.Remove("cache_inttest_store.db")
	os.Remove("server_inttest_storage.db")
}

// runCmd get command line and returns command stdout, stderr and error
func runCmd(cmdLine string, env []string) (string, string, error) {
	cmdArgs := strings.Fields(cmdLine)
	fmt.Println("running", cmdArgs)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)
	var o, e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e

	err := cmd.Run()
	return o.String(), e.String(), err
}

// startCmd startd specified process in the background
func startCmd(cmdLine string, env []string) (*exec.Cmd, error) {
	cmdArgs := strings.Fields(cmdLine)
	fmt.Println("running", cmdArgs)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)
	err := cmd.Start()
	return cmd, err
}

// stopCmd kills the previously started process
func stopCmd(cmd *exec.Cmd) error {
	fmt.Printf("------------ %v\n", cmd.Process)
	return cmd.Process.Kill()
}

func startServer() (*exec.Cmd, error) {
	env := []string{
		"SERVER_CFG=server_inttest.cfg",
	}
	cmdLine := "./server"
	return startCmd(cmdLine, env)
}

func runClient(args string) (string, string, error) {
	env := []string{
		"GOSECRET_CFG=gosecret_inttest.cfg",
	}
	cmdLine := "go run ../cmd/client/main.go " + args
	return runCmd(cmdLine, env)
}

var _ = Describe("Run server and client together", func() {
	var srvCmd *exec.Cmd
	var cmdErr error

	BeforeEach(func() {
		cleanUp()
		srvCmd, cmdErr = startServer()
		Expect(cmdErr).NotTo(HaveOccurred(), "Server should start")
		time.Sleep(time.Second * 1)
	})

	AfterEach(func() {
		cmdErr = stopCmd(srvCmd)
		Expect(cmdErr).NotTo(HaveOccurred(), "Server should stop")
	})

	Describe("Register user", func() {
		It("Should successfully register", func() {
			By("Running client with 'user -a register'")
			stdOut, stdErr, err := runClient("user -a register")
			Expect(err).NotTo(HaveOccurred(), "Client should run OK")
			fmt.Print(stdOut, stdErr)
		})
	})
})
