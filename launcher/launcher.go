package launcher

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kihyun1998/macos_zip_updater/logger"
)

// LaunchApp 앱을 실행합니다.
func LaunchApp(appPath string) {
	log := logger.GetInstance()

	// 앱이름 추출
	appName := filepath.Base(appPath)
	appName = strings.TrimSuffix(appName, ".app")

	log.Log(fmt.Sprintf("Launching Flutter app: %s (app name: %s)", appPath, appName))

	// 앱 실행
	cmd := exec.Command("open", "-a", appName)
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return
	}

	log.Log(fmt.Sprintf("Flutter app \"%s\" launched successfully", appName))
}
