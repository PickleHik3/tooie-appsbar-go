# tooie-appsbar-go

A terminal-based app bar for Android (Termux) that displays app icons using Sixel graphics.

## Requirements

- Termux with a Sixel-capable terminal
- Go 1.21+

## Build

```bash
go build -ldflags="-s -w" -o tooie-appsbar ./cmd/launcher
```

## Configuration

Config file: `~/.config/tooie-appsbar-go/config.yaml`

```yaml
# Display order - only these apps shown, in this order
display:
  - Chrome
  - Files
  - WhatsApp

grid:
  rows: 2
  columns: 4

style:
  border: true
  padding: 1
  icon_scale: 0.8  # Global icon scale (0.1-1.0)

behavior:
  close_on_launch: true

apps:
  - name: Chrome
    icon: /path/to/chrome.png
    package: com.android.chrome
  - name: Files
    icon: /path/to/files.png
    package: com.google.android.apps.nbu.files
    activity: com.google.android.apps.nbu.files.home.HomeActivity
  - name: WhatsApp
    icon: /path/to/whatsapp.png
    package: com.whatsapp
    icon_scale: 0.6  # Per-app override
```

### Options

| Field | Description |
|-------|-------------|
| `display` | App names in display order (if empty, show all apps) |
| `grid.rows` | Number of rows in the grid |
| `grid.columns` | Number of columns in the grid |
| `style.border` | Show borders around cells |
| `style.padding` | Padding inside cells (in characters) |
| `style.icon_scale` | Global icon scale 0.1-1.0 (default: 1.0) |
| `behavior.close_on_launch` | Exit after launching an app (default: false) |
| `apps[].name` | Display name (used for display order matching) |
| `apps[].icon` | Path to icon image (PNG, JPG, GIF) |
| `apps[].package` | Android package name |
| `apps[].activity` | Optional: specific activity to launch |
| `apps[].icon_scale` | Per-app icon scale override (0.1-1.0) |

## Usage

```bash
./tooie-appsbar
```

- Touch an icon to launch the app
- Press `q` or `Esc` to quit

## Version

0.1

## Roadmap

- [x] Icon scale option (global and per-app)
- [x] Display order configuration
