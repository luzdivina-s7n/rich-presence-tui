package tui
import (
	"unicode"
	tea "github.com/charmbracelet/bubbletea"
)
var ruToLat = map[rune]rune{
	'й': 'q', 'ц': 'w', 'у': 'e', 'к': 'r', 'е': 't', 'н': 'y', 'г': 'u',
	'ш': 'i', 'щ': 'o', 'з': 'p', 'х': '[', 'ъ': ']',
	'ф': 'a', 'ы': 's', 'в': 'd', 'а': 'f', 'п': 'g', 'р': 'h', 'о': 'j',
	'л': 'k', 'д': 'l', 'ж': ';', 'э': '\'',
	'я': 'z', 'ч': 'x', 'с': 'c', 'м': 'v', 'и': 'b', 'т': 'n', 'ь': 'm',
	'б': ',', 'ю': '.',
}
func normKey(k tea.KeyMsg) string {
	if k.Type == tea.KeyRunes && len(k.Runes) == 1 {
		r := unicode.ToLower(k.Runes[0])
		if lat, ok := ruToLat[r]; ok {
			return string(lat)
		}
		return string(r)
	}
	return k.String()
}