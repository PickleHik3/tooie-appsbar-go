package app

import (
	"fmt"
	"image"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"tooie-appsbar-go/internal/config"
	"tooie-appsbar-go/internal/graphics"
	"tooie-appsbar-go/internal/sys"
)

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		queryTerminal,
		loadIcons(m.DisplayApps),
	)
}

// Update handles events and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// Only set dimensions on first receive, ignore resizes (e.g., soft keyboard)
		if m.TermWidth == 0 && m.TermHeight == 0 {
			m.TermWidth = msg.Width
			m.TermHeight = msg.Height
			return m, queryTerminal
		}
		// Ignore subsequent resize events to prevent redraws
		return m, nil

	case terminalGeometryMsg:
		m.CellPx = msg.CellDim
		m.Ready = true
		m.ClearCache()
		return m, tea.ClearScreen

	case iconsLoadedMsg:
		m.Icons = msg.Icons

	case tea.MouseMsg:
		// Only handle release events, ignore press/motion to avoid extra redraws
		if msg.Action != tea.MouseActionRelease {
			return m, nil
		}
		index := m.HitTest(msg.X, msg.Y)
		if index >= 0 && index < len(m.DisplayApps) {
			// Flash visual feedback directly via ANSI (no View() redraw)
			m.flashCell(index)

			// Launch app
			go sys.LaunchApp(m.DisplayApps[index].Package, m.DisplayApps[index].Activity)

			if m.Config.Behavior.CloseOnLaunch {
				return m, tea.Quit
			}
		}
		return m, nil
	}

	return m, nil
}

// terminalGeometryMsg carries terminal pixel dimensions.
type terminalGeometryMsg struct {
	CellDim sys.CellDim
}

// iconsLoadedMsg carries loaded icon images.
type iconsLoadedMsg struct {
	Icons []image.Image
}

// queryTerminal queries terminal geometry.
func queryTerminal() tea.Msg {
	geom, err := sys.GetTerminalGeometry()
	if err != nil {
		// Use fallback dimensions
		return terminalGeometryMsg{
			CellDim: sys.CellDim{Width: 10, Height: 20},
		}
	}
	return terminalGeometryMsg{CellDim: geom.CellDim}
}

// loadIcons loads all icon images for the display apps.
func loadIcons(apps []config.AppConfig) tea.Cmd {
	return func() tea.Msg {
		icons := make([]image.Image, len(apps))
		for i, app := range apps {
			if app.Icon != "" {
				img, err := graphics.LoadImage(app.Icon)
				if err == nil {
					icons[i] = img
				} else {
					icons[i] = graphics.CreatePlaceholder(64, 64)
				}
			} else {
				icons[i] = graphics.CreatePlaceholder(64, 64)
			}
		}
		return iconsLoadedMsg{Icons: icons}
	}
}

// flashCell provides visual feedback by briefly highlighting the cell border.
// Uses direct ANSI output to avoid triggering a full View() redraw.
func (m *Model) flashCell(index int) {
	if !m.Config.Style.Border {
		return
	}

	cellW, cellH := m.GridCellSize()
	if cellW <= 0 || cellH <= 0 {
		return
	}

	col := index % m.Config.Grid.Columns
	row := index / m.Config.Grid.Columns

	// Calculate top-left position of the cell (1-indexed for ANSI)
	startX := col*cellW + 1
	startY := row*cellH + 1

	// Highlight color (bright cyan)
	highlight := "\x1b[96m" // Bright cyan
	reset := "\x1b[0m"

	// Rounded border characters
	topLeft := "╭"
	topRight := "╮"
	bottomLeft := "╰"
	bottomRight := "╯"
	horizontal := "─"
	vertical := "│"

	var output string

	// Top border
	output += fmt.Sprintf("\x1b[%d;%dH%s%s", startY, startX, highlight, topLeft)
	for x := 1; x < cellW-1; x++ {
		output += horizontal
	}
	output += topRight

	// Side borders
	for y := 1; y < cellH-1; y++ {
		output += fmt.Sprintf("\x1b[%d;%dH%s", startY+y, startX, vertical)
		output += fmt.Sprintf("\x1b[%d;%dH%s", startY+y, startX+cellW-1, vertical)
	}

	// Bottom border
	output += fmt.Sprintf("\x1b[%d;%dH%s", startY+cellH-1, startX, bottomLeft)
	for x := 1; x < cellW-1; x++ {
		output += horizontal
	}
	output += bottomRight + reset

	// Move cursor to bottom
	output += fmt.Sprintf("\x1b[%d;1H", m.TermHeight)

	// Write directly to stdout
	fmt.Fprint(os.Stdout, output)
}
