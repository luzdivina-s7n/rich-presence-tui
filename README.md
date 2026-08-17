A terminal-based UI for configuring Discord Rich Presence via local IPCs.
- Multiple presets with individual Application IDs
- Full-featured time configuration
- Direct support for asset keys and full image URLs
- System tray support, autostart, and English/Russian UI
- Zero local files created upon download — configuration is stored in your OS config directory and can be deleted directly from the main menu
`bin/rich-presence-tui.exe`
`bin/rich-presence-tui.exe -tray`
- `↑/↓` — Navigate
- `Enter` — Open / Start activity
- `→` — Settings or preset actions
- `Esc` — Back
- `?` — Help
- `Ctrl+Z` — Toggle minimize to tray
- Go (version 1.18 or higher)
- Make (optional, for using the Makefile)
- rsrc (for generating Windows icon resources)
To build the Windows executable:
`bat
build.bat`
Or using `make`:
`make
make build
`
