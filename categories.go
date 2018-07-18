package zelduh

import "github.com/miketmoore/terraform2d"

const (
	CategoryPlayer = terraform2d.EntityCategory(1 << iota)
	CategorySword
	CategoryArrow
	CategoryBomb
	CategoryEnemy
	CategoryExplosion
	CategoryHeart
	CategoryCoin
	CategoryObstacle
	CategoryMovableObstacle
	CategoryCollisionSwitch
	CategoryWarp
)
