package equipment

import (
	"github.com/faiface/pixel/imdraw"
)

type Sword struct {
	Shape *imdraw.IMDraw
}

func NewSword() Sword {
	return Sword{
		Shape: imdraw.New(nil),
	}
}
