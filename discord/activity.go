package discord

import (
	"rich-presence-tui/config"
	"time"
)

type ActivityType int

const (
	ActivityTypeGame ActivityType = iota
	ActivityTypeStreaming
	ActivityTypeListening
	ActivityTypeWatching
	ActivityTypeCustom
	ActivityTypeCompeting
)

type Activity struct {
	Name       string              `json:"name,omitempty"`
	Type       ActivityType        `json:"type"`
	Details    string              `json:"details,omitempty"`
	State      string              `json:"state,omitempty"`
	Party      *ActivityParty      `json:"party,omitempty"`
	Assets     *ActivityAssets     `json:"assets,omitempty"`
	Buttons    []*ActivityButton   `json:"buttons,omitempty"`
	Timestamps *ActivityTimestamps `json:"timestamps,omitempty"`
}
type ActivityParty struct {
	ID   string `json:"id,omitempty"`
	Size [2]int `json:"size,omitempty"`
}
type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}
type ActivityButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}
type ActivityTimestamps struct {
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}

func BuildActivity(cfg config.ActivityConfig, now time.Time) Activity {
	act := Activity{
		Name:    cfg.Name,
		Type:    ActivityType(cfg.Type),
		Details: cfg.Details,
		State:   cfg.State,
	}
	if ts, ok := buildTimestamps(cfg, now); ok {
		act.Timestamps = ts
	}
	if cfg.PartySize > 0 || cfg.PartyMax > 0 {
		act.Party = &ActivityParty{
			Size: [2]int{cfg.PartySize, cfg.PartyMax},
		}
	}
	if cfg.LargeImage != "" || cfg.LargeText != "" || cfg.SmallImage != "" || cfg.SmallText != "" {
		act.Assets = &ActivityAssets{
			LargeImage: proxyImage(cfg.LargeImage),
			LargeText:  cfg.LargeText,
			SmallImage: proxyImage(cfg.SmallImage),
			SmallText:  cfg.SmallText,
		}
	}
	for _, b := range cfg.Buttons {
		if b.Label != "" && b.URL != "" && len(act.Buttons) < 2 {
			act.Buttons = append(act.Buttons, &ActivityButton{Label: b.Label, URL: b.URL})
		}
	}
	return act
}
func proxyImage(s string) string {
	return s
}
func buildTimestamps(cfg config.ActivityConfig, now time.Time) (*ActivityTimestamps, bool) {
	var ts ActivityTimestamps
	ok := false
	switch cfg.Start.Mode {
	case "elapsed":
		ts.Start = now.Unix()
		ok = true
	case "custom":
		if cfg.Start.Value > 0 {
			ts.Start = now.Add(-time.Duration(cfg.Start.Value) * time.Second).Unix()
			ok = true
		}
	}
	switch cfg.End.Mode {
	case "duration", "countdown":
		if cfg.End.Value > 0 {
			ts.End = now.Add(time.Duration(cfg.End.Value) * time.Second).Unix()
			ok = true
		}
	}
	if !ok {
		return nil, false
	}
	return &ts, true
}
