package tui

import (
	"encoding/json"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"path/filepath"
	"rich-presence-tui/assets"
	"rich-presence-tui/config"
	"rich-presence-tui/discord"
	"strings"
	"testing"
	"time"
)

var now = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func newTestModel(t *testing.T) (Model, *config.Config) {
	t.Helper()
	dir := t.TempDir()
	cfg, err := config.Load(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	cfg.Settings.Language = "ru"
	return New(cfg, discord.NewPresence()), cfg
}
func keyMsg(s string) tea.KeyMsg {
	switch s {
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEscape}
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}
	case " ":
		return tea.KeyMsg{Type: tea.KeySpace}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}
func press(m Model, key string) Model {
	nm, cmd := m.Update(keyMsg(key))
	mm, ok := nm.(Model)
	if !ok {
		panic("root model Update returned a non-Model value")
	}
	return applyCmds(mm, cmd)
}
func typeText(m Model, s string) Model {
	return upd(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)})
}
func pressRune(m Model, r rune) Model {
	nm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	mm, ok := nm.(Model)
	if !ok {
		panic("root model Update returned a non-Model value")
	}
	return applyCmds(mm, cmd)
}
func upd(m Model, msg tea.Msg) Model {
	nm, _ := m.Update(msg)
	mm, ok := nm.(Model)
	if !ok {
		panic("root model Update returned a non-Model value")
	}
	return mm
}
func applyCmds(m Model, cmd tea.Cmd) Model {
	for _, msg := range execCmd(cmd) {
		m = upd(m, msg)
	}
	return m
}
func execCmd(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}
	ch := make(chan tea.Msg, 1)
	go func() { ch <- cmd() }()
	select {
	case msg := <-ch:
		return collectMsgs(msg)
	case <-time.After(300 * time.Millisecond):
		return nil
	}
}
func collectMsgs(msg tea.Msg) []tea.Msg {
	switch v := msg.(type) {
	case tea.BatchMsg:
		var out []tea.Msg
		for _, c := range v {
			out = append(out, execCmd(c)...)
		}
		return out
	default:
		return []tea.Msg{msg}
	}
}
func mustContain(t *testing.T, s, want string) {
	t.Helper()
	if !strings.Contains(s, want) {
		t.Fatalf("expected output to contain %q, got:\n%s", want, s)
	}
}
func sized(m Model) Model {
	return upd(m, tea.WindowSizeMsg{Width: 100, Height: 40})
}
func seedPresets(cfg *config.Config, names ...string) {
	for _, name := range names {
		if _, err := cfg.CreatePreset(name); err != nil {
			panic(err)
		}
	}
}
func TestMainMenuDashboardRenders(t *testing.T) {
	m, cfg := newTestModel(t)
	seedPresets(cfg, "Game", "Music")
	m = sized(New(cfg, discord.NewPresence()))
	v := m.View()
	mustContain(t, v, "Rich Presence TUI")
	mustContain(t, v, "Пресеты")
	mustContain(t, v, "Game")
	mustContain(t, v, "Music")
	mustContain(t, v, "Создать новый пресет")
	mustContain(t, v, "Настройки")
	mustContain(t, v, "Выход")
}
func TestLanguageToggleWithOneKey(t *testing.T) {
	m, cfg := newTestModel(t)
	m = sized(m)
	m = press(m, "l")
	if cfg.Settings.Language != "en" {
		t.Fatalf("expected language en after one L press, got %q", cfg.Settings.Language)
	}
	v := m.View()
	mustContain(t, v, "Presets")
	mustContain(t, v, "Settings")
	mustContain(t, v, "Create new preset")
	m = press(m, "l")
	if cfg.Settings.Language != "ru" {
		t.Fatalf("expected language ru after second L press, got %q", cfg.Settings.Language)
	}
}
func TestLanguageToggleFromSettingsRow(t *testing.T) {
	m, cfg := newTestModel(t)
	m = sized(m)
	m = press(m, "down")
	m = press(m, "down")
	m = press(m, "enter")
	mustContain(t, m.View(), "Настройки")
	m = press(m, "enter")
	if cfg.Settings.Language != "en" {
		t.Fatalf("expected language en after settings row toggle, got %q", cfg.Settings.Language)
	}
	mustContain(t, m.View(), "Settings")
}
func TestCreatePresetFromDashboard(t *testing.T) {
	m, cfg := newTestModel(t)
	m = sized(m)
	m = press(m, "n")
	mustContain(t, m.View(), "Название пресета")
	m = typeText(m, "My Game")
	m = press(m, "enter")
	if len(cfg.Presets) != 1 {
		t.Fatalf("expected 1 preset, got %d", len(cfg.Presets))
	}
	if cfg.Presets[0].Name != "My Game" {
		t.Fatalf("expected preset name My Game, got %q", cfg.Presets[0].Name)
	}
	v := m.View()
	mustContain(t, v, "Редактировать активность")
}
func TestSettingsDefaultAppIDReturnsToSettings(t *testing.T) {
	m, cfg := newTestModel(t)
	cfg.AppIDs = append(cfg.AppIDs, "123456789012345678")
	m = sized(m)
	m = press(m, "down")
	m = press(m, "down")
	m = press(m, "enter")
	mustContain(t, m.View(), "Настройки")
	for i := 0; i < 4; i++ {
		m = press(m, "down")
	}
	m = press(m, "enter")
	mustContain(t, m.View(), "Выбор Application ID")
	m = press(m, "esc")
	mustContain(t, m.View(), "Настройки")
	for i := 0; i < 4; i++ {
		m = press(m, "down")
	}
	m = press(m, "enter")
	m = press(m, "enter")
	mustContain(t, m.View(), "Настройки")
}
func TestSettingsMinimizeToTrayToggle(t *testing.T) {
	m, cfg := newTestModel(t)
	m = sized(m)
	m = press(m, "down")
	m = press(m, "down")
	m = press(m, "enter")
	m = press(m, "down")
	m = press(m, "down")
	m = press(m, "enter")
	if !cfg.Settings.MinimizeToTray {
		t.Fatal("expected MinimizeToTray to be enabled after toggle")
	}
	mustContain(t, m.View(), "Вкл")
}
func TestSettingsFullResetDeletesConfig(t *testing.T) {
	m, cfg := newTestModel(t)
	if _, err := cfg.CreatePreset("Preset A"); err != nil {
		t.Fatal(err)
	}
	cfg.AppIDs = append(cfg.AppIDs, "111111111111111111")
	m = sized(m)
	m = press(m, "down")
	m = press(m, "down")
	m = press(m, "enter")
	mustContain(t, m.View(), assets.T("menu.resetConfirm"))
	m = press(m, "enter")
	if len(cfg.Presets) != 0 {
		t.Fatalf("expected presets to be cleared, got %d", len(cfg.Presets))
	}
	if len(cfg.AppIDs) != 0 {
		t.Fatalf("expected app ids to be cleared, got %d", len(cfg.AppIDs))
	}
	if _, err := os.Stat(cfg.Path()); !os.IsNotExist(err) {
		t.Fatalf("expected config file to be deleted, stat err = %v", err)
	}
	mustContain(t, m.View(), assets.T("menu.section"))
}
func TestEditTypeSelectsViaOptionsPanel(t *testing.T) {
	_, cfg := newTestModel(t)
	if _, err := cfg.CreatePreset("P"); err != nil {
		t.Fatal(err)
	}
	id := cfg.Presets[0].ID
	em := newEditModel(appCtx{cfg: cfg, presence: discord.NewPresence()}, id)
	em.focus = panelFields
	em.cursor = 0
	if cmd := em.Update(tea.KeyMsg{Type: tea.KeyRight}); cmd != nil {
		t.Fatal("opening the options panel should not emit a command")
	}
	if em.focus != panelOptions {
		t.Fatalf("expected panelOptions, got %d", em.focus)
	}
	em.Update(tea.KeyMsg{Type: tea.KeyDown})
	if cmd := em.Update(tea.KeyMsg{Type: tea.KeyEnter}); cmd == nil {
		t.Fatal("expected a command after applying an option")
	} else if msg := cmd(); msg != nil {
		_ = em.Update(msg)
	}
	if cfg.Presets[0].Activity.Type != config.TypeListening {
		t.Fatalf("expected listening type, got %d", cfg.Presets[0].Activity.Type)
	}
	if cmd := em.Update(tea.KeyMsg{Type: tea.KeyRight}); cmd != nil {
		t.Fatal("unexpected command")
	}
	if cmd := em.Update(tea.KeyMsg{Type: tea.KeyLeft}); cmd != nil {
		t.Fatal("cancelling should not emit a command")
	}
	if em.focus != panelFields {
		t.Fatalf("expected focus back on fields, got %d", em.focus)
	}
	if cfg.Presets[0].Activity.Type != config.TypeListening {
		t.Fatalf("cancel changed the type to %d", cfg.Presets[0].Activity.Type)
	}
}
func TestEditorOmitsBrokenFeatures(t *testing.T) {
	_, cfg := newTestModel(t)
	if _, err := cfg.CreatePreset("P"); err != nil {
		t.Fatal(err)
	}
	em := newEditModel(appCtx{cfg: cfg, presence: discord.NewPresence()}, cfg.Presets[0].ID)
	if len(em.cats) != 6 {
		t.Fatalf("expected 6 categories, got %d", len(em.cats))
	}
	v := em.View(110, 40)
	for _, banned := range []string{"Секреты", "Флаги", "Платформы", "Party ID", "URL (для стрима)", "Emoji", "Join Secret", "Instance"} {
		if strings.Contains(v, banned) {
			t.Fatalf("editor still exposes removed feature %q", banned)
		}
	}
}
func TestEditorPartyShowsBothFields(t *testing.T) {
	_, cfg := newTestModel(t)
	if _, err := cfg.CreatePreset("P"); err != nil {
		t.Fatal(err)
	}
	em := newEditModel(appCtx{cfg: cfg, presence: discord.NewPresence()}, cfg.Presets[0].ID)
	var party *categoryDef
	for i := range em.cats {
		if em.cats[i].title == assets.T("edit.cat.party") {
			party = &em.cats[i]
			break
		}
	}
	if party == nil {
		t.Fatal("party category not found in editor")
	}
	if len(party.fields) != 2 {
		t.Fatalf("expected 2 party fields, got %d", len(party.fields))
	}
	got := em.renderFields(40, 20)
	for _, want := range []string{assets.T("edit.f.partySize"), assets.T("edit.f.partyMax")} {
		if !strings.Contains(got, want) {
			t.Fatalf("party fields missing %q in render:\n%s", want, got)
		}
	}
}
func TestEditorSingleHighlight(t *testing.T) {
	_, cfg := newTestModel(t)
	if _, err := cfg.CreatePreset("P"); err != nil {
		t.Fatal(err)
	}
	em := newEditModel(appCtx{cfg: cfg, presence: discord.NewPresence()}, cfg.Presets[0].ID)
	em.focus = panelFields
	em.cursor = 0
	if got := strings.Count(em.renderFields(38, 12), "> "); got != 1 {
		t.Fatalf("fields column must highlight exactly one row, got %d", got)
	}
	if got := strings.Count(em.renderOptions(30, 12), "> "); got != 0 {
		t.Fatalf("options column must not highlight while fields are focused, got %d", got)
	}
	em.focus = panelOptions
	if got := strings.Count(em.renderOptions(30, 12), "> "); got != 1 {
		t.Fatalf("options column must highlight exactly one option, got %d", got)
	}
}
func TestEditorBackRowReturns(t *testing.T) {
	_, cfg := newTestModel(t)
	if _, err := cfg.CreatePreset("P"); err != nil {
		t.Fatal(err)
	}
	em := newEditModel(appCtx{cfg: cfg, presence: discord.NewPresence()}, cfg.Presets[0].ID)
	em.focus = panelFields
	em.cursor = len(em.fields)
	cmd := em.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a back command")
	}
	bm, ok := cmd().(backMsg)
	if !ok || bm.to != screenPresetActions {
		t.Fatalf("expected backMsg to presetActions, got %#v", cmd())
	}
}
func TestTimeCustomDurationAppears(t *testing.T) {
	_, cfg := newTestModel(t)
	if _, err := cfg.CreatePreset("P"); err != nil {
		t.Fatal(err)
	}
	em := newEditModel(appCtx{cfg: cfg, presence: discord.NewPresence()}, cfg.Presets[0].ID)
	idx := -1
	for i, f := range em.fields {
		if f.def.id == "end" {
			idx = i
		}
	}
	if idx < 0 {
		t.Fatal("end field not found")
	}
	em.cursor = idx
	em.Update(tea.KeyMsg{Type: tea.KeyRight})
	em.Update(tea.KeyMsg{Type: tea.KeyDown})
	if cmd := em.Update(tea.KeyMsg{Type: tea.KeyEnter}); cmd != nil {
		for _, msg := range execCmd(cmd) {
			_ = em.Update(msg)
		}
	}
	if !em.editing {
		t.Fatal("expected the duration editor to open in the options panel")
	}
	if em.durEdit != "endDuration" {
		t.Fatalf("durEdit = %q, want endDuration", em.durEdit)
	}
	if cfg.Presets[0].Activity.End.Mode != "duration" {
		t.Fatalf("end mode = %q, want duration", cfg.Presets[0].Activity.End.Mode)
	}
	if secs, ok := parseDuration("100:30:00"); !ok || secs != 361800 {
		t.Fatalf("parseDuration(100:30:00) = %d, %v", secs, ok)
	}
	if secs, ok := parseDuration("999999:59:59"); !ok || secs != 999999*3600+3599 {
		t.Fatalf("parseDuration(999999:59:59) = %d, %v", secs, ok)
	}
}
func TestTimeStartCustomAppears(t *testing.T) {
	_, cfg := newTestModel(t)
	if _, err := cfg.CreatePreset("P"); err != nil {
		t.Fatal(err)
	}
	em := newEditModel(appCtx{cfg: cfg, presence: discord.NewPresence()}, cfg.Presets[0].ID)
	idx := -1
	for i, f := range em.fields {
		if f.def.id == "start" {
			idx = i
		}
	}
	if idx < 0 {
		t.Fatal("start field not found")
	}
	em.cursor = idx
	em.Update(tea.KeyMsg{Type: tea.KeyRight})
	em.Update(tea.KeyMsg{Type: tea.KeyDown})
	em.Update(tea.KeyMsg{Type: tea.KeyDown})
	if cmd := em.Update(tea.KeyMsg{Type: tea.KeyEnter}); cmd != nil {
		for _, msg := range execCmd(cmd) {
			_ = em.Update(msg)
		}
	}
	if !em.editing {
		t.Fatal("expected the duration editor to open in the options panel")
	}
	if em.durEdit != "startDuration" {
		t.Fatalf("durEdit = %q, want startDuration", em.durEdit)
	}
	if cfg.Presets[0].Activity.Start.Mode != "custom" {
		t.Fatalf("start mode = %q, want custom", cfg.Presets[0].Activity.Start.Mode)
	}
}
func TestTimeCustomValueSavedFromOptionsPanel(t *testing.T) {
	_, cfg := newTestModel(t)
	if _, err := cfg.CreatePreset("P"); err != nil {
		t.Fatal(err)
	}
	em := newEditModel(appCtx{cfg: cfg, presence: discord.NewPresence()}, cfg.Presets[0].ID)
	for i, f := range em.fields {
		if f.def.id == "start" {
			em.cursor = i
		}
	}
	em.Update(tea.KeyMsg{Type: tea.KeyRight})
	em.Update(tea.KeyMsg{Type: tea.KeyDown})
	em.Update(tea.KeyMsg{Type: tea.KeyDown})
	if cmd := em.Update(tea.KeyMsg{Type: tea.KeyEnter}); cmd != nil {
		for _, msg := range execCmd(cmd) {
			_ = em.Update(msg)
		}
	}
	if !em.editing || em.durEdit != "startDuration" {
		t.Fatalf("editor not open for startDuration (editing=%v durEdit=%q)", em.editing, em.durEdit)
	}
	em.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2:30:00")})
	if cmd := em.Update(tea.KeyMsg{Type: tea.KeyEnter}); cmd != nil {
		for _, msg := range execCmd(cmd) {
			_ = em.Update(msg)
		}
	}
	if em.editing {
		t.Fatal("editor should be closed after entering the time")
	}
	if em.durEdit != "" {
		t.Fatalf("durEdit = %q, want empty", em.durEdit)
	}
	if cfg.Presets[0].Activity.Start.Mode != "custom" {
		t.Fatalf("start mode = %q, want custom", cfg.Presets[0].Activity.Start.Mode)
	}
	if cfg.Presets[0].Activity.Start.Value != 2*3600+30*60 {
		t.Fatalf("start value = %d, want %d", cfg.Presets[0].Activity.Start.Value, 2*3600+30*60)
	}
	if got := em.renderFields(40, 20); !strings.Contains(got, assets.T("edit.time.duration")+" (2:30:00)") {
		t.Fatalf("start select should show the custom time, got:\n%s", got)
	}
}
func TestHotkeysWorkRegardlessOfLayout(t *testing.T) {
	m, cfg := newTestModel(t)
	seedPresets(cfg, "Game")
	m = sized(New(cfg, discord.NewPresence()))
	m = pressRune(m, 'у')
	mustContain(t, m.View(), "Редактирование")
}
func TestFixedLayoutAndFooter(t *testing.T) {
	m, _ := newTestModel(t)
	m = sized(m)
	v := m.View()
	if strings.Count(v, "\n") != m.height-1 {
		t.Fatalf("expected %d lines after sizing, got %d", m.height, strings.Count(v, "\n")+1)
	}
	mustContain(t, v, "остановлен")
	mustContain(t, v, "выбор")
	m = press(m, "down")
	m = press(m, "down")
	m = press(m, "enter")
	v2 := m.View()
	if strings.Count(v2, "\n") != m.height-1 {
		t.Fatalf("layout jumped to %d lines", strings.Count(v2, "\n")+1)
	}
	mustContain(t, v2, "Настройки")
	mustContain(t, v2, "вкл/выкл")
}
func TestParseAndFormatDuration(t *testing.T) {
	cases := []struct {
		in   string
		want int
		ok   bool
	}{
		{"90", 90, true},
		{"5:30", 330, true},
		{"1:30:00", 5400, true},
		{"", 0, false},
		{"abc", 0, false},
	}
	for _, c := range cases {
		got, ok := parseDuration(c.in)
		if got != c.want || ok != c.ok {
			t.Fatalf("parseDuration(%q) = %d,%v want %d,%v", c.in, got, ok, c.want, c.ok)
		}
	}
	if formatDuration(90) != "01:30" {
		t.Fatalf("formatDuration(90) = %q, want 01:30", formatDuration(90))
	}
	if formatDuration(5400) != "1:30:00" {
		t.Fatalf("formatDuration(5400) = %q, want 1:30:00", formatDuration(5400))
	}
}
func TestLeftArrowBacksOutOfSubscreens(t *testing.T) {
	m, _ := newTestModel(t)
	m = sized(m)
	m = press(m, "down")
	m = press(m, "down")
	m = press(m, "enter")
	mustContain(t, m.View(), "Настройки")
	m = press(m, "left")
	mustContain(t, m.View(), "Пресеты")
}
func TestBuildActivityConvertsAllFields(t *testing.T) {
	act := config.DefaultActivity()
	act.Name = "Game"
	act.Details = "playing"
	act.State = "ranked"
	act.Type = config.TypeCompeting
	act.LargeImage = "big"
	act.LargeText = "big tooltip"
	act.SmallImage = "small"
	act.SmallText = "small tooltip"
	act.Start = config.Timestamp{Mode: "elapsed"}
	act.End = config.Timestamp{Mode: "countdown", Value: 600}
	act.PartySize = 2
	act.PartyMax = 4
	act.Buttons[0] = config.ButtonConfig{Label: "Join", URL: "https://x.com"}
	a := discord.BuildActivity(act, now)
	if a.Name != "Game" || a.Details != "playing" || a.State != "ranked" {
		t.Fatalf("text fields not converted: %+v", a)
	}
	if a.Assets == nil || a.Assets.LargeImage != "big" || a.Assets.SmallText != "small tooltip" {
		t.Fatalf("assets not converted: %+v", a.Assets)
	}
	if a.Timestamps == nil || a.Timestamps.Start == 0 || a.Timestamps.End != now.Add(600*time.Second).Unix() {
		t.Fatalf("timestamps not converted: %+v", a.Timestamps)
	}
	if a.Party == nil || a.Party.Size != [2]int{2, 4} {
		t.Fatalf("party not converted: %+v", a.Party)
	}
	if len(a.Buttons) != 1 || a.Buttons[0].Label != "Join" {
		t.Fatalf("buttons not converted: %+v", a.Buttons)
	}
	if a.Type != 5 {
		t.Fatalf("type not converted: %d", a.Type)
	}
}
func TestBuildActivitySendsImageURLsVerbatim(t *testing.T) {
	act := config.DefaultActivity()
	act.LargeImage = "https://example.com/large.png"
	act.SmallImage = "https://example.com/small.png"
	a := discord.BuildActivity(act, now)
	if a.Assets == nil {
		t.Fatal("expected assets")
	}
	if a.Assets.LargeImage != "https://example.com/large.png" {
		t.Fatalf("large image changed: %q", a.Assets.LargeImage)
	}
	if a.Assets.SmallImage != "https://example.com/small.png" {
		t.Fatalf("small image changed: %q", a.Assets.SmallImage)
	}
}
func TestBuildActivityKeepsAssetKeysAndExistingPrefix(t *testing.T) {
	act := config.DefaultActivity()
	act.LargeImage = "numbani_map"
	act.SmallImage = "mp:https://example.com/already.png"
	a := discord.BuildActivity(act, now)
	if a.Assets == nil {
		t.Fatal("expected assets")
	}
	if a.Assets.LargeImage != "numbani_map" {
		t.Fatalf("asset key must not be prefixed: %q", a.Assets.LargeImage)
	}
	if a.Assets.SmallImage != "mp:https://example.com/already.png" {
		t.Fatalf("already-prefixed image must not be double-prefixed: %q", a.Assets.SmallImage)
	}
}
func TestBuildActivityOmitsInviteFields(t *testing.T) {
	act := config.DefaultActivity()
	a := discord.BuildActivity(act, now)
	raw, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	for _, banned := range []string{"\"instance\"", "\"flags\"", "\"supported_platforms\"", "\"url\"", "\"emoji\"", "\"secrets\"", "\"party\""} {
		if strings.Contains(string(raw), banned) {
			t.Fatalf("invite-related field %s must not be sent: %s", banned, raw)
		}
	}
}
func TestBuildActivityCustomStart(t *testing.T) {
	act := config.DefaultActivity()
	act.Start = config.Timestamp{Mode: "custom", Value: 3600}
	act.End = config.Timestamp{Mode: "duration", Value: 600}
	a := discord.BuildActivity(act, now)
	if a.Timestamps == nil {
		t.Fatal("expected timestamps")
	}
	if a.Timestamps.Start != now.Add(-time.Hour).Unix() {
		t.Fatalf("custom start = %d, want %d", a.Timestamps.Start, now.Add(-time.Hour).Unix())
	}
	if a.Timestamps.End != now.Add(600*time.Second).Unix() {
		t.Fatalf("end = %d, want %d", a.Timestamps.End, now.Add(600*time.Second).Unix())
	}
}
func TestMainMenuRightOpensPresetNotToggle(t *testing.T) {
	m, cfg := newTestModel(t)
	seedPresets(cfg, "Game")
	cfg.Presets[0].Enabled = true
	m = sized(New(cfg, discord.NewPresence()))
	m = press(m, "right")
	mustContain(t, m.View(), "Пресет: Game")
	if activePreset(cfg) == nil {
		t.Fatal("right arrow must not toggle the preset off")
	}
}
func TestHelpScreenOpensAndCloses(t *testing.T) {
	m, _ := newTestModel(t)
	m = sized(m)
	m = pressRune(m, '?')
	mustContain(t, m.View(), "Помощь")
	mustContain(t, m.View(), "Статус активности")
	m = press(m, "esc")
	mustContain(t, m.View(), "Пресеты")
}
func TestEditorHintLivesInFooterNotColumns(t *testing.T) {
	m, cfg := newTestModel(t)
	seedPresets(cfg, "Game")
	m = sized(New(cfg, discord.NewPresence()))
	m = pressRune(m, 'e')
	v := m.View()
	if strings.Contains(v, "[→] ок") {
		t.Fatal("select hint must not render inside the editor columns")
	}
	mustContain(t, v, "[↑/↓] поле")
}
func TestLanguageNameOfActiveDict(t *testing.T) {
	assets.SetLanguage("ru")
	if assets.LanguageName(assets.Languages()[0]) == "" {
		t.Fatal("expected a language name")
	}
}
