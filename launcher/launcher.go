package launcher

import (
	"fmt"
	"os/exec"

	"github.com/kihyun1998/macos_zip_updater/logger"
)

// runOpen runs `open` with args and returns its combined output.
var runOpen = func(args ...string) ([]byte, error) {
	return exec.Command("open", args...).CombinedOutput()
}

// LaunchApp 앱을 실행합니다.
func LaunchApp(appPath string) {
	log := logger.GetInstance()

	log.Log(fmt.Sprintf("Launching Flutter app: %s", appPath))

	// 앱 실행
	if output, err := runOpen(appPath); err != nil {
		log.Error(fmt.Errorf("open %s: %w: %s", appPath, err, output))
		return
	}

	log.Log(fmt.Sprintf("Flutter app \"%s\" launched successfully", appPath))
}
