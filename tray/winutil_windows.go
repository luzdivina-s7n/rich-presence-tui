//go:build windows

package tray
import "golang.org/x/sys/windows"
const (
	swHide = 0
	swShow = 5
)
func showWindow(cmd int) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	user32 := windows.NewLazySystemDLL("user32.dll")
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	showWindow := user32.NewProc("ShowWindow")
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return
	}
	showWindow.Call(hwnd, uintptr(cmd))
}
