package monitor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kihyun1998/macos_zip_updater/logger"
)

// exitedPid returns the PID of a child that has exited and been reaped.
func exitedPid(t *testing.T) int {
	t.Helper()
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	return cmd.Process.Pid
}

func TestIsProcessAlive(t *testing.T) {
	cases := []struct {
		name string
		pid  int
		want bool
	}{
		{"this process", os.Getpid(), true},
		{"exited and reaped child", exitedPid(t), false},
		{"launchd, owned by root", 1, true},
		{"zero", 0, false},
		{"negative", -5, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsProcessAlive(c.pid); got != c.want {
				t.Fatalf("IsProcessAlive(%d) = %v, want %v", c.pid, got, c.want)
			}
		})
	}
}

func TestIsProcessAliveWithoutPsOnPath(t *testing.T) {
	t.Setenv("PATH", "")

	if !IsProcessAlive(os.Getpid()) {
		t.Fatal("this process reads as exited when ps cannot be found")
	}
}

func TestMonitorParentProcessLogsNoErrorForAnExitedParent(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "updater.log")
	if err := logger.GetInstance().Initialize(logPath); err != nil {
		t.Fatal(err)
	}

	if err := MonitorParentProcess(exitedPid(t)); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "ERROR") {
		t.Fatalf("an ordinary parent exit was logged as an error:\n%s", data)
	}
}
