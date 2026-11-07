package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// bin is the ulov24 program, built once for all tests in this package.
var bin string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "ulov24-test-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	bin = filepath.Join(dir, "ulov24")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		os.RemoveAll(dir)
		fmt.Fprintln(os.Stderr, "build failed:", err)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// runCLI runs the program with the given extra environment and arguments and
// returns its trimmed standard output and exit code.
func runCLI(t *testing.T, env []string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), env...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	code := 0
	if err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			t.Fatalf("run %v: %v", args, err)
		}
		code = exit.ExitCode()
	}
	return strings.TrimSpace(out.String()), code
}

func TestVersion(t *testing.T) {
	out, code := runCLI(t, nil, "version")
	if out != "ulov24 c1-m1" || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestNoArgs(t *testing.T) {
	out, code := runCLI(t, nil)
	if !strings.HasPrefix(out, "Ishlatish: ulov24 <buyruq>") || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestUnknownCommand(t *testing.T) {
	out, code := runCLI(t, nil, "foo")
	if out != "Noma'lum buyruq: foo" || code != 1 {
		t.Fatalf("got %q, code %d", out, code)
	}
}
