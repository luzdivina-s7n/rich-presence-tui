package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"rich-presence-tui/assets"
	"rich-presence-tui/config"
	"strings"
)

type settingsModel struct {
	app                appCtx
	list               listModel
	items              []listItem
	autostartSupported bool
	autostartState     bool
}

func newSettingsModel(app appCtx) *settingsModel {
	m := &settingsModel{app: app}
	m.autostartState = config.AutostartEnabled()
	m.autostartSupported = true
	m.rebuild()
	return m
}
func (m *settingsModel) langOK() bool { return true }
func (m *settingsModel) rebuild() {
	langName := assets.LanguageName(m.app.cfg.Settings.Language)
	as := assets.T("common.off")
	if m.autostartState {
		as = assets.T("common.on")
	}
	if !m.autostartSupported {
		as = assets.T("settings.unsupported")
	}
	mt := assets.T("common.off")
	if m.app.cfg.Settings.MinimizeToTray {
		mt = assets.T("common.on")
	}
	cd := assets.T("common.off")
	if m.app.cfg.Settings.ConfirmDelete {
		cd = assets.T("common.on")
	}
	aid := m.app.cfg.Settings.ActiveAppID
	if aid == "" {
		aid = "-"
	}
	m.items = []listItem{
		{title: assets.T("settings.language"), desc: langName},
		{title: assets.T("settings.autostart"), desc: as},
		{title: assets.T("settings.minimizeTray"), desc: mt},
		{title: assets.T("settings.confirmDelete"), desc: cd},
		{title: assets.T("settings.defaultAppID"), desc: aid},
		{title: assets.T("settings.configPath"), desc: m.app.cfg.Path()},
		{title: "<- " + assets.T("settings.back")},
	}
	m.list.SetItems(m.items)
}
func (m *settingsModel) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch normKey(msg) {
		case "up", "k":
			m.list.MoveUp()
		case "down", "j":
			m.list.MoveDown()
		case "enter", " ", "right":
			switch m.list.cursor {
			case 0:
				return m.toggleLang()
			case 1:
				if !m.autostartSupported {
					return nil
				}
				m.autostartState = !m.autostartState
				m.app.cfg.Settings.Autostart = m.autostartState
				if err := config.SetAutostart(m.autostartState); err != nil {
					m.autostartState = config.AutostartEnabled()
					return statusCmd(err.Error(), statusError)
				}
				m.rebuild()
				return statusCmd(assets.T("status.saved"), statusSuccess)
			case 2:
				m.app.cfg.Settings.MinimizeToTray = !m.app.cfg.Settings.MinimizeToTray
				if err := m.app.cfg.Save(); err != nil {
					return errCmd(err)
				}
				m.rebuild()
				return statusCmd(assets.T("status.saved"), statusSuccess)
			case 3:
				m.app.cfg.Settings.ConfirmDelete = !m.app.cfg.Settings.ConfirmDelete
				if err := m.app.cfg.Save(); err != nil {
					return errCmd(err)
				}
				m.rebuild()
				return statusCmd(assets.T("status.saved"), statusSuccess)
			case 4:
				return navCmd(screenAppID, appIDRoute{target: appIDTargetActive, back: screenSettings})
			case 5:
				return nil
			default:
				return backCmd(screenMainMenu, nil)
			}
		case "esc", "left":
			return backCmd(screenMainMenu, nil)
		}
	}
	return nil
}
func (m *settingsModel) toggleLang() tea.Cmd {
	code := "en"
	if m.app.cfg.Settings.Language == "en" {
		code = "ru"
	}
	m.app.cfg.Settings.Language = code
	assets.SetLanguage(code)
	if err := m.app.cfg.Save(); err != nil {
		return errCmd(err)
	}
	m.rebuild()
	return statusCmd(fmt.Sprintf(assets.T("status.lang"), assets.LanguageName(code)), statusSuccess)
}
func (m *settingsModel) View(width, height int) string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render(assets.T("settings.title")))
	b.WriteString("\n\n")
	b.WriteString(renderList("", m.items, m.list.cursor, width))
	return b.String()
}
func (m *settingsModel) hint() string {
	return assets.T("settings.hint")
}
