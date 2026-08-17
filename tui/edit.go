package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"rich-presence-tui/assets"
	"rich-presence-tui/config"
	"strconv"
	"strings"
)

type focusPanel int

const (
	panelFields focusPanel = iota
	panelOptions
)

type editModel struct {
	app     appCtx
	preset  *config.Preset
	act     *config.ActivityConfig
	cats    []categoryDef
	fields  []formField
	catOf   []int
	cursor  int
	focus   focusPanel
	optCur  int
	editing bool
	editor  textinput.Model
	durEdit string
	scroll  int
}

func newEditModel(app appCtx, id string) *editModel {
	m := &editModel{app: app}
	m.preset = app.cfg.Preset(id)
	if m.preset == nil {
		return m
	}
	m.act = &m.preset.Activity
	m.cats = buildCategories()
	m.buildFields("")
	return m
}
func (m *editModel) buildFields(keep string) {
	var fields []formField
	var catOf []int
	for ci, cat := range m.cats {
		for _, d := range cat.fields {
			if d.visible != nil && !d.visible(m.act) {
				continue
			}
			fields = append(fields, formField{def: d, value: d.get(m.act)})
			catOf = append(catOf, ci)
		}
	}
	m.fields = fields
	m.catOf = catOf
	if keep != "" {
		for i, f := range m.fields {
			if f.def.id == keep {
				m.cursor = i
				break
			}
		}
	} else {
		m.cursor = 0
	}
	if m.cursor > len(m.fields) {
		m.cursor = len(m.fields)
	}
	if m.focus != panelFields {
		m.focus = panelFields
	}
}
func (m *editModel) langOK() bool {
	return m.preset != nil && !m.editing
}
func (m *editModel) currentField() *formField {
	if m.cursor < len(m.fields) {
		return &m.fields[m.cursor]
	}
	return nil
}
func (m *editModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case fieldChangedMsg:
		for _, cat := range m.cats {
			for _, d := range cat.fields {
				if d.id == msg.id {
					d.set(m.act, msg.value)
					break
				}
			}
		}
		m.buildFields(msg.id)
		if err := m.app.cfg.Save(); err != nil {
			return errCmd(err)
		}
		cmds := []tea.Cmd{statusCmd(assets.T("status.saved"), statusSuccess)}
		if m.preset != nil && m.preset.Enabled {
			cmds = append(cmds, refreshPresetCmd(m.app, m.preset))
		}
		return tea.Batch(cmds...)
	}
	if key, ok := msg.(tea.KeyMsg); ok {
		if m.preset == nil {
			if normKey(key) == "esc" {
				return backCmd(screenPresetActions, nil)
			}
			return nil
		}
		if m.editing {
			return m.updateEditing(key)
		}
		if m.focus == panelOptions {
			return m.updateOptions(key)
		}
		return m.updateFields(key)
	}
	return nil
}
func (m *editModel) updateFields(key tea.KeyMsg) tea.Cmd {
	switch normKey(key) {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.fields) {
			m.cursor++
		}
	case "enter", "right":
		return m.activateField()
	case " ":
		if f := m.currentField(); f != nil && f.def.kind == kindBool {
			return m.toggle(m.cursor)
		}
	case "left", "esc":
		return backCmd(screenPresetActions, m.preset.ID)
	}
	return nil
}
func (m *editModel) activateField() tea.Cmd {
	if m.cursor >= len(m.fields) {
		return backCmd(screenPresetActions, m.preset.ID)
	}
	f := &m.fields[m.cursor]
	switch f.def.kind {
	case kindText, kindNumber:
		ti, cmd := newInput("", f.value, 28)
		m.editing = true
		m.editor = ti
		return cmd
	case kindBool:
		return m.toggle(m.cursor)
	case kindSelect:
		m.focus = panelOptions
		m.optCur = optionIndex(f.def.options, f.value)
		return nil
	}
	return nil
}
func (m *editModel) toggle(idx int) tea.Cmd {
	v := m.fields[idx].value == "true"
	m.fields[idx].value = strconv.FormatBool(!v)
	id := m.fields[idx].def.id
	return func() tea.Msg { return fieldChangedMsg{id: id, value: m.fields[idx].value} }
}
func (m *editModel) updateOptions(key tea.KeyMsg) tea.Cmd {
	f := m.currentField()
	if f == nil || f.def.kind != kindSelect {
		m.focus = panelFields
		return nil
	}
	opts := f.def.options
	switch normKey(key) {
	case "up", "k":
		if m.optCur > 0 {
			m.optCur--
		}
	case "down", "j":
		if m.optCur < len(opts)-1 {
			m.optCur++
		}
	case "enter", "right":
		m.fields[m.cursor].value = opts[m.optCur].value
		id := f.def.id
		value := opts[m.optCur].value
		if durID := durationFieldID(id, value); durID != "" {
			m.durEdit = durID
			m.editing = true
			ti, cmd := newInput("", durationValue(m.act, durID), 28)
			m.editor = ti
			return tea.Batch(cmd, func() tea.Msg { return fieldChangedMsg{id: id, value: value} })
		}
		m.focus = panelFields
		return func() tea.Msg { return fieldChangedMsg{id: id, value: value} }
	case "left", "esc":
		m.focus = panelFields
	}
	return nil
}
func (m *editModel) updateEditing(key tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	m.editor, cmd = m.editor.Update(key)
	switch normKey(key) {
	case "enter":
		value := m.editor.Value()
		id := m.fields[m.cursor].def.id
		if m.durEdit != "" {
			id = m.durEdit
		}
		m.editing = false
		m.editor = textinput.Model{}
		m.durEdit = ""
		changed := fieldChangedMsg{id: id, value: value}
		return tea.Batch(cmd, func() tea.Msg { return changed })
	case "esc":
		m.editing = false
		m.editor = textinput.Model{}
		m.durEdit = ""
	}
	return cmd
}
func durationFieldID(fieldID, value string) string {
	if fieldID == "start" && value == "custom" {
		return "startDuration"
	}
	if fieldID == "end" && value == "duration" {
		return "endDuration"
	}
	return ""
}
func durationValue(a *config.ActivityConfig, id string) string {
	var secs int64
	switch id {
	case "startDuration":
		secs = a.Start.Value
	case "endDuration":
		secs = a.End.Value
	}
	if secs <= 0 {
		return ""
	}
	return formatDuration(int(secs))
}
func (m *editModel) View(width, height int) string {
	if m.preset == nil {
		return mutedStyle.Render(assets.T("status.notFound"))
	}
	fieldW := 34
	optW := width - fieldW - 1
	if optW < 24 {
		optW = 24
	}
	var b strings.Builder
	b.WriteString(sectionStyle.Render(fmt.Sprintf(assets.T("edit.title"), m.preset.Name)))
	b.WriteString("\n\n")
	colH := height - 2
	if colH < 8 {
		colH = 8
	}
	b.WriteString(lipglossJoinHorizontal(
		m.renderFields(fieldW, colH),
		" "+m.renderOptions(optW, colH),
	))
	return b.String()
}
func (m *editModel) renderFields(w, h int) string {
	rows := m.fieldRows(w)
	visible := h
	if m.cursor < m.scroll {
		m.scroll = m.cursor
	}
	if m.cursor >= m.scroll+visible {
		m.scroll = m.cursor - visible + 1
	}
	if m.scroll > len(rows)-visible {
		m.scroll = len(rows) - visible
	}
	if m.scroll < 0 {
		m.scroll = 0
	}
	out := make([]string, 0, h)
	for i := m.scroll; i < len(rows) && len(out) < h; i++ {
		out = append(out, rows[i])
	}
	for len(out) < h {
		out = append(out, strings.Repeat(" ", w))
	}
	return strings.Join(out, "\n")
}
func (m *editModel) fieldRows(w int) []string {
	var rows []string
	for i := range m.fields {
		if i == 0 || m.catOf[i] != m.catOf[i-1] {
			rows = append(rows, dimStyle.Render("-- "+m.cats[m.catOf[i]].title))
		}
		rows = append(rows, padTo(m.renderField(m.fields[i], i == m.cursor && m.focus == panelFields, w), w))
	}
	marker := "  "
	style := hintTextStyle
	if m.cursor >= len(m.fields) {
		marker = "> "
		style = selectedRow
	}
	rows = append(rows, padTo(style.Render(marker+"<- "+assets.T("edit.back")), w))
	return rows
}
func (m *editModel) renderOptions(w, h int) string {
	var b strings.Builder
	header := assets.T("edit.panelHint")
	if m.durEdit != "" {
		if d, ok := fieldDefByID(m.cats, m.durEdit); ok {
			header = d.label
		}
	} else if f := m.currentField(); f != nil {
		header = f.def.label
	}
	b.WriteString(sectionStyle.Render(header))
	b.WriteString("\n\n")
	switch {
	case m.editing:
		b.WriteString("  " + inputPromptStyle.Render("> ") + m.editor.View())
	case m.cursor >= len(m.fields):
		b.WriteString(hintTextStyle.Render("<- " + assets.T("edit.back")))
	default:
		f := m.currentField()
		switch f.def.kind {
		case kindSelect:
			for i, o := range f.def.options {
				marker := "  "
				style := normalRow
				if m.focus == panelOptions && i == m.optCur {
					marker = "> "
					style = selectedRow
				}
				b.WriteString(style.Render(marker + o.label))
				b.WriteString("\n")
			}
		case kindBool:
			value := boolOffStyle.Render("[" + assets.T("common.off") + "]")
			if f.value == "true" {
				value = boolOnStyle.Render("[" + assets.T("common.on") + "]")
			}
			b.WriteString(value)
			b.WriteString("\n\n")
		default:
			if f.value == "" {
				b.WriteString(mutedStyle.Render("-"))
			} else {
				b.WriteString(fieldValueStyle.Render(truncTo(f.value, w)))
			}
		}
	}
	return padToCol(b.String(), w, h)
}
func (m *editModel) hint() string {
	if m.preset == nil {
		return ""
	}
	if m.editing {
		return assets.T("edit.hint.edit")
	}
	if m.cursor >= len(m.fields) {
		return assets.T("edit.hint.back")
	}
	if m.focus == panelOptions {
		return assets.T("edit.hint.select")
	}
	f := m.currentField()
	if f == nil {
		return ""
	}
	switch f.def.kind {
	case kindBool:
		return assets.T("edit.hint.space")
	case kindSelect:
		return assets.T("edit.hint.fields")
	default:
		return assets.T("edit.hint.enter")
	}
}

