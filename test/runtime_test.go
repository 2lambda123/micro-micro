// +build integration

package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

type cmdFunc func() ([]byte, error)

func try(blockName string, t *testing.T, f cmdFunc, maxTime time.Duration) {
	start := time.Now()
	var outp []byte
	var err error

	for {
		if time.Since(start) > maxTime {
			_, file, line, _ := runtime.Caller(1)
			fname := filepath.Base(file)
			if err != nil {
				t.Fatalf("%v:%v, %v (failed after %v with '%v'), output: '%v'", fname, line, blockName, time.Since(start), err, string(outp))
			}
			return
		}
		outp, err = f()
		if err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func once(blockName string, t *testing.T, f cmdFunc) {
	outp, err := f()
	if err != nil {
		t.Fatalf("%v with '%v', output: %v", blockName, err, string(outp))
	}
}

type server struct {
	cmd       *exec.Cmd
	t         *testing.T
	proxyPort int
}

func newServer(t *testing.T) server {
	min := 8000
	max := 60000
	portnum := rand.Intn(max-min) + min

	return server{
		cmd: exec.Command("docker", "run",
			fmt.Sprintf("-p=%v:8081", portnum), "micro", "server"),
		t:         t,
		proxyPort: portnum,
	}
}

func (s server) launch() {
	go func() {
		if err := s.cmd.Start(); err != nil {
			s.t.Fatal(err)
		}
	}()
	// @todo find a way to know everything is up and running
	try("Calling micro server", s.t, func() ([]byte, error) {
		return exec.Command("micro", s.envFlag(), "call", "go.micro.runtime", "Runtime.Read", "{}").CombinedOutput()
	}, 10000*time.Millisecond)
}

func (s server) close() {
	if s.cmd.Process != nil {
		s.cmd.Process.Signal(syscall.SIGTERM)
	}
}

func (s server) envFlag() string {
	return fmt.Sprintf("-env=127.0.0.1:%v", s.proxyPort)
}

func TestNew(t *testing.T) {
	t.Parallel()
	defer func() {
		exec.Command("rm", "-r", "./foobar").CombinedOutput()
	}()
	outp, err := exec.Command("micro", "new", "foobar").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(outp), "protoc") {
		t.Fatalf("micro new lacks 	protobuf install instructions %v", string(outp))
	}

	lines := strings.Split(string(outp), "\n")
	// executing install instructions
	for _, line := range lines {
		if strings.HasPrefix(line, "go get") {
			parts := strings.Split(line, " ")
			getOutp, getErr := exec.Command(parts[0], parts[1:]...).CombinedOutput()
			if getErr != nil {
				t.Fatal(string(getOutp))
			}
		}
		if strings.HasPrefix(line, "protoc") {
			parts := strings.Split(line, " ")
			protocCmd := exec.Command(parts[0], parts[1:]...)
			protocCmd.Dir = "./foobar"
			pOutp, pErr := protocCmd.CombinedOutput()
			if pErr != nil {
				t.Fatal(string(pOutp))
			}
		}
	}

	buildCommand := exec.Command("go", "build")
	buildCommand.Dir = "./foobar"
	outp, err = buildCommand.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}
}

func TestServerModeCall(t *testing.T) {
	t.Parallel()
	serv := newServer(t)

	callCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.runtime", "Runtime.Read", "{}")
	outp, err := callCmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Call to server should fail, got no error, output: %v", string(outp))
	}

	serv.launch()
	defer serv.close()

	try("Calling Runtime.Read", t, func() ([]byte, error) {
		outp, err = exec.Command("micro", serv.envFlag(), "call", "go.micro.runtime", "Runtime.Read", "{}").CombinedOutput()
		if err != nil {
			return outp, errors.New("Call to runtime read should succeed")
		}
		return outp, err
	}, 2*time.Second)
}

func TestRunLocalSource(t *testing.T) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", serv.envFlag(), "run", "./example-service")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	try("Find test/example", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "test/example-service",
		// as the runtime name is the relative path inside a repo.
		if !strings.Contains(string(outp), "test/example-service") {
			return outp, errors.New("Can't find example service in runtime")
		}
		return outp, err
	}, 12*time.Second)

	try("Find go.micro.service.example in list", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "list", "services").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "go.micro.service.example") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 5*time.Second)
}

func TestLocalOutsideRepo(t *testing.T) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	dirname := "last-dir-of-path"
	folderPath := filepath.Join(os.TempDir(), dirname)

	err := os.MkdirAll(folderPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	// since copying a whole folder is rather involved and only Linux sources
	// are available, see https://stackoverflow.com/questions/51779243/copy-a-folder-in-go
	// we fall back to `cp`
	outp, err := exec.Command("cp", "-r", "example-service/.", folderPath).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	runCmd := exec.Command("micro", serv.envFlag(), "run", ".")
	runCmd.Dir = folderPath
	outp, err = runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	try("Find "+dirname, t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		lines := strings.Split(string(outp), "\n")
		found := false
		for _, line := range lines {
			if strings.HasPrefix(line, dirname) {
				found = true
			}
		}
		if !found {
			return outp, errors.New("Can't find '" + dirname + "' in runtime")
		}
		return outp, err
	}, 12*time.Second)

	try("Find go.micro.service.example in list", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "list", "services").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "go.micro.service.example") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 12*time.Second)
}

