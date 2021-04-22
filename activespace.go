package zelduh

// ActiveSpaceRectangle defines a rectangle that represents the active space on the screen
type ActiveSpaceRectangle struct {
	X, Y, Width, Height float64
}

func NewActiveSpaceRectangle(x, y, width, height float64) ActiveSpaceRectangle {
	return ActiveSpaceRectangle{x, y, width, height}
}
