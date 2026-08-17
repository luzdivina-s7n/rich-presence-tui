package tui
import (
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	"rich-presence-tui/assets"
)
type helpModel struct {
	app appCtx
}
func newHelpModel(app appCtx) *helpModel {
	return &helpModel{app: app}
}
func (m *helpModel) langOK() bool { return true }
func (m *helpModel) hint() string {
	return assets.T("help.hint")
}
func (m *helpModel) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch normKey(key) {
		case "esc", "left", "q", "enter", " ":
			return backCmd(screenMainMenu, nil)
		}
	}
	return nil
}
func (m *helpModel) View(width, height int) string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render(assets.T("help.title")))
	b.WriteString("\n\n")
	b.WriteString(assets.T("help.keys"))
	b.WriteString("\n\n")
	b.WriteString(sectionStyle.Render(assets.T("help.section.visibility")))
	b.WriteString("\n\n")
	b.WriteString(assets.T("help.visibility"))
	return b.String()
}