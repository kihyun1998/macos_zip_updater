package monitor

import (
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/kihyun1998/macos_zip_updater/logger"
)

// IsProcessAlive 프로세스가 살아있는지 검사하는 함수
// kill(pid, 0): no error or EPERM means the process exists, ESRCH means it
// does not. A pid of 0 or below reads as exited.
func IsProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	if err == nil || errors.Is(err, syscall.EPERM) {
		return true
	}
	if !errors.Is(err, syscall.ESRCH) {
		logger.GetInstance().Error(fmt.Errorf("kill(%d, 0): %w", pid, err))
	}
	return false
}

// MonitorParentProcess 부모 프로세스를 모니터링합니다.
func MonitorParentProcess(parentPid int) error {
	log := logger.GetInstance()

	log.Log("*** Starting parent process monitoring... ***")
	checkCount := 0

	for IsProcessAlive(parentPid) {
		checkCount++
		log.Log(fmt.Sprintf("Parent process check #%d: still alive", checkCount))
		time.Sleep(1 * time.Second)
	}

	log.Log("*** Parent process exit detected! ***")
	return nil
}
