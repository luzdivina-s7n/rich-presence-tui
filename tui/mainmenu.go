package tui

import (
	"errors"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"rich-presence-tui/assets"
	"rich-presence-tui/config"
	"rich-presence-tui/tray"
	"strings"
)

type mainMenuModel struct {
	app     appCtx
	list    listModel
	items   []listItem
	editing bool
	input   textinput.Model
	confirm *confirmModel
}

func newMainMenuModel(app appCtx) *mainMenuModel {
	m := &mainMenuModel{app: app}
	m.rebuild()
	return m
}
func (m *mainMenuModel) langOK() bool { return !m.editing && m.confirm == nil }
func (m *mainMenuModel) rebuild() {
	items := m.items[:0]
	items = append(items, listItem{title: assets.T("menu.presets"), header: true})
	for _, p := range m.app.cfg.Presets {
		badge := "[ ]"
		if p.Enabled {
			badge = "[x]"
		}
		items = append(items, listItem{
			title: p.Name,
			desc:  "AppID: " + p.AppID,
			badge: badge,
		})
	}
	items = append(items, listItem{title: "+ " + assets.T("presets.create")})
	items = append(items, listItem{title: assets.T("menu.section"), header: true})
	items = append(items, listItem{title: assets.T("menu.reset")})
	items = append(items, listItem{title: assets.T("menu.settings")})
	items = append(items, listItem{title: assets.T("menu.exit")})
	m.items = items
	m.list.SetItems(items)
}
func (m *mainMenuModel) presetIndex() int {
	if m.list.cursor >= 1 && m.list.cursor <= len(m.app.cfg.Presets) {
		return m.list.cursor - 1
	}
	return -1
}
func (m *mainMenuModel) selectedPreset() *config.Preset {
	i := m.presetIndex()
	if i < 0 || i >= len(m.app.cfg.Presets) {
		return nil
	}
	return m.app.cfg.Presets[i]
}
func (m *mainMenuModel) Update(msg tea.Msg) tea.Cmd {
	if m.confirm != nil {
		var cmd tea.Cmd
		*m.confirm, cmd = m.confirm.Update(msg)
		if m.confirm.confirmed() || m.confirm.dismissed() {
			m.confirm = nil
		}
		return cmd
	}
	if m.editing {
		return m.updateEditing(msg)
	}
	if key, ok := msg.(tea.KeyMsg); ok {
		switch normKey(key) {
		case "up", "k":
			m.list.MoveUp()
		case "down", "j":
			m.list.MoveDown()
		case "enter", " ":
			return m.activate(m.list.cursor)
		case "right":
			return m.open(m.list.cursor)
		case "n":
			return m.beginCreate()
		case "e":
			if p := m.selectedPreset(); p != nil {
				return navCmd(screenEdit, p.ID)
			}
			return statusCmd(assets.T("menu.nopreset"), statusError)
		case "a":
			if p := m.selectedPreset(); p != nil {
				return navCmd(screenPresetActions, p.ID)
			}
			return statusCmd(assets.T("menu.nopreset"), statusError)
		case "t":
			tray.Minimize()
			return statusCmd(assets.T("status.trayed"), statusInfo)
		case "q", "esc":
			return tea.Quit
		case "?":
			return navCmd(screenHelp, nil)
		}
	}
	return nil
}
func (m *mainMenuModel) activate(idx int) tea.Cmd {
	n := len(m.app.cfg.Presets)
	switch {
	case idx >= 1 && idx <= n:
		return m.togglePreset(idx - 1)
	case idx == n+1:
		return m.beginCreate()
	case idx == n+3:
		return m.resetCmd()
	case idx == n+4:
		return navCmd(screenSettings, nil)
	case idx == n+5:
		return tea.Quit
	}
	return nil
}
func (m *mainMenuModel) open(idx int) tea.Cmd {
	n := len(m.app.cfg.Presets)
	switch {
	case idx >= 1 && idx <= n:
		return navCmd(screenPresetActions, m.app.cfg.Presets[idx-1].ID)
	case idx == n+1:
		return m.beginCreate()
	case idx == n+3:
		return m.resetCmd()
	case idx == n+4:
		return navCmd(screenSettings, nil)
	case idx == n+5:
		return tea.Quit
	}
	return nil
}
func (m *mainMenuModel) resetCmd() tea.Cmd {
	m.confirm = &confirmModel{
		message: assets.T("menu.resetConfirm"),
		yes:     true,
		onYes: func() tea.Cmd {
			m.app.cfg.Settings = config.Settings{Language: "ru"}
			m.app.cfg.Presets = nil
			m.app.cfg.AppIDs = nil
			if err := m.app.cfg.Delete(); err != nil {
				return errCmd(err)
			}
			assets.SetLanguage(m.app.cfg.Settings.Language)
			m.rebuild()
			return statusCmd(assets.T("menu.resetDone"), statusSuccess)
		},
	}
	return nil
}
func (m *mainMenuModel) togglePreset(i int) tea.Cmd {
	if i < 0 || i >= len(m.app.cfg.Presets) {
		return nil
	}
	p := m.app.cfg.Presets[i]
	if p.Enabled {
		return stopPresetCmd(m.app, p)
	}
	return startPresetCmd(m.app, p)
}
func (m *mainMenuModel) beginCreate() tea.Cmd {
	if m.app.cfg.Full() {
		return statusCmd(assets.T("presets.full"), statusError)
	}
	var cmd tea.Cmd
	m.input, cmd = newInput(assets.T("presets.name"), "", 24)
	m.editing = true
	return cmd
}
func (m *mainMenuModel) updateEditing(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			name := strings.TrimSpace(m.input.Value())
			m.editing = false
			if name == "" {
				return nil
			}
			if m.app.cfg.PresetByName(name) != nil {
				return statusCmd(name+": "+assets.T("presets.exists"), statusError)
			}
			p, err := m.app.cfg.CreatePreset(name)
			if err != nil {
				if errors.Is(err, config.ErrPresetLimit) {
					return statusCmd(assets.T("presets.full"), statusError)
				}
				return errCmd(err)
			}
			if err := m.app.cfg.Save(); err != nil {
				return errCmd(err)
			}
			m.rebuild()
			return tea.Batch(
				statusCmd(assets.T("presets.created"), statusSuccess),
				navCmd(screenPresetActions, p.ID),
			)
		case "esc":
			m.editing = false
			return nil
		}
	}
	return cmd
}
func (m *mainMenuModel) View(width, height int) string {
	if m.editing {
		var b strings.Builder
		b.WriteString(sectionStyle.Render(assets.T("presets.create")))
		b.WriteString("\n\n  ")
		b.WriteString(inputPromptStyle.Render("> "))
		b.WriteString(m.input.View())
		return b.String()
	}
	var b strings.Builder
	b.WriteString(m.header())
	b.WriteString("\n\n")
	content := renderList("", m.items, m.list.cursor, width)
	if m.confirm != nil {
		content = m.confirm.View(width)
	}
	b.WriteString(content)
	return b.String()
}
func (m *mainMenuModel) hint() string {
	if m.confirm != nil {
		return ""
	}
	if m.editing {
		return assets.T("edit.hint.edit")
	}
	return assets.T("menu.hint")
}
func (m *mainMenuModel) header() string {
	h := titleStyle.Render(assets.T("app.title"))
	h += "   " + mutedStyle.Render("["+assets.LanguageName(m.app.cfg.Settings.Language)+"]")
	if p := activePreset(m.app.cfg); p != nil {
		h += "   " + badgeActive.Render("[x] "+p.Name)
	}
	return h
}
func activePreset(cfg *config.Config) *config.Preset {
	for _, p := range cfg.Presets {
		if p.Enabled {
			return p
		}
	}
	return nil
}
