package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// deadPid is a PID no process holds, so the parent wait returns at once.
const deadPid = "99999999"

// TestMain runs main() in a child process when UPDATER_TEST_ARGS is set,
// recording each relaunch to UPDATER_TEST_LAUNCHED instead of running `open`.
func TestMain(m *testing.M) {
	if args := os.Getenv("UPDATER_TEST_ARGS"); args != "" {
		marker := os.Getenv("UPDATER_TEST_LAUNCHED")
		launchApp = func(appPath string) {
			f, err := os.OpenFile(marker, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			f.WriteString(appPath + "\n")
			f.Close()
		}
		os.Args = append([]string{"updater"}, strings.Split(args, "\n")...)
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// runUpdater runs main() with args in a child process and returns the app
// paths it relaunched.
func runUpdater(t *testing.T, args ...string) []string {
	t.Helper()
	marker := filepath.Join(t.TempDir(), "launched")
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	cmd.Env = append(os.Environ(),
		"UPDATER_TEST_ARGS="+strings.Join(args, "\n"),
		"UPDATER_TEST_LAUNCHED="+marker,
	)
	out, err := cmd.CombinedOutput()
	if _, ok := err.(*exec.ExitError); err != nil && !ok {
		t.Fatal(err)
	}
	t.Logf("updater output:\n%s", out)

	data, err := os.ReadFile(marker)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
}

// unwritableLogPath returns a log path whose parent is a regular file.
func unwritableLogPath(t *testing.T) string {
	t.Helper()
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, nil, 0644); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(blocker, "updater.log")
}

func TestRestartRelaunchesWhenLogCannotBeOpened(t *testing.T) {
	appPath := filepath.Join(t.TempDir(), "ACRA Point.app")

	launched := runUpdater(t, "--restart", deadPid, unwritableLogPath(t), appPath)

	if len(launched) != 1 || launched[0] != appPath {
		t.Fatalf("relaunched %q, want exactly [%q]", launched, appPath)
	}
}

func TestUpdateSwapsAndRelaunchesWhenLogCannotBeOpened(t *testing.T) {
	root := t.TempDir()
	appPath := filepath.Join(root, "ACRA Point.app")
	extracted := filepath.Join(root, "extracted")
	backupPath := filepath.Join(root, "backup.app")
	if err := os.MkdirAll(appPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(appPath, "version"), []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	newApp := filepath.Join(extracted, "ACRA Point.app")
	if err := os.MkdirAll(newApp, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(newApp, "version"), []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}

	launched := runUpdater(t, "--update", deadPid, unwritableLogPath(t), appPath, extracted, backupPath)

	if len(launched) != 1 || launched[0] != appPath {
		t.Fatalf("relaunched %q, want exactly [%q]", launched, appPath)
	}
	got, err := os.ReadFile(filepath.Join(appPath, "version"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "new" {
		t.Fatalf("app at app_path is %q, want the new bundle", got)
	}
}
