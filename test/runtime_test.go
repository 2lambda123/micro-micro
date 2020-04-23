// +build integration

package test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

type cmdFunc func() ([]byte, error)

func try(blockName string, t *testing.T, f cmdFunc, maxTime time.Duration) {
	elapsed := 0 * time.Millisecond
	var outp []byte
	var err error
	for {
		if elapsed > maxTime {
			if err != nil {
				t.Fatalf("%v (failed after %v with '%v'), output: %v", blockName, elapsed, err, string(outp))
			}
		}
		outp, err = f()
		if err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
		elapsed += 100 * time.Millisecond
	}
}

type server struct {
	cmd *exec.Cmd
	t   *testing.T
}

func newServer(t *testing.T) server {
	// @todo this is a dangerous move, should instead specify a branch new
	// folder for tests and only nuke those
	outp, err := exec.Command("rm", "-rf", "/tmp/micro/store").CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}
	return server{cmd: exec.Command("micro", "server"), t: t}
}

func (s server) launch() {
	go func() {
		if err := s.cmd.Start(); err != nil {
			s.t.Fatal(err)
		}
	}()
	try("Calling micro server", s.t, func() ([]byte, error) {
		return exec.Command("micro", "call", "go.micro.runtime", "Runtime.Read", "{}").CombinedOutput()
	}, 3000*time.Millisecond)
}

func (s server) close() {
	if s.cmd.Process != nil {
		s.cmd.Process.Signal(syscall.SIGTERM)
	}
}

func TestNew(t *testing.T) {
	defer func() {
		exec.Command("rm", "-r", "./foobar").CombinedOutput()
	}()
	outp, err := exec.Command("micro", "new", "foobar").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(outp), "protoc") {
		t.Fatalf("micro new lacks protobuf install instructions %v", string(outp))
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
	outp, err := exec.Command("micro", "env", "set", "server").CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to set env to server, err: %v, output: %v", err, string(outp))
	}

	callCmd := exec.Command("micro", "call", "go.micro.runtime", "Runtime.Read", "{}")
	outp, err = callCmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Call to server should fail, got no error, output: %v", string(outp))
	}

	serv := newServer(t)
	serv.launch()
	defer serv.close()

	try("Calling Runtime.Read", t, func() ([]byte, error) {
		outp, err = exec.Command("micro", "call", "go.micro.runtime", "Runtime.Read", "{}").CombinedOutput()
		if err != nil {
			return outp, errors.New("Call to runtime read should succeed")
		}
		return outp, err
	}, 2*time.Second)

}

func TestRunLocalSource(t *testing.T) {
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", "run", "./example-service")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	try("Find test/example", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", "ps")
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
	}, 30*time.Second)

	try("Find go.micro.service.example in list", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", "list", "services").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "go.micro.service.example") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 5*time.Second)
}

func TestRunGithubSource(t *testing.T) {
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

	runCmd := exec.Command("micro", "run", "helloworld")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	try("Find hello world", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", "ps")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "helloworld") {
			return outp, errors.New("Output should contain hello world")
		}
		return outp, nil
	}, 30*time.Second)

	try("Call hello world", t, func() ([]byte, error) {
		callCmd := exec.Command("micro", "call", "go.micro.service.helloworld", "Helloworld.Call", `{"name": "Joe"}`)
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
	}, 15*time.Second)

}

func TestRunLocalUpdateAndCall(t *testing.T) {
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	// Run the example service
	runCmd := exec.Command("micro", "run", "./example-service")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
	}

	try("Finding example service with micro ps", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", "ps")
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
		callCmd := exec.Command("micro", "call", "go.micro.service.example", "Example.Call", `{"name": "Joe"}`)
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
	}, 15*time.Second)

	replaceStringInFile(t, "./example-service/handler/handler.go", "Hello", "Hi")
	defer func() {
		// Change file back
		replaceStringInFile(t, "./example-service/handler/handler.go", "Hi", "Hello")
	}()

	updateCmd := exec.Command("micro", "update", "./example-service")
	outp, err = updateCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	try("Call example service after modification", t, func() ([]byte, error) {
		callCmd := exec.Command("micro", "call", "go.micro.service.example", "Example.Call", `{"name": "Joe"}`)
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
	}, 15*time.Second)
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
