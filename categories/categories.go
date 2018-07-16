package categories

import (
	"github.com/miketmoore/terraform2d"
)

const (
	Player = terraform2d.EntityCategory(1 << iota)
	Sword
	Arrow
	Bomb
	Enemy
	Explosion
	Heart
	Coin
	Obstacle
	MovableObstacle
	CollisionSwitch
	Warp
)
