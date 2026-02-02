package sys

import (
	"os"

	"golang.org/x/sys/unix"
)

// CellDim represents the pixel dimensions of a terminal cell.
type CellDim struct {
	Width  int
	Height int
}

// TerminalGeometry holds terminal dimensions in both cells and pixels.
type TerminalGeometry struct {
	Cols     int
	Rows     int
	XPixel   int
	YPixel   int
	CellDim  CellDim
}

// GetTerminalGeometry queries the terminal for its dimensions using ioctl.
func GetTerminalGeometry() (TerminalGeometry, error) {
	fd := int(os.Stdout.Fd())
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return TerminalGeometry{}, err
	}

	geom := TerminalGeometry{
		Cols:   int(ws.Col),
		Rows:   int(ws.Row),
		XPixel: int(ws.Xpixel),
		YPixel: int(ws.Ypixel),
	}

	// Calculate cell dimensions
	if geom.XPixel > 0 && geom.Cols > 0 {
		geom.CellDim.Width = geom.XPixel / geom.Cols
	} else {
		geom.CellDim.Width = 10 // Fallback
	}

	if geom.YPixel > 0 && geom.Rows > 0 {
		geom.CellDim.Height = geom.YPixel / geom.Rows
	} else {
		geom.CellDim.Height = 20 // Fallback
	}

	return geom, nil
}
