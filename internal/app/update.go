package app

import (
	"image"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"tooie-appsbar-go/internal/config"
	"tooie-appsbar-go/internal/graphics"
	"tooie-appsbar-go/internal/sys"
)

// clearSelectionMsg clears the selection highlight after a delay.
type clearSelectionMsg struct{ index int }

// clearErrorMsg clears the error flash after a delay.
type clearErrorMsg struct{ index int }

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		queryTerminal,
		loadIcons(m.Config),
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
		m.TermWidth = msg.Width
		m.TermHeight = msg.Height
		m.ClearCache()
		return m, tea.Batch(tea.ClearScreen, queryTerminal)

	case terminalGeometryMsg:
		m.CellPx = msg.CellDim
		m.Ready = true
		m.ClearCache()
		return m, tea.ClearScreen

	case iconsLoadedMsg:
		m.Icons = msg.Icons

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease {
			index := m.HitTest(msg.X, msg.Y)
			if index >= 0 && index < len(m.Config.Apps) {
				m.Selected = index
				return m, tea.Batch(
					launchApp(index, m.Config.Apps[index]),
					clearSelectionAfter(300*time.Millisecond, index),
				)
			}
		}

	case launchResultMsg:
		if msg.Err != nil {
			// Show error flash
			if msg.Index >= 0 && msg.Index < len(m.ErrorFlash) {
				m.ErrorFlash[msg.Index] = true
				return m, clearErrorAfter(500*time.Millisecond, msg.Index)
			}
		} else if m.Config.Behavior.CloseOnLaunch {
			return m, tea.Quit
		}

	case clearSelectionMsg:
		if m.Selected == msg.index {
			m.Selected = -1
		}

	case clearErrorMsg:
		if msg.index >= 0 && msg.index < len(m.ErrorFlash) {
			m.ErrorFlash[msg.index] = false
		}
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

// loadIcons loads all icon images.
func loadIcons(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		icons := make([]image.Image, len(cfg.Apps))
		for i, app := range cfg.Apps {
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

// clearSelectionAfter returns a command that clears selection after a delay.
func clearSelectionAfter(d time.Duration, index int) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return clearSelectionMsg{index}
	})
}

// clearErrorAfter returns a command that clears error flash after a delay.
func clearErrorAfter(d time.Duration, index int) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return clearErrorMsg{index}
	})
}

// launchApp launches an Android app.
func launchApp(index int, app config.AppConfig) tea.Cmd {
	return func() tea.Msg {
		err := sys.LaunchApp(app.Package, app.Activity)
		return launchResultMsg{Index: index, Err: err}
	}
}
