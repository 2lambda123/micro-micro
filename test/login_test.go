// +build kind

package test

import (
	"strings"
	"testing"
	"time"

	"github.com/micro/micro/v3/client/cli/token"
)

func TestCorruptedTokenLogin(t *testing.T) {
	TrySuite(t, testCorruptedLogin, retryCount)
}

func testCorruptedLogin(t *T) {
	serv := NewServer(t)
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	t.Parallel()

	// get server command
	cmd := serv.Command()

	outp, _ := cmd.Exec("status")
	if !strings.Contains(string(outp), "Unauthorized") {
		t.Fatalf("Call should need authorization")
	}
	outp, _ = cmd.Exec("login", "--email", serv.Env(), "--password", "password")
	if !strings.Contains(string(outp), "Successfully logged in.") {
		t.Fatalf("Login failed: %s", outp)
	}
	outp, _ = cmd.Exec("status")
	if string(outp) != "" {
		t.Fatalf("Call should receive no output: %s", outp)
	}
	// munge token
	tok, _ := token.Get(serv.Env())
	tok.Expiry = time.Now().Add(-1 * time.Hour)
	tok.RefreshToken = tok.RefreshToken + "a"
	token.Save(serv.Env(), tok)

	outp, _ = cmd.Exec("status")
	if !strings.Contains(string(outp), "Account can't be found for refresh token") {
		t.Fatalf("Call should have failed: %s", outp)
	}
	outp, _ = cmd.Exec("login", "--email", serv.Env(), "--password", "password")
	if !strings.Contains(string(outp), "Successfully logged in.") {
		t.Fatalf("Login failed: %s", outp)
	}

}
