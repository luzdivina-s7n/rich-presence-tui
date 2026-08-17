package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"rich-presence-tui/assets"
	"strings"
)

type appIDTarget string

const (
	appIDTargetActive appIDTarget = "active"
	appIDTargetPreset appIDTarget = "preset"
)

type appIDRoute struct {
	target appIDTarget
	back   screen
}
type appIDModel struct {
	app     appCtx
	target  appIDTarget
	preset  string
	back    screen
	list    listModel
	items   []listItem
	editing bool
	input   textinput.Model
}

func newAppIDModel(app appCtx, payload interface{}) *appIDModel {
	m := &appIDModel{app: app}
	m.target = appIDTargetActive
	m.back = screenSettings
	switch v := payload.(type) {
	case appIDRoute:
		m.target = v.target
		m.back = v.back
	case string:
		m.target = appIDTargetPreset
		m.preset = v
		m.back = screenPresetActions
	}
	m.rebuild()
	return m
}
func (m *appIDModel) langOK() bool { return !m.editing }
func (m *appIDModel) rebuild() {
	m.items = m.items[:0]
	cur := m.currentID()
	for _, id := range m.app.cfg.AppIDs {
		badge := ""
		if id == cur {
			badge = "[x]"
		}
		m.items = append(m.items, listItem{title: id, badge: badge})
	}
	m.items = append(m.items,
		listItem{title: "+ " + assets.T("appid.enter")},
		listItem{title: "<- " + assets.T("common.back")},
	)
	m.list.SetItems(m.items)
}
func (m *appIDModel) currentID() string {
	if m.target == appIDTargetPreset {
		if p := m.app.cfg.Preset(m.preset); p != nil {
			return p.AppID
		}
		return ""
	}
	return m.app.cfg.Settings.ActiveAppID
}
func (m *appIDModel) apply(id string) tea.Cmd {
	if m.target == appIDTargetPreset {
		if p := m.app.cfg.Preset(m.preset); p != nil {
			p.AppID = id
		}
	} else {
		m.app.cfg.Settings.ActiveAppID = id
	}
	m.app.cfg.AddAppID(id)
	if err := m.app.cfg.Save(); err != nil {
		return errCmd(err)
	}
	var payload interface{}
	if m.target == appIDTargetPreset {
		payload = m.preset
	}
	return tea.Batch(backCmd(m.back, payload), statusCmd(assets.T("status.saved"), statusSuccess))
}
func (m *appIDModel) Update(msg tea.Msg) tea.Cmd {
	if m.editing {
		return m.updateEditing(msg)
	}
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch normKey(msg) {
		case "up", "k":
			m.list.MoveUp()
		case "down", "j":
			m.list.MoveDown()
		case "enter", " ", "right":
			idx := m.list.cursor
			switch {
			case idx < len(m.app.cfg.AppIDs):
				return m.apply(m.app.cfg.AppIDs[idx])
			case idx == len(m.app.cfg.AppIDs):
				var cmd tea.Cmd
				m.input, cmd = newInput(assets.T("appid.prompt"), "", 24)
				m.editing = true
				return cmd
			default:
				return m.backCmd()
			}
		case "d":
			idx := m.list.cursor
			if idx < len(m.app.cfg.AppIDs) {
				id := m.app.cfg.AppIDs[idx]
				m.app.cfg.RemoveAppID(id)
				if m.currentID() == id {
					if m.target == appIDTargetPreset {
						if p := m.app.cfg.Preset(m.preset); p != nil {
							p.AppID = ""
						}
					} else {
						m.app.cfg.Settings.ActiveAppID = ""
					}
				}
				if err := m.app.cfg.Save(); err != nil {
					return errCmd(err)
				}
				m.rebuild()
				return nil
			}
		case "esc", "left":
			return m.backCmd()
		}
	}
	return nil
}
func (m *appIDModel) updateEditing(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			id := strings.TrimSpace(m.input.Value())
			m.editing = false
			if id == "" {
				return nil
			}
			if !isNumeric(id) {
				return statusCmd(assets.T("appid.invalid"), statusError)
			}
			return m.apply(id)
		case "esc":
			m.editing = false
			return nil
		}
	}
	return cmd
}
func (m *appIDModel) backCmd() tea.Cmd {
	var payload interface{}
	if m.target == appIDTargetPreset {
		payload = m.preset
	}
	return backCmd(m.back, payload)
}
func (m *appIDModel) View(width, height int) string {
	if m.editing {
		var b strings.Builder
		b.WriteString(sectionStyle.Render(assets.T("appid.title")))
		b.WriteString("\n\n  ")
		b.WriteString(inputPromptStyle.Render("> "))
		b.WriteString(m.input.View())
		return b.String()
	}
	var b strings.Builder
	b.WriteString(sectionStyle.Render(assets.T("appid.title")))
	b.WriteString("\n\n")
	b.WriteString(renderList(assets.T("appid.title"), m.items, m.list.cursor, width))
	return b.String()
}
func (m *appIDModel) hint() string {
	if m.editing {
		return assets.T("edit.hint.edit")
	}
	return assets.T("appid.hint")
}
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
