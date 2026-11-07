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

func TestUnknownCommand(t *testing.T) {
	out, code := runCLI(t, nil, "foo")
	if out != "Noma'lum buyruq: foo" || code != 1 {
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

func TestCities(t *testing.T) {
	out, code := runCLI(t, nil, "cities")
	want := "1 Toshkent\n2 Samarqand\n3 Buxoro\n4 Farg'ona\n5 Andijon\n6 Namangan"
	if out != want || code != 0 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

// fixedNow is the clock every test runs at: 20 March 2026, 06:00 in Tashkent.
const fixedNow = "ULOV24_NOW=2026-03-20T06:00:00+05:00"

// dataEnv gives a test its own empty data file and the fixed clock.
func dataEnv(t *testing.T) ([]string, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "d.json")
	return []string{fixedNow, "ULOV24_DATA=" + path}, path
}

// mustRun runs a command that must succeed and print want.
func mustRun(t *testing.T, env []string, want string, args ...string) {
	t.Helper()
	out, code := runCLI(t, env, args...)
	if out != want || code != 0 {
		t.Fatalf("%v: got %q, code %d; want %q", args, out, code, want)
	}
}

func TestBookingFlow(t *testing.T) {
	env, _ := dataEnv(t)
	mustRun(t, env, "Yo'nalish #1 qo'shildi",
		"route", "add", "toshkent", "samarqand", "2026-03-20 08:00", "3", "today", "150000")
	mustRun(t, env, "So'rov #1 qo'shildi",
		"request", "add", "toshkent", "samarqand", "2026-03-20 07:00", "2026-03-20 09:00", "2", "today")
	mustRun(t, env, "Moslik #1: so'rov #1 → yo'nalish #1", "match")
	mustRun(t, env, "Yangi moslik yo'q", "match")
	mustRun(t, env, "Moslik #1: partial", "accept", "1", "passenger")
	mustRun(t, env, "Moslik #1: accepted · yo'nalishda 1 joy qoldi", "accept", "1", "driver")

	out, code := runCLI(t, env, "accept", "1", "driver")
	if out != "Moslik #1 allaqachon yopilgan" || code != 1 {
		t.Fatalf("accept twice: got %q, code %d", out, code)
	}
}

const wantFeed = "#2 Toshkent → Samarqand · 20.03 06:30 · 1 joy · 140 000 so'm · hozir\n" +
	"#1 Toshkent → Samarqand · 20.03 08:00 · 3 joy · 150 000 so'm · bugun\n" +
	"#3 Toshkent → Buxoro · 21.03 08:00 · 4 joy · 200 000 so'm · rejali"

func addFeedRoutes(t *testing.T, env []string) {
	t.Helper()
	mustRun(t, env, "Yo'nalish #1 qo'shildi", "route", "add", "toshkent", "samarqand", "2026-03-20 08:00", "3", "today", "150000")
	mustRun(t, env, "Yo'nalish #2 qo'shildi", "route", "add", "Toshkent", "Samarqand", "2026-03-20 06:30", "1", "now", "140000")
	mustRun(t, env, "Yo'nalish #3 qo'shildi", "route", "add", "toshkent", "buxoro", "2026-03-21 08:00", "4", "planned", "200000")
}

func TestFeedFromFile(t *testing.T) {
	env, _ := dataEnv(t)
	addFeedRoutes(t, env)
	mustRun(t, env, wantFeed+"\nKo'rsatildi: 3 / 3", "feed")
	mustRun(t, append(env, "TZ=America/New_York"), wantFeed+"\nKo'rsatildi: 3 / 3", "feed")
	mustRun(t, env, strings.Join(strings.Split(wantFeed, "\n")[:2], "\n")+"\nKo'rsatildi: 2 / 2", "feed", "toshkent", "samarqand")
}

func TestRouteWithoutPrice(t *testing.T) {
	env, _ := dataEnv(t)
	mustRun(t, env, "Yo'nalish #1 qo'shildi", "route", "add", "toshkent", "samarqand", "2026-03-20 08:00", "3", "today")
	mustRun(t, env, "#1 Toshkent → Samarqand · 20.03 08:00 · 3 joy · narx kelishiladi · bugun\nKo'rsatildi: 1 / 1", "feed")
}

func TestInputErrors(t *testing.T) {
	env, path := dataEnv(t)
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"route", "add", "toshkent", "samarqand", "2026-03-20 08:00", "9", "today"}, "joylar soni 1 dan 8 gacha bo'lishi kerak: 9"},
		{[]string{"route", "add", "toshkent", "samarqand", "2026-03-20 08:00", "3", "today", "-1"}, "Narx manfiy bo'lishi mumkin emas: -1"},
		{[]string{"route", "add", "toshkent", "samarqand", "20.03 08:00", "3", "today"}, "Vaqt formati: YYYY-MM-DD HH:MM"},
		{[]string{"route", "add", "toshkent", "moskva", "2026-03-20 08:00", "3", "today"}, "Noma'lum shahar: moskva"},
		{[]string{"feed", "toshkent", "moskva"}, "Noma'lum shahar: moskva"},
		{[]string{"route", "add", "toshkent", "toshkent", "2026-03-20 08:00", "3", "today"}, "Qayerdan va qayerga bir xil bo'lmasin"},
		{[]string{"route", "add", "toshkent", "samarqand", "2026-03-20 08:00", "3", "tez"}, "noma'lum rejim: tez"},
		{[]string{"request", "add", "toshkent", "samarqand", "2026-03-20 09:00", "2026-03-20 07:00", "2", "today"}, "Vaqt oralig'i noto'g'ri: boshlanishi oxiridan keyin"},
		{[]string{"accept", "7", "driver"}, "Bunday moslik yo'q: 7"},
		{[]string{"accept", "1", "boss"}, "Rol driver yoki passenger bo'lishi kerak: boss"},
	}
	for _, c := range cases {
		out, code := runCLI(t, env, c.args...)
		if out != c.want || code != 1 {
			t.Errorf("%v: got %q, code %d; want %q", c.args, out, code, c.want)
		}
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("failed commands must not create the data file: %v", err)
	}
}

