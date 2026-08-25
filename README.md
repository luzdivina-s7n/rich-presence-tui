A terminal-based UI for configuring Discord Rich Presence via local IPCs.

## Features

- Multiple presets with individual Application IDs
- Full-featured time configuration
- Direct support for asset keys and full image URLs
- System tray support, autostart, and English/Russian UI
- Configuration is stored in your OS config directory and can be deleted directly from the main menu

## Requirements

- Go (version 1.18 or higher)
- Make (optional, for using the Makefile)
- rsrc (for generating Windows icon resources)

## Building 

To build the Windows executable, run the batch script:

```cmd
build.bat
```