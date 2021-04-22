package zelduh

type WindowConfig struct {
	X, Y, Width, Height float64
}

func NewWindowConfig(x, y, width, height float64) WindowConfig {
	return WindowConfig{x, y, width, height}
}
