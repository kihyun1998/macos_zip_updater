package backup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/macos_zip_updater/logger"
)

// MoveAppToBackup 백업 함수
func MoveAppToBackup(appPath, backupPath string) bool {
	log := logger.GetInstance()

	log.Log("Starting backup (rename) of current app...")

	// app path 있는지 검사
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		log.Log(fmt.Sprintf("App does not exist at: %s", appPath))
		return false
	}

	// 기존 백업 삭제
	if _, err := os.Stat(backupPath); err == nil {
		log.Log("Removing existing backup...")
		if err := os.RemoveAll(backupPath); err != nil {
			log.Error(err)
			return false
		}
	}

	// .app 을 .app.backup으로 바꾸는 것이다.
	// backupPath는 .app.backup의 fullpath가 오는 것이다.
	// 그래서 부모 경로 만드는 것이다.
	backupParent := filepath.Dir(backupPath)
	if _, err := os.Stat(backupParent); os.IsNotExist(err) {
		if err := os.MkdirAll(backupParent, 0755); err != nil {
			log.Error(err)
			return false
		}
	}

	log.Log("Moving app to backup location...")
	if err := os.Rename(appPath, backupPath); err != nil {
		log.Error(err)
		return false
	}

	// 백업 검증
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		log.Log("Backup verification failed: backup directory does not exist")
		return false
	}

	log.Log(fmt.Sprintf("Backup completed successfully (renamed): %s", backupPath))
	return true
}

// RollbackFromBackup 롤백 함수
func RollbackFromBackup(appPath, backupPath string) bool {
	log := logger.GetInstance()

	log.Log("Rollback from backup...")

	// 백업 폴더 검사
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		log.Log(fmt.Sprintf("Backup does not exist: %s", backupPath))
		return false
	}

	// 실패한 앱이 있다면 삭제
	if _, err := os.Stat(appPath); err == nil {
		log.Log("Removing failed app...")
		if err := os.RemoveAll(appPath); err != nil {
			log.Error(err)
			return false
		}
	}

	// 백업을 원래 위치로 rename
	if err := os.Rename(backupPath, appPath); err != nil {
		log.Error(err)
		return false
	}

	// 롤백 검증
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		log.Log(fmt.Sprintf("Rollback verification failed: app not found at %s", appPath))
		return false
	}

	log.Log("Rollback completed successfully")
	return true
}

// DeleteBackup 백업을 삭제합니다 (성공 시 정리용).
func DeleteBackup(backupPath string) {
	log := logger.GetInstance()

	log.Log("Deleting backup...")

	if _, err := os.Stat(backupPath); err == nil {
		if err := os.RemoveAll(backupPath); err != nil {
			log.Error(err)
			return
		}
		log.Log("Backup deleted successfully")
	} else {
		log.Log("No backup to delete")
	}
}

// CleanupUpdateDir 임시 디렉토리 제거
func CleanupUpdateDir(extractedPath string) {
	log := logger.GetInstance()

	log.Log(fmt.Sprintf("Cleaning up extracted directory: %s", extractedPath))

	if _, err := os.Stat(extractedPath); err == nil {
		if err := os.RemoveAll(extractedPath); err != nil {
			log.Error(err)
			return
		}
		log.Log("Extracted directory deleted")
	}
}
