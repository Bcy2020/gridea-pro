//go:build windows

package notify

import (
	"sync"

	toast "git.sr.ht/~jackmordaunt/go-toast/v2"
)

// sendPlatform 走 WinRT toast。
//
//  1. 首次调用前必须用 toast.SetAppData 把 AppID / GUID / IconPath 写进 HKCU 注册表，
//     COM 通知路径才能 attribute 到 Gridea Pro。否则 COM 必败，只有 PowerShell
//     fallback 兜底，且某些 Windows 版本上 PS fallback 也不会弹出。
//  2. SetAppData 是 idempotent 的（库内部判等），重复调用没副作用。
//  3. 完整的角标小图标显示仍要求 Start Menu 有 AUMID 关联的快捷方式
//     （NSIS 安装器的事），但本步已经让通知能弹出 + 来源名显示成 "Gridea Pro"。
var setAppDataOnce sync.Once

func ensureAppData() {
	setAppDataOnce.Do(func() {
		_ = toast.SetAppData(toast.AppData{
			AppID:    appDisplayName,
			GUID:     appGUID,
			IconPath: appIconPath(),
		})
	})
}

func sendPlatform(title, body string) error {
	ensureAppData()
	n := toast.Notification{
		AppID: appDisplayName,
		Title: title,
		Body:  body,
		Icon:  appIconPath(),
	}
	return n.Push()
}
