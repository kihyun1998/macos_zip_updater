package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"
)

// Logger 싱글톤 인스턴스
type Logger struct {
	logFile *os.File
	mu      sync.Mutex
}

var (
	instance *Logger
	once     sync.Once
)

// GetInstance Logger 인스턴스를 가져옵니다.
func GetInstance() *Logger {
	once.Do(func() {
		instance = &Logger{}
	})
	return instance
}

// Initialize 로그 파일을 초기화합니다.
func (l *Logger) Initialize(logFilePath string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 로그 디렉토리 생성
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 로그 파일 열기 (append mode)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.logFile = file
	return nil
}

// Log 로그 메시지를 기록합니다.
func (l *Logger) Log(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logFile == nil {
		fmt.Fprintln(os.Stderr, "Logger not initialized")
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	if _, err := l.logFile.WriteString(logMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
		return
	}

	// flush
	l.logFile.Sync()
}

// Error 에러 로그를 기록합니다.
func (l *Logger) Error(err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logFile == nil {
		fmt.Fprintln(os.Stderr, "Logger not initialized")
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	stackTrace := string(debug.Stack())
	logMessage := fmt.Sprintf("[%s]   ERROR=%s\nSTACKTRACE=%s\n", timestamp, err.Error(), stackTrace)

	if _, writeErr := l.logFile.WriteString(logMessage); writeErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", writeErr)
		return
	}

	// flush
	l.logFile.Sync()
}

// Close 로그 파일을 닫습니다.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
