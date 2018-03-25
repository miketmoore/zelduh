package render

import (
	"fmt"

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
	Win          *pixelgl.Window
	playerEntity renderEntity
	coins        []renderEntity
	enemies      []renderEntity
	obstacles    []renderEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(appearance *components.AppearanceComponent, spatial *components.SpatialComponent) {
	s.playerEntity = renderEntity{
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

// Update changes spatial data based on movement data
func (s *System) Update() {
	player := s.playerEntity
	player.Shape.Clear()
	player.Shape.Color = player.AppearanceComponent.Color
	player.Shape.Push(player.SpatialComponent.Rect.Min)
	player.Shape.Push(player.SpatialComponent.Rect.Max)
	player.Shape.Rectangle(0)
	player.Shape.Draw(s.Win)

	for _, coin := range s.coins {
		coin.Shape.Clear()
		coin.Shape.Color = coin.AppearanceComponent.Color
		coin.Shape.Push(coin.SpatialComponent.Rect.Min)
		coin.Shape.Push(coin.SpatialComponent.Rect.Max)
		coin.Shape.Rectangle(0)
		coin.Shape.Draw(s.Win)
	}

	for _, enemy := range s.enemies {
		enemy.Shape.Clear()
		enemy.Shape.Color = enemy.AppearanceComponent.Color
		enemy.Shape.Push(enemy.SpatialComponent.Rect.Min)
		enemy.Shape.Push(enemy.SpatialComponent.Rect.Max)
		enemy.Shape.Rectangle(0)
		enemy.Shape.Draw(s.Win)
	}

	for _, obstacle := range s.obstacles {
		obstacle.Shape.Clear()
		obstacle.Shape.Color = obstacle.AppearanceComponent.Color
		obstacle.Shape.Push(obstacle.SpatialComponent.Rect.Min)
		obstacle.Shape.Push(obstacle.SpatialComponent.Rect.Max)
		obstacle.Shape.Rectangle(1)
		obstacle.Shape.Draw(s.Win)
	}
}
