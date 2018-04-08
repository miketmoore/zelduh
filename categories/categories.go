package categories

// Category is used to group entities
type Category uint

const (
	Player = Category(1 << iota)
	Sword
	Arrow
	Enemy
	Explosion
	Coin
	Obstacle
	MovableObstacle
	CollisionSwitch
	Warp
)
