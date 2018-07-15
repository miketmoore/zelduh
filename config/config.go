package config

const FrameRate int = 5

// TileSize defines the width and height of a tile
const TileSize float64 = 48

const (
	WinX float64 = 0
	WinY float64 = 0
	WinW float64 = 800
	WinH float64 = 800
)

const (
	MapW float64 = 672 // 48 * 14
	MapH float64 = 576 // 48 * 12
	MapX         = (WinW - MapW) / 2
	MapY         = (WinH - MapH) / 2
)
