package launcher

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kihyun1998/macos_zip_updater/logger"
)

// fakeOpen replaces runOpen for the test, recording each call's args and
// answering with output and err.
func fakeOpen(t *testing.T, output string, err error) *[][]string {
	t.Helper()
	var calls [][]string
	orig := runOpen
	runOpen = func(args ...string) ([]byte, error) {
		calls = append(calls, args)
		return []byte(output), err
	}
	t.Cleanup(func() { runOpen = orig })
	return &calls
}

func TestLaunchAppOpensTheBundleAtAppPath(t *testing.T) {
	calls := fakeOpen(t, "", nil)
	appPath := "/Users/someone/Downloads/ACRAPointClient.app"

	LaunchApp(appPath)

	want := [][]string{{appPath}}
	if !reflect.DeepEqual(*calls, want) {
		t.Fatalf("open called with %q, want %q", *calls, want)
	}
}

func TestLaunchAppLogsAFailedOpenWithItsOutput(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "updater.log")
	if err := logger.GetInstance().Initialize(logPath); err != nil {
		t.Fatal(err)
	}
	openOutput := "The file /nowhere/ACRAPointClient.app does not exist."
	fakeOpen(t, openOutput, errors.New("exit status 1"))

	LaunchApp("/nowhere/ACRAPointClient.app")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	logged := string(data)
	if !strings.Contains(logged, openOutput) {
		t.Errorf("log does not carry open's output %q:\n%s", openOutput, logged)
	}
	if strings.Contains(logged, "launched successfully") {
		t.Errorf("log reports success for a failed open:\n%s", logged)
	}
}