type formField struct {
	def   fieldDef
	value string
}

func (m *editModel) renderField(field formField, selected bool, width int) string {
	def := field.def
	marker := "  "
	if selected {
		marker = "> "
	}
	labelStyle := fieldLabelStyle
	if selected {
		labelStyle = fieldSelectedStyle
	}
	label := marker + labelStyle.Render(def.label+":")
	var value string
	switch def.kind {
	case kindBool:
		if field.value == "true" {
			value = boolOnStyle.Render("[" + assets.T("common.on") + "]")
		} else {
			value = boolOffStyle.Render("[" + assets.T("common.off") + "]")
		}
	case kindSelect:
		label := optionLabel(def.options, field.value)
		if extra := m.timeDisplay(field); extra != "" {
			label += " (" + extra + ")"
		}
		value = fieldValueStyle.Render(label)
	default:
		if field.value == "" {
			value = mutedStyle.Render("-")
		} else {
			value = fieldValueStyle.Render(truncTo(field.value, width-lipgloss.Width(label)-1))
		}
	}
	return label + " " + value
}
func (m *editModel) timeDisplay(field formField) string {
	if m.act == nil {
		return ""
	}
	switch field.def.id {
	case "start":
		if m.act.Start.Mode == "custom" && m.act.Start.Value > 0 {
			return formatDuration(int(m.act.Start.Value))
		}
	case "end":
		if m.act.End.Mode == "duration" && m.act.End.Value > 0 {
			return formatDuration(int(m.act.End.Value))
		}
	}
	return ""
}
func optionLabel(opts []option, value string) string {
	for _, o := range opts {
		if o.value == value {
			return o.label
		}
	}
	return ""
}
func fieldDefByID(cats []categoryDef, id string) (fieldDef, bool) {
	for _, cat := range cats {
		for _, d := range cat.fields {
			if d.id == id {
				return d, true
			}
		}
	}
	return fieldDef{}, false
}
func optionIndex(opts []option, value string) int {
	for i, o := range opts {
		if o.value == value {
			return i
		}
	}
	return 0
}
