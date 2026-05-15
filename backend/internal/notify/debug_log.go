package notify

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 简单的诊断日志 —— Windows GUI 应用没有 stdout，靠它定位为什么通知没弹出。
// 文件位置：<UserCacheDir>/gridea-pro/notify.log，单文件追加写，不轮转。
// 只在 Send 失败时写一行，正常情况下不产生任何文件。

var (
	debugLogPathOnce sync.Once
	debugLogPath     string
	debugLogMu       sync.Mutex
)

func debugLogFile() string {
	debugLogPathOnce.Do(func() {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return
		}
		dir := filepath.Join(cacheDir, "gridea-pro")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return
		}
		debugLogPath = filepath.Join(dir, "notify.log")
	})
	return debugLogPath
}

func writeDebugLog(title, body string, err error) {
	if err == nil {
		return // 正常成功不留痕，只记录失败
	}
	p := debugLogFile()
	if p == "" {
		return
	}
	debugLogMu.Lock()
	defer debugLogMu.Unlock()
	f, openErr := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if openErr != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "[%s] err=%v | title=%q body=%q\n", time.Now().Format(time.RFC3339), err, title, body)
}
