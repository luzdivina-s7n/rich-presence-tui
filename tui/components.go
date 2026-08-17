package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"rich-presence-tui/assets"
	"strings"
	"time"
)

type listItem struct {
	title  string
	desc   string
	badge  string
	dim    bool
	header bool
}
type listModel struct {
	items  []listItem
	cursor int
}

func (l *listModel) SetItems(items []listItem) {
	l.items = items
	if len(items) == 0 {
		l.cursor = 0
		return
	}
	if l.cursor >= len(items) {
		l.cursor = len(items) - 1
	}
	if items[l.cursor].header {
		l.MoveDown()
	}
}
func (l *listModel) MoveUp() {
	for i := l.cursor - 1; i >= 0; i-- {
		if !l.items[i].header {
			l.cursor = i
			return
		}
	}
}
func (l *listModel) MoveDown() {
	for i := l.cursor + 1; i < len(l.items); i++ {
		if !l.items[i].header {
			l.cursor = i
			return
		}
	}
}
func (l *listModel) reset() {
	l.cursor = 0
	if len(l.items) > 0 && l.items[l.cursor].header {
		l.MoveDown()
	}
}
func renderList(title string, items []listItem, cursor int, width int) string {
	var b strings.Builder
	if title != "" {
		b.WriteString(sectionStyle.Render(title))
		b.WriteString("\n")
	}
	if len(items) == 0 {
		b.WriteString(dimRow.Render("  " + assets.T("presets.empty")))
		return b.String()
	}
	rowW := width - 4
	if rowW < 20 {
		rowW = 20
	}
	for i, it := range items {
		if it.header {
			b.WriteString("\n")
			b.WriteString(dimStyle.Render("-- " + it.title))
			b.WriteString("\n")
			continue
		}
		marker := "  "
		style := normalRow
		if i == cursor {
			marker = "> "
			style = selectedRow
		} else if it.dim {
			style = dimRow
		}
		var row strings.Builder
		if it.badge == "[x]" {
			row.WriteString(badgeActive.Render(it.badge))
			row.WriteString(" ")
		} else if it.badge == "[ ]" {
			row.WriteString(badgeInactive.Render(it.badge))
			row.WriteString(" ")
		}
		row.WriteString(it.title)
		if it.desc != "" {
			row.WriteString("  ")
			row.WriteString(mutedStyle.Render(it.desc))
		}
		content := fit(row.String(), rowW-2)
		b.WriteString(style.Render(marker + content))
		b.WriteString("\n")
	}
	return b.String()
}
func fit(s string, w int) string {
	if w <= 0 || lipgloss.Width(s) <= w {
		return s
	}
	runes := []rune(s)
	var b strings.Builder
	for _, r := range runes {
		if lipgloss.Width(b.String()+string(r)) > w-3 {
			break
		}
		b.WriteRune(r)
	}
	return b.String() + "..."
}

type statusModel struct {
	text  string
	kind  statusKind
	until time.Time
}

func (s *statusModel) set(text string, k statusKind) {
	s.text = text
	s.kind = k
	s.until = time.Now().Add(6 * time.Second)
}
func (s *statusModel) expired() bool {
	return time.Now().After(s.until)
}
func (s *statusModel) View() string {
	if s.text == "" {
		return ""
	}
	var style = infoStyle
	switch s.kind {
	case statusSuccess:
		style = successStyle
	case statusError:
		style = errorStyle
	}
	return style.Render(s.text)
}

type confirmState int

const (
	confirmOpen confirmState = iota
	confirmYes
	confirmNo
)

type confirmModel struct {
	message string
	yes     bool
	state   confirmState
	onYes   func() tea.Cmd
}

func (c *confirmModel) Update(msg tea.Msg) (confirmModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "left", "right", "tab":
			c.yes = !c.yes
		case "enter":
			if c.yes {
				c.state = confirmYes
				if c.onYes != nil {
					return *c, c.onYes()
				}
				return *c, nil
			}
			c.state = confirmNo
			return *c, nil
		case "esc", "q":
			c.state = confirmNo
			return *c, nil
		}
	}
	return *c, nil
}
func (c *confirmModel) confirmed() bool {
	return c.state == confirmYes
}
func (c *confirmModel) dismissed() bool {
	return c.state == confirmNo
}
func (c *confirmModel) View(width int) string {
	yes := c.choice(assets.T("common.yes"), c.yes)
	no := c.choice(assets.T("common.no"), !c.yes)
	body := strings.TrimSpace(c.message) + "\n\n  " + yes + "   " + no
	box := overlayStyle.Render(body)
	return lipglossCenter(box, width)
}
func (c *confirmModel) choice(label string, selected bool) string {
	marker := "  "
	style := normalRow
	if selected {
		marker = "> "
		style = selectedRow
	}
	return style.Render(marker + "[" + label + "]")
}
func lipglossCenter(box string, width int) string {
	bw := lipgloss.Width(box)
	if bw >= width {
		return box
	}
	pad := strings.Repeat(" ", (width-bw)/2)
	lines := strings.Split(box, "\n")
	for i := range lines {
		lines[i] = pad + lines[i]
	}
	return strings.Join(lines, "\n")
}
func newInput(placeholder, value string, width int) (textinput.Model, tea.Cmd) {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = placeholder
	ti.Width = width
	if value != "" {
		ti.SetValue(value)
	}
	ti.PromptStyle = inputPromptStyle
	ti.TextStyle = inputTextStyle
	ti.PlaceholderStyle = inputPlaceholderStyle
	ti.Cursor.Style = inputCursorStyle
	return ti, ti.Focus()
}
func truncTo(s string, w int) string {
	if w <= 0 || lipgloss.Width(s) <= w {
		return s
	}
	runes := []rune(s)
	var b strings.Builder
	for _, r := range runes {
		if lipgloss.Width(b.String()+string(r)) > w {
			break
		}
		b.WriteRune(r)
	}
	return b.String()
}
func padTo(s string, w int) string {
	if w <= 0 {
		return s
	}
	return truncTo(s, w) + strings.Repeat(" ", w-lipgloss.Width(truncTo(s, w)))
}
func padFrame(content string, w, h int) string {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	out := make([]string, 0, h)
	for i := 0; i < h; i++ {
		line := ""
		if i < len(lines) {
			line = lines[i]
		}
		out = append(out, padTo(line, w))
	}
	return strings.Join(out, "\n")
}
func padToCol(content string, w, h int) string {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	out := make([]string, 0, h)
	for i := 0; i < h; i++ {
		line := ""
		if i < len(lines) {
			line = lines[i]
		}
		out = append(out, padTo(line, w))
	}
	return strings.Join(out, "\n")
}
func navCmd(s screen, payload interface{}) tea.Cmd {
	return func() tea.Msg { return navMsg{target: s, payload: payload} }
}
func backCmd(to screen, payload interface{}) tea.Cmd {
	return func() tea.Msg { return backMsg{to: to, payload: payload} }
}
func statusCmd(text string, k statusKind) tea.Cmd {
	return func() tea.Msg { return statusMsg{text: text, kind: k} }
}
func errCmd(err error) tea.Cmd {
	if err == nil {
		return nil
	}
	return func() tea.Msg { return errorMsg{err: err} }
}
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
func lipglossJoinHorizontal(left, right string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
