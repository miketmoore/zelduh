package mvmt

import "github.com/faiface/pixel"

func Normalize(target, position pixel.Vec) pixel.Vec {
	return pixel.V(target.X-position.X, target.Y-position.Y)
}
