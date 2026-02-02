### Project Structure
```text
sixel-launcher/
├── cmd/
│   └── launcher/
│       └── main.go           # Entry point
├── config/
│   ├── config.yaml           # User definition
│   └── parser.go             # Struct definitions & YAML loader
├── internal/
│   ├── app/
│   │   ├── model.go          # Bubbletea State Machine
│   │   ├── update.go         # Event Loop (Mouse, Keys, Physics)
│   │   └── view.go           # Render Logic (Layout calc)
│   ├── graphics/
│   │   ├── sixel.go          # Image -> Sixel Encoder
│   │   └── resizer.go        # Smart scaling logic (Lanczos)
│   └── sys/
│       ├── exec.go           # Process spawning (setsid)
│       └── terminal.go       # ioctl calls for pixel geometry
└── assets/                   # Default icons (optional)
```

---

### Phase 0: Foundation & Dependencies
**Objective**: Initialize the environment and secure necessary libraries for TUI, math, and imaging.

1.  **Initialize Module**: `go mod init github.com/username/sixel-launcher`
2.  **Install Core Packages**:
    *   `github.com/charmbracelet/bubbletea` (The TUI Runtime)
    *   `github.com/charmbracelet/lipgloss` (Layout & Styling)
    *   `github.com/charmbracelet/harmonica` (Physics/Animation)
    *   `github.com/mattn/go-sixel` (Sixel Encoding)
    *   `github.com/nfnt/resize` (High-quality image resampling - standard `image/draw` is too ugly for icons)
    *   `gopkg.in/yaml.v3` (Config parsing)

---

### Phase 1: The "Responsive" Engine (Math & IO)
**Objective**: Accurately calculate how big an image *should* be in pixels based on how big the terminal is in text cells.

1.  **Terminal Geometry Service** (`internal/sys/terminal.go`):
    *   Implement `GetWindowSize()` using `unix.IoctlGetWinsize`.
    *   **Crucial Step**: Extract `ws_xpixel` and `ws_ypixel` alongside `ws_col` and `ws_row`.
    *   Calculate `CellWidthPx = ws_xpixel / ws_col` and `CellHeightPx = ws_ypixel / ws_row`. *Without this, images will stretch.*

2.  **Configuration Loader** (`config/parser.go`):
    *   Define structs for `Config`, `Style`, `App`.
    *   Load `config.yaml`.
    *   Validate that paths to PNGs exist.

---

### Phase 2: The Graphics Pipeline
**Objective**: Load PNGs and convert them to Sixel strings dynamically.

1.  **Image Loading**:
    *   Load PNG files into `image.Image` structs using standard library.
    *   Keep the *original* high-res `image.Image` in memory. Do not overwrite it. We need the source for resizing later.

2.  **Dynamic Sixel Encoder** (`internal/graphics/sixel.go`):
    *   Create a function `RenderIcon(source image.Image, targetWidthCells int, targetHeightCells int, cellGeometry sys.CellDim) string`.
    *   **Substep A (Calc)**: Convert `targetWidthCells` to target *pixels* using the geometry from Phase 1.
    *   **Substep B (Scale)**: Use `resize.Resize` (Lanczos3) to scale the source image to the target pixel box while **maintaining aspect ratio**.
    *   **Substep C (Encode)**: Encode the result to a Sixel escape sequence string.

---

### Phase 3: The State Machine (Bubbletea + Harmonica)
**Objective**: Set up the application loop and physics state.

1.  **Model Definition** (`internal/app/model.go`):
    ```go
    type Model struct {
        Width, Height int         // Terminal size
        CellPx, CellPy int        // Pixel size of a single char
        Apps          []AppConfig
        Springs       []harmonica.Spring // One physics spring per app
        Positions     []float64          // Current Y-offset for animation
        SixelCache    map[int]string     // Cache key: "index_width_height"
    }
    ```

2.  **The Resize Handler** (`internal/app/update.go`):
    *   Listen for `tea.WindowSizeMsg`.
    *   On trigger:
        1.  Re-query `ioctl` to get new pixel density.
        2.  Calculate layout (1 row, 5 columns).
        3.  Determine the available width per cell (subtracting border padding).
        4.  **Flush the Sixel Cache**.
        5.  Trigger a "Lazy Load" or immediate re-render of icons at the new size.

---

### Phase 4: Physics & Interaction
**Objective**: Make it feel alive.

1.  **Hit Testing Logic**:
    *   The terminal gives us Mouse X/Y in *cells*.
    *   We know the grid layout (e.g., Column 1 is x=0 to x=15).
    *   Map Mouse Click -> App Index.

2.  **Animation Integration**:
    *   Implement `tea.Tick` msg to step the simulation.
    *   Update `harmonica.Spring` logic.
    *   Update `Model.Positions[i]` based on spring velocity.
    *   **Optimization**: Only return a command to re-render if the position changed significantly (`delta > 0.1`).

3.  **Process Spawning** (`internal/sys/exec.go`):
    *   Function `Launch(cmd string, args []string)`.
    *   Use `syscall.Setsid` to detach.
    *   Use `exec.Command` with `Stdin/out/err` set to `nil` to prevent it from hijacking the TUI.

---

### Phase 5: The View Layer (Layout)
**Objective**: Combine Lipgloss borders with Sixel payloads.

1.  **Grid Layout**:
    *   Use `lipgloss.JoinHorizontal`.
    *   Calculate dynamic widths: `ColumnWidth = TerminalWidth / 5`.

2.  **The Cell Renderer**:
    *   **Layer 1 (Container)**: A Lipgloss box with Borders.
    *   **Layer 2 (Padding)**: Apply `PaddingTop` based on the `harmonica` calculated position. This creates the "bounce" effect.
    *   **Layer 3 (Content)**: The Sixel string.
    *   *Warning*: Lipgloss measures string length. Sixel strings have length 0 or varying length depending on implementation. We must force the Lipgloss box to have a fixed Width/Height so it doesn't collapse around the "invisible" Sixel string.

### Phase 6: Compilation & Verification
1.  **Build**: `go build -ldflags="-s -w" -o launcher cmd/launcher/main.go`
2.  **Test Scenarios**:
    *   Maximize terminal -> Icons should grow and stay sharp.
    *   Shrink to 76x36 -> Icons should shrink, borders tighten.
    *   Click -> Icon bounces down, app opens.

This plan ensures high-performance rendering by caching Sixels and only regenerating them when the terminal physically changes size, while handling animations via lightweight layout offsets.

**Shall I begin Phase 0 and 1 (Project setup, Config, and Terminal Geometry)?**
