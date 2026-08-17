package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	MaxPresets = 10
)

type Settings struct {
	Language       string `json:"language"`
	ActiveAppID    string `json:"active_app_id"`
	ConfirmDelete  bool   `json:"confirm_delete"`
	Autostart      bool   `json:"autostart"`
	MinimizeToTray bool   `json:"minimize_to_tray"`
}
type Config struct {
	Settings Settings  `json:"settings"`
	Presets  []*Preset `json:"presets"`
	AppIDs   []string  `json:"app_ids"`
	path     string
}
type Preset struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	AppID    string         `json:"app_id"`
	Enabled  bool           `json:"enabled"`
	Activity ActivityConfig `json:"activity"`
}
type ActivityType int

const (
	TypeGame ActivityType = iota
	TypeStreaming
	TypeListening
	TypeWatching
	TypeCustom
	TypeCompeting
)

type ActivityConfig struct {
	Type       ActivityType   `json:"type"`
	Name       string         `json:"name"`
	State      string         `json:"state"`
	Details    string         `json:"details"`
	LargeImage string         `json:"large_image"`
	LargeText  string         `json:"large_text"`
	SmallImage string         `json:"small_image"`
	SmallText  string         `json:"small_text"`
	Start      Timestamp      `json:"start"`
	End        Timestamp      `json:"end"`
	PartySize  int            `json:"party_size"`
	PartyMax   int            `json:"party_max"`
	Buttons    []ButtonConfig `json:"buttons"`
}
type Timestamp struct {
	Mode  string `json:"mode"`
	Value int64  `json:"value,omitempty"`
}
type ButtonConfig struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

func DefaultActivity() ActivityConfig {
	return ActivityConfig{
		Type:    TypeGame,
		Buttons: []ButtonConfig{{}, {}},
	}
}

var ErrPresetLimit = errors.New("preset limit reached")

func DefaultPath() string {
	if p := os.Getenv("RICH_PRESENCE_TUI_CONFIG"); p != "" {
		return p
	}
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "rich-presence-tui", "config.json")
	}
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), "config.json")
	}
	return "config.json"
}
func Load(path string) (*Config, error) {
	cfg := &Config{path: path}
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if cfg.Settings.Language == "" {
		cfg.Settings.Language = "ru"
	}
	return cfg, nil
}
func (c *Config) Save() error {
	if c.path == "" {
		c.path = DefaultPath()
	}
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	tmp := c.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, c.path)
}
func (c *Config) Path() string {
	if c.path == "" {
		c.path = DefaultPath()
	}
	return c.path
}
func (c *Config) Delete() error {
	path := c.Path()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	_ = os.Remove(path + ".tmp")
	return nil
}
func (c *Config) Preset(id string) *Preset {
	for _, p := range c.Presets {
		if p.ID == id {
			return p
		}
	}
	return nil
}
func (c *Config) PresetByName(name string) *Preset {
	for _, p := range c.Presets {
		if p.Name == name {
			return p
		}
	}
	return nil
}
func (c *Config) Full() bool {
	return len(c.Presets) >= MaxPresets
}
func (c *Config) CreatePreset(name string) (*Preset, error) {
	if c.Full() {
		return nil, fmt.Errorf("%w: %d", ErrPresetLimit, MaxPresets)
	}
	p := &Preset{
		ID:       fmt.Sprintf("p%d", time.Now().UnixNano()),
		Name:     name,
		AppID:    c.Settings.ActiveAppID,
		Activity: DefaultActivity(),
	}
	c.Presets = append(c.Presets, p)
	return p, nil
}
func (c *Config) DeletePreset(id string) {
	for i, p := range c.Presets {
		if p.ID == id {
			c.Presets = append(c.Presets[:i], c.Presets[i+1:]...)
			return
		}
	}
}
func (c *Config) AddAppID(id string) {
	for _, v := range c.AppIDs {
		if v == id {
			return
		}
	}
	c.AppIDs = append(c.AppIDs, id)
}
func (c *Config) RemoveAppID(id string) {
	for i, v := range c.AppIDs {
		if v == id {
			c.AppIDs = append(c.AppIDs[:i], c.AppIDs[i+1:]...)
			return
		}
	}
}
