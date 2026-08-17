package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"rich-presence-tui/assets"
	"rich-presence-tui/config"
	"rich-presence-tui/discord"
	"rich-presence-tui/tray"
	"strings"
	"time"
)

type presetActionsModel struct {
	app     appCtx
	preset  *config.Preset
	list    listModel
	items   []listItem
	confirm *confirmModel
}

func newPresetActionsModel(app appCtx, id string) *presetActionsModel {
	m := &presetActionsModel{app: app}
	m.preset = app.cfg.Preset(id)
	m.rebuild()
	return m
}
func (m *presetActionsModel) rebuild() {
	m.items = m.items[:0]
	label := assets.T("pact.start")
	if m.preset != nil && m.preset.Enabled {
		label = assets.T("pact.restart")
	}
	m.items = []listItem{
		{title: assets.T("pact.editActivity")},
		{title: assets.T("pact.changeAppID")},
		{title: label},
		{title: assets.T("pact.stop")},
		{title: assets.T("pact.tray")},
		{title: assets.T("pact.delete")},
		{title: "<- " + assets.T("pact.back")},
	}
	m.list.SetItems(m.items)
}
func (m *presetActionsModel) Update(msg tea.Msg) tea.Cmd {
	if m.confirm != nil {
		var cmd tea.Cmd
		*m.confirm, cmd = m.confirm.Update(msg)
		if m.confirm.confirmed() {
			m.confirm = nil
			return cmd
		}
		if m.confirm.dismissed() {
			m.confirm = nil
			return cmd
		}
		return cmd
	}
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch normKey(msg) {
		case "up", "k":
			m.list.MoveUp()
		case "down", "j":
			m.list.MoveDown()
		case "enter", " ", "right":
			return m.activate(m.list.cursor)
		case "esc", "left":
			return backCmd(screenMainMenu, nil)
		}
	}
	return nil
}
func (m *presetActionsModel) activate(idx int) tea.Cmd {
	if m.preset == nil {
		return backCmd(screenMainMenu, nil)
	}
	switch idx {
	case 0:
		return navCmd(screenEdit, m.preset.ID)
	case 1:
		return navCmd(screenAppID, m.preset.ID)
	case 2:
		return startPresetCmd(m.app, m.preset)
	case 3:
		return stopPresetCmd(m.app, m.preset)
	case 4:
		tray.Minimize()
		return statusCmd(assets.T("status.trayed"), statusInfo)
	case 5:
		if m.app.cfg.Settings.ConfirmDelete {
			m.confirm = &confirmModel{
				message: fmt.Sprintf(assets.T("pact.confirm"), m.preset.Name),
				yes:     true,
				onYes:   func() tea.Cmd { return m.deleteCmd() },
			}
			return nil
		}
		return m.deleteCmd()
	default:
		return backCmd(screenMainMenu, nil)
	}
}
func (m *presetActionsModel) deleteCmd() tea.Cmd {
	id := m.preset.ID
	m.app.cfg.DeletePreset(id)
	if err := m.app.cfg.Save(); err != nil {
		return errCmd(err)
	}
	return tea.Batch(
		statusCmd(assets.T("pact.deleted"), statusSuccess),
		backCmd(screenMainMenu, nil),
	)
}
func startPresetCmd(app appCtx, preset *config.Preset) tea.Cmd {
	return func() tea.Msg {
		act := discord.BuildActivity(preset.Activity, time.Now())
		if act.Name == "" {
			act.Name = preset.Name
		}
		if err := app.presence.Set(preset.AppID, act); err != nil {
			return errorMsg{err: fmt.Errorf("%s (%w)", assets.T("status.discordOff"), err)}
		}
		return startedPresetMsg{id: preset.ID, name: preset.Name}
	}
}
func refreshPresetCmd(app appCtx, preset *config.Preset) tea.Cmd {
	return func() tea.Msg {
		act := discord.BuildActivity(preset.Activity, time.Now())
		if act.Name == "" {
			act.Name = preset.Name
		}
		if err := app.presence.Set(preset.AppID, act); err != nil {
			return errorMsg{err: err}
		}
		return presetRefreshedMsg{id: preset.ID}
	}
}
func stopPresetCmd(app appCtx, preset *config.Preset) tea.Cmd {
	return func() tea.Msg {
		if err := app.presence.Clear(); err != nil {
			return errorMsg{err: err}
		}
		return stoppedPresetMsg{id: preset.ID}
	}
}
func (m *presetActionsModel) View(width, height int) string {
	if m.preset == nil {
		return mutedStyle.Render(assets.T("status.notFound"))
	}
	state := assets.T("pact.inactive")
	stateStyle := badgeInactive
	if m.preset.Enabled {
		state = assets.T("pact.active")
		stateStyle = badgeActive
	}
	var b strings.Builder
	b.WriteString(sectionStyle.Render(fmt.Sprintf(assets.T("pact.title"), m.preset.Name)))
	b.WriteString("\n  ")
	b.WriteString(fieldLabelStyle.Render("AppID: "))
	b.WriteString(fieldValueStyle.Render(m.preset.AppID))
	b.WriteString("   ")
	marker := "[ ]"
	if m.preset.Enabled {
		marker = "[x]"
	}
	b.WriteString(stateStyle.Render(marker + " " + state))
	b.WriteString("\n\n")
	rows := make([]listItem, len(m.items))
	copy(rows, m.items)
	b.WriteString(renderList("", rows, m.list.cursor, width))
	content := b.String()
	if m.confirm != nil {
		content = m.confirm.View(width)
	}
	return content
}
func (m *presetActionsModel) hint() string {
	if m.confirm != nil {
		return ""
	}
	return assets.T("pact.hint")
}