func TestCorruptDataFile(t *testing.T) {
	env, path := dataEnv(t)
	if err := os.WriteFile(path, []byte("{bad"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"feed"}, {"route", "add", "toshkent", "samarqand", "2026-03-20 08:00", "3", "today"}} {
		out, code := runCLI(t, env, args...)
		if out != "Ma'lumot fayli buzilgan: "+path || code != 1 {
			t.Errorf("%v: got %q, code %d", args, out, code)
		}
	}
	if data, _ := os.ReadFile(path); string(data) != "{bad" {
		t.Fatalf("data file changed: %q", data)
	}
}

func TestSameUserDoesNotMatch(t *testing.T) {
	env, _ := dataEnv(t)
	me := append(env, "ULOV24_USER=7")
	mustRun(t, me, "Yo'nalish #1 qo'shildi", "route", "add", "toshkent", "samarqand", "2026-03-20 08:00", "3", "today")
	mustRun(t, me, "So'rov #1 qo'shildi", "request", "add", "toshkent", "samarqand", "2026-03-20 07:00", "2026-03-20 09:00", "1", "today")
	mustRun(t, env, "Yangi moslik yo'q", "match")
}

func TestBadClock(t *testing.T) {
	out, code := runCLI(t, []string{"ULOV24_NOW=ertaga"}, "feed")
	if out != "ULOV24_NOW noto'g'ri: ertaga" || code != 1 {
		t.Fatalf("got %q, code %d", out, code)
	}
}

func TestFeedFirstPage(t *testing.T) {
	env, _ := dataEnv(t)
	for i := 1; i <= 12; i++ {
		mustRun(t, env, fmt.Sprintf("Yo'nalish #%d qo'shildi", i),
			"route", "add", "toshkent", "samarqand", fmt.Sprintf("2026-03-21 %02d:00", i), "3", "planned")
	}
	out, code := runCLI(t, env, "feed")
	lines := strings.Split(out, "\n")
	if code != 0 || len(lines) != 11 || lines[10] != "Ko'rsatildi: 10 / 12" {
		t.Fatalf("got %d lines, code %d, last %q", len(lines), code, lines[len(lines)-1])
	}
}

func TestStoreLog(t *testing.T) {
	env, _ := dataEnv(t)
	cmd := exec.Command(bin, "route", "add", "toshkent", "samarqand", "2026-03-20 08:00", "3", "today")
	cmd.Env = append(append(os.Environ(), env...), "ULOV24_LOG=1")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr.String(), "store: SaveRoute #1") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
