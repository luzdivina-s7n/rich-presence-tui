package tray
type Callbacks struct {
	OnOpen func()
	OnQuit func()
}
var active Callbacks
func Start(cb Callbacks) {
	active = cb
	startPlatform()
}
func Minimize() {
	minimizePlatform()
}
func Restore() {
	restorePlatform()
}
