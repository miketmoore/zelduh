package render

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/systems"
)

type renderEntity struct {
	ID int
	*components.SpatialComponent
	*systems.AppearanceComponent
}

// System is a custom system
type System struct {
	Win          *pixelgl.Window
	playerEntity renderEntity
	coins        []renderEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(appearance *systems.AppearanceComponent, spatial *components.SpatialComponent) {
	s.playerEntity = renderEntity{
		AppearanceComponent: appearance,
		SpatialComponent:    spatial,
	}
}

// AddCoin adds the player to the system
func (s *System) AddCoin(id int, appearance *systems.AppearanceComponent, spatial *components.SpatialComponent) {
	s.coins = append(s.coins, renderEntity{
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
			i--
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

	// fmt.Printf("render system update length of coins: %d\n", len(s.coins))
	for _, coin := range s.coins {
		coin.Shape.Clear()
		coin.Shape.Color = coin.AppearanceComponent.Color
		coin.Shape.Push(coin.SpatialComponent.Rect.Min)
		coin.Shape.Push(coin.SpatialComponent.Rect.Max)
		coin.Shape.Rectangle(0)
		coin.Shape.Draw(s.Win)
	}
}
