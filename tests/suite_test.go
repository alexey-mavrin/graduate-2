package integration_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/rand"
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

	Describe("Basic usage", func() {
		It("Should successfully register user", func() {
			By("Running 'client user'")
			stdOut, stdErr, err := runClient("user -a register")
			fmt.Print(stdOut, stdErr)
			Expect(err).NotTo(HaveOccurred(), "Client should register")
		})

		It("Should store and retrieve account", func() {
			By("Running 'client acc'")
			_, _, err := runClient("user -a register")
			Expect(err).NotTo(HaveOccurred(), "Client should register")

			stdOut, stdErr, err := runClient(
				"acc -a store -n acc_name_1, -u u1 -p pass1",
			)
			fmt.Print(stdOut, stdErr)
			Expect(err).NotTo(HaveOccurred(), "Client should run OK")
			stdOut, stdErr, err = runClient(
				"acc -a get -i 1",
			)
			fmt.Print(stdOut, stdErr)
			Expect(err).NotTo(HaveOccurred(), "Client should run OK")
		})

		It("Should store and retrieve note", func() {
			By("Running 'client note'")
			_, _, err := runClient("user -a register")
			Expect(err).NotTo(HaveOccurred(), "Client should register")

			stdOut, stdErr, err := runClient(
				"note -a store -n note_name_1, -t text1",
			)
			fmt.Print(stdOut, stdErr)
			Expect(err).NotTo(HaveOccurred(), "Client should run OK")
			stdOut, stdErr, err = runClient(
				"note -a get -i 1",
			)
			fmt.Print(stdOut, stdErr)
			Expect(err).NotTo(HaveOccurred(), "Client should run OK")
		})

		It("Should store and retrieve card", func() {
			By("Running 'client card'")
			_, _, err := runClient("user -a register")
			Expect(err).NotTo(HaveOccurred(), "Client should register")

			stdOut, stdErr, err := runClient(
				"card -a store -n c1 -num 1111222233334444 -c 123",
			)
			fmt.Print(stdOut, stdErr)
			Expect(err).NotTo(HaveOccurred(), "Client should run OK")
			stdOut, stdErr, err = runClient(
				"card -a get -i 1",
			)
			fmt.Print(stdOut, stdErr)
			Expect(err).NotTo(HaveOccurred(), "Client should run OK")
		})

		It("Should store and retrieve binary data", func() {
			By("Running 'client bin'")
			_, _, err := runClient("user -a register")
			Expect(err).NotTo(HaveOccurred(), "Client should register")

			file, err := makeBinFile(256)
			Expect(err).NotTo(HaveOccurred(), "Binary file is created")
			stdOut, stdErr, err := runClient(
				fmt.Sprintf("bin -a store -n bin1 -f %s", file),
			)
			fmt.Print(stdOut, stdErr)

			fileOut := file + ".out"
			Expect(err).NotTo(HaveOccurred(), "Client should run OK")
			stdOut, stdErr, err = runClient(
				fmt.Sprintf("bin -a get -i 1 -f %s", fileOut),
			)
			fmt.Print(stdOut, stdErr)
			Expect(err).NotTo(HaveOccurred(), "Client should run OK")

			h1, err := getHash(file)
			Expect(err).NotTo(HaveOccurred(), "Get original file hash")
			h2, err := getHash(fileOut)
			Expect(err).NotTo(HaveOccurred(), "Get saved file hash")
			Expect(h1).To(Equal(h2), "Files should be the same")
			os.Remove(file)
			os.Remove(fileOut)
		})
	})
})

func makeBinFile(size uint) (string, error) {
	file, err := ioutil.TempFile("/tmp", "bin")
	if err != nil {
		return "", err
	}
	_, err = file.Write(randomBin(size))
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func randomBin(n uint) []byte {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, n)
	for i := range buf {
		buf[i] = byte(rand.Intn(256))
	}
	return buf
}

func getHash(file string) ([sha256.Size]byte, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return [sha256.Size]byte{}, err
	}
	hash := sha256.Sum256(buf)
	return hash, nil
}
