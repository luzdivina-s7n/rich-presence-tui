package tui
type screen int
const (
	screenMainMenu screen = iota
	screenAppID
	screenSettings
	screenPresetActions
	screenEdit
	screenHelp
)
type navMsg struct {
	target  screen
	payload interface{}
}
type backMsg struct {
	to      screen
	payload interface{}
}
type statusKind int
const (
	statusInfo statusKind = iota
	statusSuccess
	statusError
)
type statusMsg struct {
	text string
	kind statusKind
}
type statusClearMsg struct{}
type errorMsg struct {
	err error
}
type savedMsg struct{}
type startedPresetMsg struct {
	id   string
	name string
}
type stoppedPresetMsg struct {
	id string
}
type presetRefreshedMsg struct {
	id string
}
type TrayOpenMsg struct{}
type fieldChangedMsg struct {
	id    string
	value string
}
type langChangedMsg struct {
	code string
}
type presetNameMsg struct {
	name string
}
