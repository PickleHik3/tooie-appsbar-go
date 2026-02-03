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
grid:
  rows: 2
  columns: 4

style:
  border: true
  padding: 1

apps:
  - name: Chrome
    icon: /path/to/chrome.png
    package: com.android.chrome
  - name: Files
    icon: /path/to/files.png
    package: com.google.android.apps.nbu.files
    activity: com.google.android.apps.nbu.files.home.HomeActivity
```

### Options

| Field | Description |
|-------|-------------|
| `grid.rows` | Number of rows in the grid |
| `grid.columns` | Number of columns in the grid |
| `style.border` | Show borders around cells |
| `style.padding` | Padding inside cells (in characters) |
| `apps[].name` | Display name (unused currently) |
| `apps[].icon` | Path to icon image (PNG, JPG, GIF) |
| `apps[].package` | Android package name |
| `apps[].activity` | Optional: specific activity to launch |
| `behavior.close_on_launch` | Exit after launching an app (default: false) |

## Usage

```bash
./tooie-appsbar
```

- Touch an icon to launch the app
- Press `q` or `Esc` to quit

## Version

0.1

## Roadmap

- [ ] Fixed icon size option (override auto-scaling)
