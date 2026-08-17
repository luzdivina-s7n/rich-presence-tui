package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"rich-presence-tui/assets"
	"rich-presence-tui/config"
	"rich-presence-tui/discord"
	"rich-presence-tui/tray"
	"strings"
	"time"
)

type appCtx struct {
	cfg      *config.Config
	presence *discord.Presence
}
type view interface {
	Update(msg tea.Msg) tea.Cmd
	View(width, height int) string
}
type langToggleable interface {
	langOK() bool
}
type hinter interface {
	hint() string
}
type Model struct {
	app           appCtx
	screen        screen
	payload       interface{}
	child         view
	width, height int
	status        statusModel
	quitting      bool
}

func New(cfg *config.Config, presence *discord.Presence) Model {
	assets.SetLanguage(cfg.Settings.Language)
	m := Model{
		app: appCtx{cfg: cfg, presence: presence},
	}
	m.screen = screenMainMenu
	m.child = m.makeChild(screenMainMenu, nil)
	return m
}
func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, m.forward(msg)
	case tea.KeyMsg:
		switch normKey(msg) {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+z":
			if m.app.cfg.Settings.MinimizeToTray {
				tray.Minimize()
				m.status.set(assets.T("status.trayed"), statusInfo)
				return m, nil
			}
		case "l":
			if t, ok := m.child.(langToggleable); ok && t.langOK() {
				return m, m.toggleLangCmd()
			}
		}
		return m, m.forward(msg)
	case navMsg:
		m.screen = msg.target
		m.payload = msg.payload
		m.child = m.makeChild(msg.target, msg.payload)
		return m, nil
	case backMsg:
		m.screen = msg.to
		m.payload = msg.payload
		m.child = m.makeChild(msg.to, msg.payload)
		return m, nil
	case statusMsg:
		m.status.set(msg.text, msg.kind)
		return m, expireStatusCmd()
	case statusClearMsg:
		m.status.text = ""
		return m, nil
	case errorMsg:
		m.status.set(msg.err.Error(), statusError)
		return m, expireStatusCmd()
	case savedMsg:
		m.status.set(assets.T("status.saved"), statusSuccess)
		return m, expireStatusCmd()
	case startedPresetMsg:
		if p := m.app.cfg.Preset(msg.id); p != nil {
			p.Enabled = true
		}
		_ = m.app.cfg.Save()
		m.refreshChild()
		m.status.set(assets.T("status.started")+" "+msg.name, statusSuccess)
		return m, expireStatusCmd()
	case stoppedPresetMsg:
		if p := m.app.cfg.Preset(msg.id); p != nil {
			p.Enabled = false
		}
		_ = m.app.cfg.Save()
		m.refreshChild()
		m.status.set(assets.T("status.stopped"), statusInfo)
		return m, expireStatusCmd()
	case presetRefreshedMsg:
		if p := m.app.cfg.Preset(msg.id); p != nil {
			p.Enabled = true
		}
		return m, nil
	case langChangedMsg:
		if msg.code != m.app.cfg.Settings.Language {
			m.app.cfg.Settings.Language = msg.code
			assets.SetLanguage(msg.code)
			_ = m.app.cfg.Save()
			m.child = m.makeChild(m.screen, m.payload)
			m.status.set(fmt.Sprintf(assets.T("status.lang"), assets.LanguageName(msg.code)), statusSuccess)
		}
		return m, expireStatusCmd()
	case TrayOpenMsg:
		return m, nil
	default:
		return m, m.forward(msg)
	}
}
func (m Model) forward(msg tea.Msg) tea.Cmd {
	if m.child == nil {
		return nil
	}
	return m.child.Update(msg)
}

const (
	uiW     = 76
	uiH     = 26
	innerW  = uiW - 6
	innerH  = uiH - 4
	footerH = 3
	childH  = innerH - footerH
)

func (m Model) View() string {
	var b strings.Builder
	if m.child != nil {
		b.WriteString(strings.TrimRight(m.child.View(innerW, childH), "\n"))
		b.WriteString("\n")
	}
	b.WriteString(m.footer())
	content := appStyle.Render(padFrame(b.String(), innerW, innerH))
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}
func (m Model) footer() string {
	rule := dimStyle.Render(strings.Repeat("-", innerW))
	hint := ""
	if h, ok := m.child.(hinter); ok {
		hint = h.hint()
	}
	hintLine := hintTextStyle.Render(truncTo(hint, innerW))
	left := badgeInactive.Render("[ ] " + assets.T("status.idle"))
	if p := activePreset(m.app.cfg); p != nil {
		left = badgeActive.Render("[x] " + truncTo(p.Name, 24) + "  " + assets.T("status.running"))
	}
	line := left
	if msg := m.status.text; msg != "" {
		var st lipgloss.Style
		switch m.status.kind {
		case statusError:
			st = errorStyle
		case statusSuccess:
			st = successStyle
		default:
			st = infoStyle
		}
		line += "   " + st.Render(truncTo(msg, innerW-lipgloss.Width(left)-3))
	}
	return rule + "\n" + hintLine + "\n" + line
}
func (m Model) makeChild(s screen, payload interface{}) view {
	switch s {
	case screenAppID:
		return newAppIDModel(m.app, payload)
	case screenSettings:
		return newSettingsModel(m.app)
	case screenHelp:
		return newHelpModel(m.app)
	case screenPresetActions:
		return newPresetActionsModel(m.app, payload.(string))
	case screenEdit:
		return newEditModel(m.app, payload.(string))
	default:
		return newMainMenuModel(m.app)
	}
}
func (m *Model) refreshChild() {
	m.child = m.makeChild(m.screen, m.payload)
}
func (m Model) toggleLangCmd() tea.Cmd {
	newLang := "en"
	if m.app.cfg.Settings.Language == "en" {
		newLang = "ru"
	}
	return func() tea.Msg { return langChangedMsg{code: newLang} }
}
func expireStatusCmd() tea.Cmd {
	return tea.Tick(6*time.Second, func(time.Time) tea.Msg {
		return statusClearMsg{}
	})
}
