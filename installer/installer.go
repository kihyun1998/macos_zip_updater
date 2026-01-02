package installer

import (
	"fmt"
	"os"

	"github.com/kihyun1998/macos_zip_updater/logger"
)

// InstallNewApp 새 앱을 설치합니다.
func InstallNewApp(updateAppPath, installPath string) bool {
	log := logger.GetInstance()

	log.Log("Installing new app...")
	log.Log(fmt.Sprintf("From: %s", updateAppPath))
	log.Log(fmt.Sprintf("To: %s", installPath))

	// update app 있나 검사
	if _, err := os.Stat(updateAppPath); os.IsNotExist(err) {
		log.Log(fmt.Sprintf("Update app does not exist: %s", updateAppPath))
		return false
	}

	// Rename으로 이동
	if err := os.Rename(updateAppPath, installPath); err != nil {
		log.Error(err)
		return false
	}

	// 설치 검증
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		log.Log(fmt.Sprintf("Installation verification failed: app not found at %s", installPath))
		return false
	}

	log.Log("New app installed successfully")
	return true
}
