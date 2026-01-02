package monitor

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/acra/flutter_macos_updater/logger"
)

// IsProcessAlive 프로세스가 살아있는지 검사하는 함수
func IsProcessAlive(pid int) bool {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid))
	err := cmd.Run()
	if err != nil {
		log := logger.GetInstance()
		log.Error(err)
		return false
	}
	return true
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
