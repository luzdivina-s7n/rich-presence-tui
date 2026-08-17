package main
import (
	"flag"
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"rich-presence-tui/config"
	"rich-presence-tui/discord"
	"rich-presence-tui/tray"
	"rich-presence-tui/tui"
)
func main() {
	trayFlag := flag.Bool("tray", false, "start minimized to the system tray")
	flag.Parse()
	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	presence := discord.NewPresence()
	defer presence.Close()
	m := tui.New(cfg, presence)
	p := tea.NewProgram(m, tea.WithAltScreen())
	tray.Start(tray.Callbacks{
		OnOpen: func() {
			tray.Restore()
			p.Send(tui.TrayOpenMsg{})
		},
		OnQuit: func() {
			p.Send(tea.Quit)
		},
	})
	if *trayFlag {
		tray.Minimize()
	}
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
