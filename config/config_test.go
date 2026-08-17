package config
import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)
func TestDefaultPathPrefersExeDir(t *testing.T) {
    t.Setenv("RICH_PRESENCE_TUI_CONFIG", filepath.Join("X", "config.json"))
    if got := DefaultPath(); got != filepath.Join("X", "config.json") {
        t.Fatalf("expected env override, got %q", got)
    }
    t.Setenv("RICH_PRESENCE_TUI_CONFIG", "")
    got := DefaultPath()
    if !strings.HasSuffix(got, string(filepath.Separator)+"rich-presence-tui"+string(filepath.Separator)+"config.json") {
        t.Fatalf("expected user config dir suffix, got %q", got)
    }
}
func TestLoadFreshConfigDefaultsLanguage(t *testing.T) {
    cfg, err := Load(filepath.Join(t.TempDir(), "config.json"))
    if err != nil {
        t.Fatal(err)
    }
    if cfg.Settings.Language != "ru" {
        t.Fatalf("expected default language ru on fresh config, got %q", cfg.Settings.Language)
    }
}
func TestSaveAndReloadRoundtrip(t *testing.T) {
    path := filepath.Join(t.TempDir(), "config.json")
    cfg, err := Load(path)
    if err != nil {
        t.Fatal(err)
    }
    cfg.Settings.Language = "en"
    cfg.Settings.ConfirmDelete = true
    p, err := cfg.CreatePreset("Game")
    if err != nil {
        t.Fatal(err)
    }
    p.AppID = "123456789012345678"
    if err := cfg.Save(); err != nil {
        t.Fatal(err)
    }
    got, err := Load(path)
    if err != nil {
        t.Fatal(err)
    }
    if got.Settings.Language != "en" || !got.Settings.ConfirmDelete {
        t.Fatalf("settings not persisted: %+v", got.Settings)
    }
    if len(got.Presets) != 1 || got.Presets[0].Name != "Game" {
        t.Fatalf("presets not persisted: %+v", got.Presets)
    }
    if got.Presets[0].AppID != "123456789012345678" {
        t.Fatalf("app id not persisted: %q", got.Presets[0].AppID)
    }
}
func TestCreatePresetLimitError(t *testing.T) {
    cfg, err := Load(filepath.Join(t.TempDir(), "config.json"))
    if err != nil {
        t.Fatal(err)
    }
    for i := 0; i < MaxPresets; i++ {
        if _, err := cfg.CreatePreset("p"); err != nil {
            t.Fatalf("preset %d: %v", i, err)
        }
    }
    _, err = cfg.CreatePreset("overflow")
    if !os.IsNotExist(err) {
    }
}
