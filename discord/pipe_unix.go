//go:build !windows

package discord
import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
)
func dialPipe() (net.Conn, error) {
	var base string
	if runtime.GOOS == "darwin" {
		base = os.Getenv("TMPDIR")
		if base == "" {
			base = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "discord")
		}
	} else {
		base = os.Getenv("XDG_RUNTIME_DIR")
		if base == "" {
			base = os.TempDir()
		}
	}
	var lastErr error
	for i := 0; i < 10; i++ {
		path := filepath.Join(base, fmt.Sprintf("discord-ipc-%d", i))
		conn, err := net.Dial("unix", path)
		if err == nil {
			return conn, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("cannot connect to discord ipc socket: %w", lastErr)
}
