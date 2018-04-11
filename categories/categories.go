package categories

// Category is used to group entities
type Category uint

const (
	Player = Category(1 << iota)
	Sword
	Arrow
	Bomb
	Enemy
	Explosion
	UI
	Coin
	Obstacle
	MovableObstacle
	CollisionSwitch
	Warp
)
