# Rich Presence TUI

A terminal-based UI for configuring Discord Rich Presence via local IPCs.

## Features

- Multiple presets with individual Application IDs
- Full-featured time configuration
- Direct support for asset keys and full image URLs
- System tray support, autostart, and English/Russian UI
- Zero local files created upon download — configuration is stored in your OS config directory and can be deleted directly from the main menu

## Usage

### Launch Commands

```bash
bin/rich-presence-tui.exe
bin/rich-presence-tui.exe -tray
```

### Keybindings

- ↑/↓ — Navigate
- Enter — Open / Start activity
- → — Settings or preset actions
- Esc — Back
- ? — Help
- Ctrl+Z — Toggle minimize to tray

## Requirements

- Go (version 1.18 or higher)
- Make (optional, for using the Makefile)
- rsrc (for generating Windows icon resources)

## Building from Source

To build the Windows executable, run the batch script:

```cmd
build.bat
```

Or use make:

```bash
make
make build
```
