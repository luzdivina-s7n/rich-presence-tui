package main
import (
	"os"
	"path/filepath"
	"rich-presence-tui/assets"
)
func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	_ = os.MkdirAll(dir, 0o755)
	data := assets.IconICO()
	if data == nil {
		panic("failed to render icon.ico")
	}
	if err := os.WriteFile(filepath.Join(dir, "icon.ico"), data, 0o644); err != nil {
		panic(err)
	}
}