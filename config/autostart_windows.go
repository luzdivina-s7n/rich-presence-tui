//go:build windows

package config
import (
	"os"
	"path/filepath"
	"strings"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)
const autostartKey = `Software\Microsoft\Windows\CurrentVersion\Run`
func AutostartEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, autostartKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	val, _, err := k.GetStringValue("rich-presence-tui")
	if err != nil {
		return false
	}
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(val), strings.ToLower(exe))
}
func SetAutostart(enabled bool) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.Abs(exe)
	if err != nil {
		return err
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, autostartKey, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		if err == windows.ERROR_FILE_NOT_FOUND {
			k, _, err = registry.CreateKey(registry.CURRENT_USER, autostartKey, registry.SET_VALUE)
		}
		if err != nil {
			return err
		}
	}
	defer k.Close()
	if enabled {
		return k.SetStringValue("rich-presence-tui", `"`+exe+`" --tray`)
	}
	return k.DeleteValue("rich-presence-tui")
}
