package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kihyun1998/macos_zip_updater/backup"
	"github.com/kihyun1998/macos_zip_updater/installer"
	"github.com/kihyun1998/macos_zip_updater/launcher"
	"github.com/kihyun1998/macos_zip_updater/logger"
	"github.com/kihyun1998/macos_zip_updater/monitor"
)

func main() {
	// 첫 번째 인자로 모드 확인
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "--restart":
		handleRestartMode()
	case "--update":
		handleUpdateMode()
	default:
		fmt.Fprintf(os.Stderr, "Unknown mode: %s\n", mode)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  Restart mode: ACRA_Point_Client_Updater --restart <parent_pid> <log_file_path> <app_path>")
	fmt.Fprintln(os.Stderr, "  Update mode:  ACRA_Point_Client_Updater --update <parent_pid> <log_file_path> <app_path> <extracted_folder_path> <backup_path>")
}

func handleRestartMode() {
	// 인자: [--restart, 부모 PID, 로그 파일 경로, 앱 경로]
	if len(os.Args) < 5 {
		fmt.Fprintln(os.Stderr, "Restart mode requires: <parent_pid> <log_file_path> <app_path>")
		os.Exit(1)
	}

	parentPid, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parent PID: %v\n", err)
		os.Exit(1)
	}

	logFilePath := os.Args[3]
	appPath := os.Args[4]

	// Logger 초기화
	log := logger.GetInstance()
	if err := log.Initialize(logFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	currentPid := os.Getpid()

	// 시작 로그 기록
	log.Log(fmt.Sprintf("Restarter started (PID: %d)", currentPid))
	log.Log(fmt.Sprintf("Parent PID: %d", parentPid))
	log.Log(fmt.Sprintf("App path: %s", appPath))

	// 재시작 실행
	runRestart(parentPid, appPath)
}

func handleUpdateMode() {
	// 인자: [--update, 부모 PID, 로그 파일 경로, Flutter 앱 경로, 압축 해제된 폴더 경로, 백업 경로]
	if len(os.Args) < 7 {
		fmt.Fprintln(os.Stderr, "Update mode requires: <parent_pid> <log_file_path> <app_path> <extracted_folder_path> <backup_path>")
		os.Exit(1)
	}

	parentPid, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid parent PID: %v\n", err)
		os.Exit(1)
	}

	logFilePath := os.Args[3]
	appPath := os.Args[4]
	extractedFolderPath := os.Args[5]
	backupPath := os.Args[6]

	// Logger 초기화
	log := logger.GetInstance()
	if err := log.Initialize(logFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	currentPid := os.Getpid()

	// 시작 로그 기록
	log.Log(fmt.Sprintf("Updater started (PID: %d)", currentPid))
	log.Log(fmt.Sprintf("Parent PID: %d", parentPid))
	log.Log(fmt.Sprintf("App path: %s", appPath))
	log.Log(fmt.Sprintf("Extracted folder path: %s", extractedFolderPath))
	log.Log(fmt.Sprintf("Backup path: %s", backupPath))

	// 부모 프로세스 모니터링
	if err := monitor.MonitorParentProcess(parentPid); err != nil {
		log.Error(err)
	}

	// 업데이트 실행
	runUpdate(appPath, extractedFolderPath, backupPath)
}

func runRestart(parentPid int, appPath string) {
	log := logger.GetInstance()

	// 부모 프로세스 모니터링
	if err := monitor.MonitorParentProcess(parentPid); err != nil {
		log.Error(err)
	}

	// 앱 재실행
	launcher.LaunchApp(appPath)
	log.Log("Restart completed, exiting...")
	os.Exit(0)
}

func runUpdate(appPath, extractedPath, backupPath string) {
	log := logger.GetInstance()

	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Errorf("panic: %v", r))
			backup.CleanupUpdateDir(extractedPath)
			backup.RollbackFromBackup(appPath, backupPath)
			launcher.LaunchApp(appPath)
			os.Exit(1)
		}
	}()

	appName := filepath.Base(appPath)
	appName = strings.TrimSuffix(appName, ".app")

	// 1. 백업 (rename)
	if !backup.MoveAppToBackup(appPath, backupPath) {
		launcher.LaunchApp(appPath)
		log.Log("Daemon exiting due to backup failure...")
		os.Exit(1)
	}

	// 2. 업데이트 앱 존재 확인
	updateAppPath := filepath.Join(extractedPath, appName+".app")
	if !fileExists(updateAppPath) {
		log.Log(fmt.Sprintf("Extracted app not found: %s", updateAppPath))
		backup.CleanupUpdateDir(extractedPath)
		backup.RollbackFromBackup(appPath, backupPath)
		launcher.LaunchApp(appPath)
		os.Exit(1)
	}
	log.Log(fmt.Sprintf("Extracted app verified: %s", updateAppPath))

	// 3. 새 앱 설치 (rename)
	if !installer.InstallNewApp(updateAppPath, appPath) {
		log.Log("Failed to install new app, rolling back...")
		backup.CleanupUpdateDir(extractedPath)
		backup.RollbackFromBackup(appPath, backupPath)
		launcher.LaunchApp(appPath)
		os.Exit(1)
	}

	// 4. 성공 - 정리
	log.Log("Update successful!")
	backup.CleanupUpdateDir(extractedPath)
	backup.DeleteBackup(backupPath)

	// 5. 앱 재실행
	launcher.LaunchApp(appPath)
	log.Log("Update exiting...")
	os.Exit(0)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
