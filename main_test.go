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
	if out != "ulov24 "+Version || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestNoArgs(t *testing.T) {
	out, code := runCLI(t, nil)
	if !strings.HasPrefix(out, "Ishlatish: ulov24 <buyruq>") || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestPrice(t *testing.T) {
	cases := []struct {
		arg  string
		want string
		code int
	}{
		{"150000", "150 000 so'm", 0},
		{"abc", "Narx butun son bo'lishi kerak: abc", 1},
		{"-5", "Narx manfiy bo'lishi mumkin emas: -5", 1},
	}
	for _, c := range cases {
		out, code := runCLI(t, nil, "price", c.arg)
		if out != c.want || code != c.code {
			t.Errorf("price %s: got %q, code %d", c.arg, out, code)
		}
	}
}

func TestModes(t *testing.T) {
	out, code := runCLI(t, nil, "modes")
	want := "now — hozir\ntoday — bugun\nplanned — rejali"
	if out != want || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

// fixedNow is the clock every test runs at: 20 March 2026, 06:00 in Tashkent.
const fixedNow = "ULOV24_NOW=2026-03-20T06:00:00+05:00"

const wantSampleMatches = "So'rov #1 (2 kishi): yo'nalish #1\n" +
	"So'rov #2 (1 kishi): yo'nalish #2\n" +
	"So'rov #3 (5 kishi): mos yo'nalish yo'q\n" +
	"So'rov #4 (1 kishi): mos yo'nalish yo'q"

func TestMatchSample(t *testing.T) {
	out, code := runCLI(t, []string{fixedNow}, "match")
	if out != wantSampleMatches || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestMatchIgnoresMachineZone(t *testing.T) {
	out, code := runCLI(t, []string{fixedNow, "TZ=America/New_York"}, "match")
	if out != wantSampleMatches || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestRoutesSample(t *testing.T) {
	out, code := runCLI(t, []string{fixedNow, "TZ=America/New_York"}, "routes")
	lines := strings.Split(out, "\n")
	if code != 0 || len(lines) != 4 {
		t.Fatalf("got %q, code %d", out, code)
	}
	if want := "#1 Toshkent → Samarqand · 20.03 08:00 · 3 joy · 150 000 so'm · bugun"; lines[0] != want {
		t.Errorf("line 1 = %q, want %q", lines[0], want)
	}
	if !strings.HasSuffix(lines[3], "· eskirgan") {
		t.Errorf("line 4 = %q, want it to end with · eskirgan", lines[3])
	}
}

const wantSampleFeed = "#2 Toshkent → Samarqand · 20.03 06:30 · 1 joy · 140 000 so'm · hozir\n" +
	"#1 Toshkent → Samarqand · 20.03 08:00 · 3 joy · 150 000 so'm · bugun\n" +
	"#3 Toshkent → Buxoro · 21.03 08:00 · 4 joy · 200 000 so'm · rejali"

func TestFeedSample(t *testing.T) {
	out, code := runCLI(t, []string{fixedNow}, "feed")
	if out != wantSampleFeed || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestFeedOneDirection(t *testing.T) {
	out, code := runCLI(t, []string{fixedNow}, "feed", "toshkent", "samarqand")
	want := strings.Join(strings.Split(wantSampleFeed, "\n")[:2], "\n")
	if out != want || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestFeedUnknownCity(t *testing.T) {
	out, code := runCLI(t, []string{fixedNow}, "feed", "toshkent", "moskva")
	if out != "Noma'lum shahar: moskva" || code != 1 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestCities(t *testing.T) {
	out, code := runCLI(t, nil, "cities")
	want := "1 Toshkent\n2 Samarqand\n3 Buxoro\n4 Farg'ona\n5 Andijon\n6 Namangan"
	if out != want || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestBadClock(t *testing.T) {
	out, code := runCLI(t, []string{"ULOV24_NOW=ertaga"}, "match")
	if out != "ULOV24_NOW noto'g'ri: ertaga" || code != 1 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestUnknownCommand(t *testing.T) {
	out, code := runCLI(t, nil, "foo")
	if out != "Noma'lum buyruq: foo" || code != 1 {
		t.Fatalf("got %q, code %d", out, code)
	}
}