func TestLocalEnvRunGithubSource(t *testing.T) {
	t.Parallel()
	outp, err := exec.Command("micro", "env", "set", "local").CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to set env to local, err: %v, output: %v", err, string(outp))
	}
	var cmd *exec.Cmd
	go func() {
		cmd = exec.Command("micro", "run", "location")
		// fire and forget as this will run forever
		cmd.CombinedOutput()
	}()
	time.Sleep(100 * time.Millisecond)
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}()

	try("Find location", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", "list", "services")
		outp, err := psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "location") {
			return outp, errors.New("Output should contain location")
		}
		return outp, nil
	}, 20*time.Second)
}

func TestRunGithubSource(t *testing.T) {
	t.Parallel()
	p, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	if len(p) == 0 {
		t.Fatalf("Git is not available %v", p)
	}
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", serv.envFlag(), "run", "helloworld")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	try("Find hello world", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "helloworld") {
			return outp, errors.New("Output should contain hello world")
		}
		return outp, nil
	}, 20*time.Second)

	try("Call hello world", t, func() ([]byte, error) {
		callCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.service.helloworld", "Helloworld.Call", `{"name": "Joe"}`)
		outp, err := callCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		rsp := map[string]string{}
		err = json.Unmarshal(outp, &rsp)
		if err != nil {
			return outp, err
		}
		if rsp["msg"] != "Hello Joe" {
			return outp, errors.New("Helloworld resonse is unexpected")
		}
		return outp, err
	}, 30*time.Second)

}

func TestRunLocalUpdateAndCall(t *testing.T) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	// Run the example service
	runCmd := exec.Command("micro", serv.envFlag(), "run", "./example-service")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	try("Finding example service with micro status", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "test/example-service",
		// as the runtime name is the relative path inside a repo.
		if !strings.Contains(string(outp), "test/example-service") {
			return outp, errors.New("can't find service in runtime")
		}
		return outp, err
	}, 30*time.Second)

	try("Call example service", t, func() ([]byte, error) {
		callCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.service.example", "Example.Call", `{"name": "Joe"}`)
		outp, err := callCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		rsp := map[string]string{}
		err = json.Unmarshal(outp, &rsp)
		if err != nil {
			return outp, err
		}
		if rsp["msg"] != "Hello Joe" {
			return outp, errors.New("Resonse is unexpected")
		}
		return outp, err
	}, 8*time.Second)

	replaceStringInFile(t, "./example-service/handler/handler.go", "Hello", "Hi")
	defer func() {
		// Change file back
		replaceStringInFile(t, "./example-service/handler/handler.go", "Hi", "Hello")
	}()

	updateCmd := exec.Command("micro", serv.envFlag(), "update", "./example-service")
	outp, err = updateCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	try("Call example service after modification", t, func() ([]byte, error) {
		callCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.service.example", "Example.Call", `{"name": "Joe"}`)
		outp, err = callCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		rsp := map[string]string{}
		err = json.Unmarshal(outp, &rsp)
		if err != nil {
			return outp, err
		}
		if rsp["msg"] != "Hi Joe" {
			return outp, errors.New("Response is not what's expected")
		}
		return outp, err
	}, 8*time.Second)
}

func TestExistingLogs(t *testing.T) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", serv.envFlag(), "run", "github.com/crufter/micro-services/logspammer")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	try("logspammer logs", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "logs", "-n", "5", "crufter/micro-services/logspammer")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "Listening on") || !strings.Contains(string(outp), "never stopping") {
			return outp, errors.New("Output does not contain expected")
		}
		return outp, nil
	}, 25*time.Second)
}

func TestStreamLogsAndThirdPartyRepo(t *testing.T) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", serv.envFlag(), "run", "github.com/crufter/micro-services/logspammer")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	try("logspammer logs", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "logs", "-n", "5", "crufter/micro-services/logspammer")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "Listening on") || !strings.Contains(string(outp), "never stopping") {
			return outp, errors.New("Output does not contain expected")
		}
		return outp, nil
	}, 25*time.Second)

	// Test streaming logs
	cmd := exec.Command("micro", serv.envFlag(), "logs", "-n", "1", "-f", "crufter-micro-services-logspammer")

	go func() {
		outp, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(err)
		}
		if len(outp) == 0 {
			t.Fatal("No log lines streamed")
		}
		if !strings.Contains(string(outp), "never stopping") {
			t.Fatalf("Unexpected logs: %v", string(outp))
		}
		// Logspammer logs every 2 seconds, so we need 2 different
		now := time.Now()
		// leaving the hour here to fix a docker issue
		// when the containers clock is a few hours behind
		stampA := now.Add(-2 * time.Second).Format("04:05")
		stampB := now.Add(-1 * time.Second).Format("104:05")
		if !strings.Contains(string(outp), stampA) && !strings.Contains(string(outp), stampB) {
			t.Fatalf("Timestamp %v or %v not found in logs: %v", stampA, stampB, string(outp))
		}
	}()

	time.Sleep(6 * time.Second)
	err = cmd.Process.Kill()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
}

func replaceStringInFile(t *testing.T, filepath string, original, newone string) {
	input, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	output := strings.ReplaceAll(string(input), original, newone)
	err = ioutil.WriteFile(filepath, []byte(output), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
