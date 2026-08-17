//go:build windows

package discord
import (
	"fmt"
	"net"
	"github.com/Microsoft/go-winio"
)
func dialPipe() (net.Conn, error) {
	var lastErr error
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf(`\\.\pipe\discord-ipc-%d`, i)
		conn, err := winio.DialPipe(path, nil)
		if err == nil {
			return conn, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("cannot connect to discord IPC pipe: %w", lastErr)
}
