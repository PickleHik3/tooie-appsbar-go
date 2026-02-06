package app

import (
	"image"

	"tooie-shelf/internal/config"
	"tooie-shelf/internal/graphics"
	"tooie-shelf/internal/sys"
)

// Model represents the application state.
type Model struct {
	Config      config.Config
	DisplayApps []config.AppConfig // Apps in display order
	TermWidth   int                // Terminal columns
	TermHeight  int                // Terminal rows
	CellPx      sys.CellDim        // Pixel dimensions per cell

	Icons      []image.Image                  // Original high-res images
	SixelCache map[string]graphics.SixelResult // Cached sixel data with dimensions

	ErrorFlash []bool // Per-app error indicator
	Selected   int    // Currently selected app index (-1 for none)

	Ready           bool // Terminal geometry acquired
	NeedsFullRedraw bool // When true, redraw icons; when false, only redraw borders
	SixelsDrawn     bool // True if sixels have been drawn to screen (static mode)
}

// launchResultMsg carries the result of an app launch attempt.
type launchResultMsg struct {
	Index int
	Err   error
}

// NewModel creates a new launcher model.
func NewModel(cfg config.Config) Model {
	displayApps := cfg.GetDisplayApps()
	numApps := len(displayApps)

	return Model{
		Config:          cfg,
		DisplayApps:     displayApps,
		Icons:           make([]image.Image, numApps),
		SixelCache:      make(map[string]graphics.SixelResult),
		ErrorFlash:      make([]bool, numApps),
		Selected:        -1,
		Ready:           false,
		NeedsFullRedraw: true,
		SixelsDrawn:     false,
	}
}

// CacheKey generates a cache key for a sixel render.
func CacheKey(appIndex, widthCells, heightCells int) string {
	return string(rune(appIndex)) + "_" + string(rune(widthCells)) + "_" + string(rune(heightCells))
}

// ClearCache invalidates all cached sixel data.
func (m *Model) ClearCache() {
	m.SixelCache = make(map[string]graphics.SixelResult)
}

// TopRowHeight returns the height of the top row (clock area) in terminal cells.
// This takes all remaining space above the icon grid.
func (m *Model) TopRowHeight() int {
	_, iconGridHeight := m.IconGridDimensions()
	if iconGridHeight >= m.TermHeight-1 {
		return 0
	}
	return m.TermHeight - iconGridHeight
}

// IconGridDimensions calculates the total dimensions of the icon grid area.
// The grid is constrained to the bottom of the terminal with visually square cells.
func (m *Model) IconGridDimensions() (width, height int) {
	if m.Config.Grid.Columns <= 0 || m.Config.Grid.Rows <= 0 {
		return 0, 0
	}
	// Width spans full terminal width
	width = m.TermWidth
	// Height is based on cell height * number of rows
	_, cellHeight := m.GridCellSize()
	height = cellHeight * m.Config.Grid.Rows
	// Ensure we don't exceed terminal height
	if height > m.TermHeight-1 {
		height = m.TermHeight - 1
	}
	return
}

// GridCellSize calculates the size of each grid cell in terminal cells.
// Cells are sized to appear square visually (accounting for terminal cell aspect ratio).
func (m *Model) GridCellSize() (width, height int) {
	if m.Config.Grid.Columns <= 0 || m.Config.Grid.Rows <= 0 {
		return 0, 0
	}
	// Width based on terminal width divided by columns
	width = m.TermWidth / m.Config.Grid.Columns
	// Height: terminal cells are typically ~2x as tall as they are wide,
	// so we use half the width to make cells appear visually square
	height = width / 2
	if height < 1 {
		height = 1
	}
	return
}

// IconCellSize calculates the available space for icons within a cell.
func (m *Model) IconCellSize() (width, height int) {
	cellW, cellH := m.GridCellSize()

	// Subtract padding and borders
	padding := m.Config.Style.Padding
	borderSize := 0
	if m.Config.Style.Border {
		borderSize = 2 // 1 char on each side
	}

	width = cellW - 2*padding - borderSize
	height = cellH - 2*padding - borderSize

	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return
}

// HitTest returns the app index at the given terminal coordinates, or -1 if none.
func (m *Model) HitTest(x, y int) int {
	cellW, cellH := m.GridCellSize()
	if cellW <= 0 || cellH <= 0 {
		return -1
	}

	// Account for top row offset
	topHeight := m.TopRowHeight()
	if y < topHeight {
		return -1 // Click in top row area
	}

	// Adjust y to be relative to the icon grid
	adjustedY := y - topHeight

	col := x / cellW
	row := adjustedY / cellH

	if col < 0 || col >= m.Config.Grid.Columns {
		return -1
	}
	if row < 0 || row >= m.Config.Grid.Rows {
		return -1
	}

	index := row*m.Config.Grid.Columns + col
	if index >= len(m.DisplayApps) {
		return -1
	}

	return index
}

// GetIconScale returns the icon scale for the app at the given display index.
func (m *Model) GetIconScale(index int) float64 {
	if index < 0 || index >= len(m.DisplayApps) {
		return 1.0
	}
	return m.Config.GetIconScale(m.DisplayApps[index])
}

// GetCellWidthForColumn returns the cell width for a specific column.
// Remainder is distributed: first 'remainder' columns get +1 width.
func (m *Model) GetCellWidthForColumn(col int) int {
	cellW, _ := m.GridCellSize()
	totalBaseWidth := cellW * m.Config.Grid.Columns
	remainder := m.TermWidth - totalBaseWidth
	if col < remainder {
		return cellW + 1
	}
	return cellW
}

// GetCellXPosition returns the starting X position (0-indexed) for a given column.
// Accounts for varying column widths when remainder is distributed.
func (m *Model) GetCellXPosition(col int) int {
	cellW, _ := m.GridCellSize()
	totalBaseWidth := cellW * m.Config.Grid.Columns
	remainder := m.TermWidth - totalBaseWidth

	// First 'remainder' columns are 1 wider
	if col <= remainder {
		// All columns up to this one have the extra width
		return col * (cellW + 1)
	}
	// First 'remainder' columns are wide, rest are base width
	return remainder*(cellW+1) + (col-remainder)*cellW
}
