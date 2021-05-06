package zelduh

import "github.com/miketmoore/zelduh/core/entity"

const (
	CategoryPlayer = entity.EntityCategory(1 << iota)
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
	CategoryIgnore
	CategoryRectangle
)
