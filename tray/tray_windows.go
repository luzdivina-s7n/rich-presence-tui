//go:build windows

package tray
import (
	"sync"
	"github.com/getlantern/systray"
	"rich-presence-tui/assets"
)
var (
	startOnce sync.Once
)
func startPlatform() {
	startOnce.Do(func() {
		go func() {
			systray.Run(onReady, onExit)
		}()
	})
}
func onReady() {
	defer func() {
		if r := recover(); r != nil {
			_ = r
		}
	}()
	systray.SetIcon(assets.TrayIcon())
	systray.SetTitle("Rich Presence TUI")
	systray.SetTooltip("Rich Presence TUI")
	mOpen := systray.AddMenuItem(assets.T("tray.open"), "")
	mQuit := systray.AddMenuItem(assets.T("tray.quit"), "")
	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				restorePlatform()
				if active.OnOpen != nil {
					active.OnOpen()
				}
			case <-mQuit.ClickedCh:
				if active.OnQuit != nil {
					active.OnQuit()
				}
			}
		}
	}()
}
func onExit() {}
func minimizePlatform() {
	showWindow(swHide)
}
func restorePlatform() {
	showWindow(swShow)
}
