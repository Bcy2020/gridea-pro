// Package notify 跨平台 OS 通知中心发送 —— 三端各自走原生 API。
//
//   - macOS：CGO + NSUserNotification（来源 / 图标自动从 Info.plist 取）。
//     deprecated 但截至 macOS 26 仍可用，等做正式签名分发后再迁 UNUserNotificationCenter。
//   - Windows：jackmordaunt/go-toast/v2，AUMID = appUserModelID。完整图标显示
//     需要在安装器（NSIS）里把 AUMID 关联到 .exe + icon；未关联时通知能发出但角标可能缺。
//   - Linux：esiqveland/notify 直接走 D-Bus org.freedesktop.Notifications，
//     显式传 app_name + app_icon。覆盖所有主流桌面环境（GNOME / KDE / XFCE / Cinnamon / dunst 等）。
//
// 通知失败不应阻塞主流程，调用方可忽略 error。
package notify

// appUserModelID Windows 上的 AppUserModelID（AUMID）。安装器需要把这个值
// 写进 Start Menu 快捷方式的 PropertyStore，否则角标取不到应用图标。
const appUserModelID = "com.gridea.pro"

// appDisplayName 通知显示的应用名（Linux D-Bus app_name）。
// macOS 来源名来自 Info.plist CFBundleName，Windows 来自 AUMID 注册。
const appDisplayName = "Gridea Pro"

// appGUID 给 Windows 注册表用的稳定标识，不要再改（一旦改用户机器上重新注册）。
const appGUID = "8b0a3c4f-1d2e-4f5a-9b6c-7d8e9f0a1b2c"

// Send 向系统通知中心发送一条通知。title 是粗体标题，body 是正文。
// 不管成败都追加一行日志到用户缓存目录下的 notify.log，方便事后诊断。
func Send(title, body string) error {
	err := sendPlatform(title, body)
	writeDebugLog(title, body, err)
	return err
}
