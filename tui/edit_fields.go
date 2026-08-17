package tui
import (
	"fmt"
	"strconv"
	"strings"
	"rich-presence-tui/assets"
	"rich-presence-tui/config"
)
type fieldKind int
const (
	kindText fieldKind = iota
	kindNumber
	kindBool
	kindSelect
)
type fieldDef struct {
	id      string
	label   string
	kind    fieldKind
	options []option
	visible func(*config.ActivityConfig) bool
	get     func(*config.ActivityConfig) string
	set     func(*config.ActivityConfig, string)
}
type categoryDef struct {
	title  string
	fields []fieldDef
}
type option struct {
	value string
	label string
}
func tf(id, key string, get func(*config.ActivityConfig) string, set func(*config.ActivityConfig, string)) fieldDef {
	return fieldDef{id: id, label: assets.T(key), kind: kindText, get: get, set: set}
}
func nf(id, key string, get func(*config.ActivityConfig) int, set func(*config.ActivityConfig, int)) fieldDef {
	return fieldDef{
		id:    id,
		label: assets.T(key),
		kind:  kindNumber,
		get:   func(a *config.ActivityConfig) string { return fmt.Sprint(get(a)) },
		set: func(a *config.ActivityConfig, v string) {
			if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n >= 0 {
				set(a, n)
			}
		},
	}
}
func bf(id, key string, get func(*config.ActivityConfig) bool, set func(*config.ActivityConfig, bool)) fieldDef {
	return fieldDef{
		id:    id,
		label: assets.T(key),
		kind:  kindBool,
		get:   func(a *config.ActivityConfig) string { return strconv.FormatBool(get(a)) },
		set:   func(a *config.ActivityConfig, v string) { set(a, v == "true") },
	}
}
func sf(id, key string, opts []option, get func(*config.ActivityConfig) string, set func(*config.ActivityConfig, string)) fieldDef {
	return fieldDef{id: id, label: assets.T(key), kind: kindSelect, options: opts, get: get, set: set}
}
func only(pred func(*config.ActivityConfig) bool, d fieldDef) fieldDef {
	d.visible = pred
	return d
}
func neverVisible(d fieldDef) fieldDef {
	d.visible = func(*config.ActivityConfig) bool { return false }
	return d
}
func typeOptions() []option {
	modes := []struct {
		typ config.ActivityType
		key string
	}{
		{config.TypeGame, "edit.type.playing"},
		{config.TypeListening, "edit.type.listening"},
		{config.TypeWatching, "edit.type.watching"},
		{config.TypeCompeting, "edit.type.competing"},
	}
	opts := make([]option, 0, len(modes))
	for _, m := range modes {
		opts = append(opts, option{value: strconv.Itoa(int(m.typ)), label: assets.T(m.key)})
	}
	return opts
}
func startModeOptions() []option {
	return []option{
		{value: "", label: assets.T("edit.time.off")},
		{value: "elapsed", label: assets.T("edit.time.sinceStart")},
		{value: "custom", label: assets.T("edit.time.duration")},
	}
}
func endModeOptions() []option {
	return []option{
		{value: "", label: assets.T("edit.time.off")},
		{value: "duration", label: assets.T("edit.time.duration")},
	}
}
func parseDuration(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	parts := strings.Split(s, ":")
	var secs int
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || n < 0 {
			return 0, false
		}
		secs = secs*60 + n
	}
	return secs, secs > 0
}
func formatDuration(secs int) string {
	h := secs / 3600
	m := (secs % 3600) / 60
	s := secs % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
func ensureButtons(a *config.ActivityConfig, n int) {
	for len(a.Buttons) < n {
		a.Buttons = append(a.Buttons, config.ButtonConfig{})
	}
}
func buildCategories() []categoryDef {
	return []categoryDef{
		{
			title: assets.T("edit.cat.type"),
			fields: []fieldDef{
				sf("type", "edit.f.type", typeOptions(),
					func(a *config.ActivityConfig) string { return strconv.Itoa(int(a.Type)) },
					func(a *config.ActivityConfig, v string) {
						if i, err := strconv.Atoi(v); err == nil {
							switch config.ActivityType(i) {
							case config.TypeGame, config.TypeListening, config.TypeWatching, config.TypeCompeting:
								a.Type = config.ActivityType(i)
							}
						}
					}),
			},
		},
		{
			title: assets.T("edit.cat.text"),
			fields: []fieldDef{
				tf("name", "edit.f.name", func(a *config.ActivityConfig) string { return a.Name }, func(a *config.ActivityConfig, v string) { a.Name = v }),
				tf("details", "edit.f.details", func(a *config.ActivityConfig) string { return a.Details }, func(a *config.ActivityConfig, v string) { a.Details = v }),
				tf("state", "edit.f.state", func(a *config.ActivityConfig) string { return a.State }, func(a *config.ActivityConfig, v string) { a.State = v }),
			},
		},
		{
			title: assets.T("edit.cat.images"),
			fields: []fieldDef{
				tf("largeImage", "edit.f.largeImage", func(a *config.ActivityConfig) string { return a.LargeImage }, func(a *config.ActivityConfig, v string) { a.LargeImage = v }),
				tf("largeTooltip", "edit.f.largeTooltip", func(a *config.ActivityConfig) string { return a.LargeText }, func(a *config.ActivityConfig, v string) { a.LargeText = v }),
				tf("smallImage", "edit.f.smallImage", func(a *config.ActivityConfig) string { return a.SmallImage }, func(a *config.ActivityConfig, v string) { a.SmallImage = v }),
				tf("smallTooltip", "edit.f.smallTooltip", func(a *config.ActivityConfig) string { return a.SmallText }, func(a *config.ActivityConfig, v string) { a.SmallText = v }),
			},
		},
		{
			title: assets.T("edit.cat.time"),
			fields: []fieldDef{
			sf("start", "edit.f.start", startModeOptions(),
				func(a *config.ActivityConfig) string { return a.Start.Mode },
				func(a *config.ActivityConfig, v string) {
					a.Start.Mode = v
					if v != "custom" {
						a.Start.Value = 0
					}
				}),
			neverVisible(tf("startDuration", "edit.f.duration",
				func(a *config.ActivityConfig) string { return formatDuration(int(a.Start.Value)) },
				func(a *config.ActivityConfig, v string) {
					if secs, ok := parseDuration(v); ok {
						a.Start.Value = int64(secs)
					}
				})),
			sf("end", "edit.f.end", endModeOptions(),
				func(a *config.ActivityConfig) string { return a.End.Mode },
				func(a *config.ActivityConfig, v string) {
					a.End.Mode = v
					if v != "duration" {
						a.End.Value = 0
					}
				}),
			neverVisible(tf("endDuration", "edit.f.duration",
				func(a *config.ActivityConfig) string { return formatDuration(int(a.End.Value)) },
				func(a *config.ActivityConfig, v string) {
					if secs, ok := parseDuration(v); ok {
						a.End.Value = int64(secs)
					}
				})),
			},
		},
		{
			title: assets.T("edit.cat.party"),
			fields: []fieldDef{
				nf("partySize", "edit.f.partySize", func(a *config.ActivityConfig) int { return a.PartySize }, func(a *config.ActivityConfig, n int) { a.PartySize = n }),
				nf("partyMax", "edit.f.partyMax", func(a *config.ActivityConfig) int { return a.PartyMax }, func(a *config.ActivityConfig, n int) { a.PartyMax = n }),
			},
		},
		{
			title: assets.T("edit.cat.buttons"),
			fields: []fieldDef{
				tf("btn1Label", "edit.f.btn1Label", func(a *config.ActivityConfig) string {
					if len(a.Buttons) > 0 {
						return a.Buttons[0].Label
					}
					return ""
				}, func(a *config.ActivityConfig, v string) {
					ensureButtons(a, 1)
					a.Buttons[0].Label = v
				}),
				tf("btn1Url", "edit.f.btn1Url", func(a *config.ActivityConfig) string {
					if len(a.Buttons) > 0 {
						return a.Buttons[0].URL
					}
					return ""
				}, func(a *config.ActivityConfig, v string) {
					ensureButtons(a, 1)
					a.Buttons[0].URL = v
				}),
				tf("btn2Label", "edit.f.btn2Label", func(a *config.ActivityConfig) string {
					if len(a.Buttons) > 1 {
						return a.Buttons[1].Label
					}
					return ""
				}, func(a *config.ActivityConfig, v string) {
					ensureButtons(a, 2)
					a.Buttons[1].Label = v
				}),
				tf("btn2Url", "edit.f.btn2Url", func(a *config.ActivityConfig) string {
					if len(a.Buttons) > 1 {
						return a.Buttons[1].URL
					}
					return ""
				}, func(a *config.ActivityConfig, v string) {
					ensureButtons(a, 2)
					a.Buttons[1].URL = v
				}),
			},
		},
	}
}