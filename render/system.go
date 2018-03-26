package render

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
)

type renderEntity struct {
	ID int
	*components.SpatialComponent
	*components.AppearanceComponent
}

// System is a custom system
type System struct {
	Win               *pixelgl.Window
	playerEntity      renderEntity
	sword             renderEntity
	coins             []renderEntity
	enemies           []renderEntity
	obstacles         []renderEntity
	moveableObstacles []renderEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(appearance *components.AppearanceComponent, spatial *components.SpatialComponent) {
	s.playerEntity = renderEntity{
		AppearanceComponent: appearance,
		SpatialComponent:    spatial,
	}
}

// AddSword adds the sword to the system
func (s *System) AddSword(appearance *components.AppearanceComponent, spatial *components.SpatialComponent) {
	s.sword = renderEntity{
		AppearanceComponent: appearance,
		SpatialComponent:    spatial,
	}
}

// AddObstacle adds an enemy to the system
func (s *System) AddObstacle(id int, appearance *components.AppearanceComponent, spatial *components.SpatialComponent) {
	s.obstacles = append(s.obstacles, renderEntity{
		ID:                  id,
		AppearanceComponent: appearance,
		SpatialComponent:    spatial,
	})
}

// AddMoveableObstacle adds a moveable obstacle to the system
func (s *System) AddMoveableObstacle(id int, appearance *components.AppearanceComponent, spatial *components.SpatialComponent) {
	s.moveableObstacles = append(s.moveableObstacles, renderEntity{
		ID:                  id,
		AppearanceComponent: appearance,
		SpatialComponent:    spatial,
	})
}

// AddCoin adds the player to the system
func (s *System) AddCoin(id int, appearance *components.AppearanceComponent, spatial *components.SpatialComponent) {
	fmt.Printf("render.System.AddCoin() id %d\n", id)
	s.coins = append(s.coins, renderEntity{
		ID:                  id,
		AppearanceComponent: appearance,
		SpatialComponent:    spatial,
	})
}

// AddEnemy adds an enemy to the system
func (s *System) AddEnemy(id int, appearance *components.AppearanceComponent, spatial *components.SpatialComponent) {
	s.enemies = append(s.enemies, renderEntity{
		ID:                  id,
		AppearanceComponent: appearance,
		SpatialComponent:    spatial,
	})
}

// RemoveCoin removes the specified coin from the system
func (s *System) RemoveCoin(id int) {
	for i := len(s.coins) - 1; i >= 0; i-- {
		coin := s.coins[i]
		if coin.ID == id {
			s.coins = append(s.coins[:i], s.coins[i+1:]...)
		}
	}
}

// RemoveEnemy removes the specified enemy from the system
func (s *System) RemoveEnemy(id int) {
	for i := len(s.enemies) - 1; i >= 0; i-- {
		enemy := s.enemies[i]
		if enemy.ID == id {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

// Update changes spatial data based on movement data
func (s *System) Update() {

	sword := s.sword
	sword.Shape.Clear()
	sword.Shape.Color = sword.AppearanceComponent.Color
	sword.Shape.Push(sword.SpatialComponent.Rect.Min)
	sword.Shape.Push(sword.SpatialComponent.Rect.Max)
	sword.Shape.Rectangle(0)
	sword.Shape.Draw(s.Win)

	player := s.playerEntity
	player.Shape.Clear()
	player.Shape.Color = player.AppearanceComponent.Color
	player.Shape.Push(player.SpatialComponent.Rect.Min)
	player.Shape.Push(player.SpatialComponent.Rect.Max)
	player.Shape.Rectangle(0)
	player.Shape.Draw(s.Win)

	for _, enemy := range s.enemies {
		enemy.Shape.Clear()
		enemy.Shape.Color = enemy.AppearanceComponent.Color
		enemy.Shape.Push(enemy.SpatialComponent.Rect.Min)
		enemy.Shape.Push(enemy.SpatialComponent.Rect.Max)
		enemy.Shape.Rectangle(0)
		enemy.Shape.Draw(s.Win)
	}

	for _, coin := range s.coins {
		coin.Shape.Clear()
		coin.Shape.Color = coin.AppearanceComponent.Color
		mod := coin.SpatialComponent.Width / 2
		coin.Shape.Push(coin.SpatialComponent.Rect.Moved(pixel.V(mod, mod)).Min)
		coin.Shape.Circle(24, 0)
		coin.Shape.Draw(s.Win)
	}

	for _, obstacle := range s.obstacles {
		obstacle.Shape.Clear()
		obstacle.Shape.Color = obstacle.AppearanceComponent.Color
		obstacle.Shape.Push(obstacle.SpatialComponent.Rect.Min)
		obstacle.Shape.Push(obstacle.SpatialComponent.Rect.Max)
		obstacle.Shape.Rectangle(0)
		obstacle.Shape.Draw(s.Win)
	}
}
