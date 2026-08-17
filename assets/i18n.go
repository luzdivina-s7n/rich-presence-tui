package assets
type dict struct {
	Code string
	Name string
	M    map[string]string
}
var current *dict
var languages = []*dict{ru, en}
var ru = &dict{Code: "ru", Name: "Русский", M: map[string]string{
	"app.title": "Rich Presence TUI",
	"menu.presets": "Пресеты",
	"menu.section": "Меню",
	"menu.settings": "Настройки",
	"menu.reset":    "Удалить конфигурацию",
	"menu.resetConfirm": "Удалить конфигурацию (пресеты и настройки)?",
	"menu.resetDone":    "Конфигурация удалена",
	"menu.exit":     "Выход",
	"menu.nopreset": "Сначала выберите пресет",
	"menu.hint":     "[↑/↓] выбор   [Enter] старт/стоп   [→] настройки   [?] помощь",
	"appid.title":   "Выбор Application ID",
	"appid.enter":   "Ввести свой",
	"appid.prompt":  "Application ID",
	"appid.hint":    "[↑/↓] выбор   [Enter] ок   [D] удалить   [Esc] назад",
	"appid.invalid": "Application ID должен состоять только из цифр",
	"presets.create":  "Создать новый пресет",
	"presets.name":    "Название пресета",
	"presets.exists":  "уже существует",
	"presets.empty":   "Пресетов пока нет",
	"presets.full":    "Достигнут лимит — максимум 10 пресетов",
	"presets.created": "Пресет создан",
	"pact.title":        "Пресет: %s",
	"pact.editActivity": "Редактировать активность",
	"pact.changeAppID":  "Изменить Application ID",
	"pact.start":        "Запустить Rich Presence",
	"pact.restart":      "Перезапустить Rich Presence",
	"pact.stop":         "Выключить Rich Presence",
	"pact.tray":         "Свернуть в трей",
	"pact.delete":       "Удалить пресет",
	"pact.back":         "Назад",
	"pact.confirm":      "Удалить пресет «%s»?",
	"pact.deleted":      "Пресет удалён",
	"pact.active":       "активен",
	"pact.inactive":     "выключен",
	"pact.hint":         "[↑/↓] выбор   [Enter] выполнить   [Esc] назад",
	"edit.title":         "Редактирование: %s",
	"edit.cat.type":      "Тип активности",
	"edit.cat.text":      "Текст",
	"edit.cat.images":    "Изображения",
	"edit.cat.time":      "Время",
	"edit.cat.party":     "Группа",
	"edit.cat.buttons":   "Кнопки",
	"edit.f.type":         "Тип активности",
	"edit.f.name":         "Название приложения/игры",
	"edit.f.details":      "Details",
	"edit.f.state":        "State",
	"edit.f.largeImage":   "Large Image",
	"edit.f.largeTooltip": "Large Tooltip",
	"edit.f.smallImage":   "Small Image",
	"edit.f.smallTooltip": "Small Tooltip",
	"edit.f.start":        "Начало",
	"edit.f.end":          "Конец",
	"edit.f.duration":     "Часы:Минуты:Секунды",
	"edit.f.partySize":    "Участники",
	"edit.f.partyMax":     "Максимум участников",
	"edit.f.btn1Label":    "Кнопка 1 — подпись",
	"edit.f.btn1Url":      "Кнопка 1 — URL",
	"edit.f.btn2Label":    "Кнопка 2 — подпись",
	"edit.f.btn2Url":      "Кнопка 2 — URL",
	"edit.type.playing":   "Играет",
	"edit.type.listening": "Слушает",
	"edit.type.watching":  "Смотрит",
	"edit.type.competing": "Соревнуется",
	"edit.time.off":        "Выключено",
	"edit.time.sinceStart": "С момента запуска",
	"edit.time.duration":   "Ввести свое",
	"edit.back":            "Назад",
	"edit.panelHint":       "Значения",
	"edit.hint.fields":     "[↑/↓] поле   [→/Enter] изменить   [←/Esc] назад",
	"edit.hint.select":     "[↑/↓] выбор   [→] ок   [←] отмена",
	"edit.hint.space":      "[Пробел] переключить",
	"edit.hint.enter":      "[Enter] изменить значение",
	"edit.hint.edit":       "[Enter] применить   [Esc] отмена",
	"edit.hint.back":       "[Enter] назад",
	"settings.title":         "Настройки",
	"settings.language":      "Язык",
	"settings.autostart":     "Автозапуск при входе в систему",
	"settings.minimizeTray":  "Сворачивать в трей (Ctrl+Z)",
	"settings.confirmDelete": "Подтверждать удаление пресета",
	"settings.defaultAppID":  "Application ID по умолчанию",
	"settings.configPath":    "Файл конфигурации",
	"settings.back":          "Назад",
	"settings.hint":          "[↑/↓] выбор   [Enter] открыть   [Пробел] вкл/выкл   [Esc] назад",
	"settings.unsupported":   "недоступно на этой платформе",
	"help.title":            "Помощь",
	"help.hint":             "[←/Esc] назад",
	"help.section.visibility": "Почему другие не видят Rich Presence?",
	"help.keys": "[Enter] — старт/стоп   [→] — настройки пресета\n" +
		"[E] — редактор   [A] — действия   [N] — новый пресет\n" +
		"[L] — язык   [T] — трей   [?] — помощь   [Q] — выход\n" +
		"[Esc]/[←] — назад   [↑/↓] — наведение",
	"help.visibility": "1. «Статус активности» — отдельный тумблер в Конфиденциальности,\n" +
		"   если он выключен, активность видите только вы.\n" +
		"2. Rich Presence виден ТОЛЬКО друзьям, незнакомцам — никогда.\n" +
		"3. У зрителя: Настройки → Активность игры → скрытые приложения.\n" +
		"4. Режим стримера на вашем аккаунте скрывает активность от всех.\n" +
		"5. У зрителя может быть выключен показ активности других.",
	"common.yes":     "Да",
	"common.no":      "Нет",
	"common.on":      "Вкл",
	"common.off":     "Выкл",
	"common.back":    "Назад",
	"status.saved":      "Сохранено",
	"status.started":    "Rich Presence запущен:",
	"status.stopped":    "Rich Presence выключен",
	"status.trayed":     "Свёрнуто в трей",
	"status.lang":       "Язык: %s",
	"status.notFound":   "Пресет не найден",
	"status.discordOff": "Discord не запущен или не поддерживает IPC",
	"status.running":    "запущен",
	"status.idle":       "остановлен",
	"tray.open": "Открыть",
	"tray.quit": "Выход",
}}
var en = &dict{Code: "en", Name: "English", M: map[string]string{
	"app.title": "Rich Presence TUI",
	"menu.presets":  "Presets",
	"menu.section":  "Menu",
	"menu.settings": "Settings",
	"menu.reset":    "Delete configuration",
	"menu.resetConfirm": "Delete the configuration (presets and settings)?",
	"menu.resetDone":    "Configuration deleted",
	"menu.exit":     "Exit",
	"menu.nopreset": "Select a preset first",
	"menu.hint":     "[up/down] select   [Enter] start/stop   [right] open   [?] help",
	"appid.title":   "Select Application ID",
	"appid.enter":   "Enter custom",
	"appid.prompt":  "Application ID",
	"appid.hint":    "[up/down] select   [Enter] ok   [D] delete   [Esc] back",
	"appid.invalid": "Application ID must contain digits only",
	"presets.create":  "Create new preset",
	"presets.name":    "Preset name",
	"presets.exists":  "already exists",
	"presets.empty":   "No presets yet",
	"presets.full":    "Limit reached — 10 presets max",
	"presets.created": "Preset created",
	"pact.title":        "Preset: %s",
	"pact.editActivity": "Edit activity",
	"pact.changeAppID":  "Change Application ID",
	"pact.start":        "Start Rich Presence",
	"pact.restart":      "Restart Rich Presence",
	"pact.stop":         "Turn off Rich Presence",
	"pact.tray":         "Minimize to tray",
	"pact.delete":       "Delete preset",
	"pact.back":         "Back",
	"pact.confirm":      "Delete preset \"%s\"?",
	"pact.deleted":      "Preset deleted",
	"pact.active":       "active",
	"pact.inactive":     "off",
	"pact.hint":         "[up/down] select   [Enter] run   [Esc] back",
	"edit.title":         "Editing: %s",
	"edit.cat.type":      "Activity type",
	"edit.cat.text":      "Text",
	"edit.cat.images":    "Images",
	"edit.cat.time":      "Time",
	"edit.cat.party":     "Party",
	"edit.cat.buttons":   "Buttons",
	"edit.f.type":         "Activity type",
	"edit.f.name":         "Application/game name",
	"edit.f.details":      "Details",
	"edit.f.state":        "State",
	"edit.f.largeImage":   "Large Image",
	"edit.f.largeTooltip": "Large Tooltip",
	"edit.f.smallImage":   "Small Image",
	"edit.f.smallTooltip": "Small Tooltip",
	"edit.f.start":        "Start",
	"edit.f.end":          "End",
	"edit.f.duration":     "Hours:Minutes:Seconds",
	"edit.f.partySize":    "Party Size",
	"edit.f.partyMax":     "Party Max",
	"edit.f.btn1Label":    "Button 1 — label",
	"edit.f.btn1Url":      "Button 1 — URL",
	"edit.f.btn2Label":    "Button 2 — label",
	"edit.f.btn2Url":      "Button 2 — URL",
	"edit.type.playing":   "Playing",
	"edit.type.listening": "Listening",
	"edit.type.watching":  "Watching",
	"edit.type.competing": "Competing",
	"edit.time.off":        "Off",
	"edit.time.sinceStart": "Since start",
	"edit.time.duration":   "Enter your own",
	"edit.back":            "Back",
	"edit.panelHint":       "Values",
	"edit.hint.fields":     "[up/down] field   [right/Enter] edit   [left/Esc] back",
	"edit.hint.select":     "[up/down] pick   [right] ok   [left] cancel",
	"edit.hint.space":      "[Space] toggle",
	"edit.hint.enter":      "[Enter] edit value",
	"edit.hint.edit":       "[Enter] apply   [Esc] cancel",
	"edit.hint.back":       "[Enter] back",
	"settings.title":         "Settings",
	"settings.language":      "Language",
	"settings.autostart":     "Run at OS login",
	"settings.minimizeTray":  "Minimize to tray (Ctrl+Z)",
	"settings.confirmDelete": "Confirm before deleting a preset",
	"settings.defaultAppID":  "Default Application ID",
	"settings.configPath":    "Configuration file",
	"settings.back":          "Back",
	"settings.hint":          "[up/down] select   [Enter] open   [Space] on/off   [Esc] back",
	"settings.unsupported":   "unsupported on this platform",
	"help.title":              "Help",
	"help.hint":               "[left/Esc] back",
	"help.section.visibility": "Why can't others see Rich Presence?",
	"help.keys": "[Enter] — start/stop   [right] — preset settings\n" +
		"[E] — edit   [A] — actions   [N] — new preset\n" +
		"[L] — language   [T] — tray   [?] — help   [Q] — quit\n" +
		"[Esc]/[left] — back   [up/down] — move",
	"help.visibility": "1. \"Activity Status\" under Privacy & Safety is a separate toggle\n" +
		"   from your public profile — when off, only you see your activity.\n" +
		"2. Rich Presence is shown only to FRIENDS — never to strangers.\n" +
		"3. Viewer side: Settings → Game Activity → hidden applications.\n" +
		"4. Streamer Mode on your account hides activity from everyone.\n" +
		"5. The viewer may have disabled showing others' activity.",
	"common.yes":     "Yes",
	"common.no":      "No",
	"common.on":      "On",
	"common.off":     "Off",
	"common.back":    "Back",
	"status.saved":      "Saved",
	"status.started":    "Rich Presence started:",
	"status.stopped":    "Rich Presence turned off",
	"status.trayed":     "Minimized to tray",
	"status.lang":       "Language: %s",
	"status.notFound":   "Preset not found",
	"status.discordOff": "Discord is not running or IPC is unavailable",
	"status.running":    "running",
	"status.idle":       "idle",
	"tray.open": "Open",
	"tray.quit": "Quit",
}}
func Languages() []string {
	out := make([]string, 0, len(languages))
	for _, l := range languages {
		out = append(out, l.Code)
	}
	return out
}
func LanguageName(code string) string {
	for _, l := range languages {
		if l.Code == code {
			return l.Name
		}
	}
	return code
}
func SetLanguage(code string) {
	for _, l := range languages {
		if l.Code == code {
			current = l
			return
		}
	}
	if current == nil {
		current = ru
	}
}
func T(key string) string {
	if current == nil {
		current = ru
	}
	if s, ok := current.M[key]; ok {
		return s
	}
	return key
}
