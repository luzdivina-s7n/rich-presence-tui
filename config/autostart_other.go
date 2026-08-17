//go:build !windows

package config
func AutostartEnabled() bool {
	return false
}
func SetAutostart(enabled bool) error {
	return errAutostartUnsupported
}
var errAutostartUnsupported = &autostartError{}
type autostartError struct{}
func (*autostartError) Error() string {
	return "autostart is not supported on this platform"
}
