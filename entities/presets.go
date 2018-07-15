package entities

import (
	"fmt"

	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/config"
	"github.com/miketmoore/zelduh/direction"
	"github.com/miketmoore/zelduh/rooms"
	"github.com/miketmoore/zelduh/sprites"
)

// GetPreset gets an entity config preset function by key
func GetPreset(key string) entityConfigPresetFn {
	return entityPresets[key]
}

type entityConfigPresetFn = func(xTiles, yTiles float64) rooms.EntityConfig

var entityPresets = map[string]entityConfigPresetFn{
	"arrow": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Arrow,
			Movement: &rooms.MovementConfig{
				Direction: direction.Down,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"up":    sprites.GetSet("arrowUp"),
				"right": sprites.GetSet("arrowRight"),
				"down":  sprites.GetSet("arrowDown"),
				"left":  sprites.GetSet("arrowLeft"),
			},
			Hitbox: &rooms.HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"bomb": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Bomb,
			Movement: &rooms.MovementConfig{
				Direction: direction.Down,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("bomb"),
			},
			Hitbox: &rooms.HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"coin": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Coin,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("coin"),
			},
		}
	},
	"explosion": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category:   categories.Explosion,
			Expiration: 12,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("explosion"),
			},
		}
	},
	"obstacle": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Obstacle,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
		}
	},
	"player": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Player,
			Health:   3,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Hitbox: &rooms.HitboxConfig{
				Box:                  imdraw.New(nil),
				Radius:               15,
				CollisionWithRectMod: 5,
			},
			Movement: &rooms.MovementConfig{
				Direction: direction.Down,
				MaxSpeed:  7.0,
				Speed:     0.0,
			},
			Coins: true,
			Dash: &rooms.DashConfig{
				Charge:    0,
				MaxCharge: 50,
				SpeedMod:  7,
			},
			Animation: rooms.AnimationConfig{
				"up":               sprites.GetSet("playerUp"),
				"right":            sprites.GetSet("playerRight"),
				"down":             sprites.GetSet("playerDown"),
				"left":             sprites.GetSet("playerLeft"),
				"swordAttackUp":    sprites.GetSet("playerSwordUp"),
				"swordAttackRight": sprites.GetSet("playerSwordRight"),
				"swordAttackLeft":  sprites.GetSet("playerSwordLeft"),
				"swordAttackDown":  sprites.GetSet("playerSwordDown"),
			},
		}
	},
	"sword": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Sword,
			Movement: &rooms.MovementConfig{
				Direction: direction.Down,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"up":    sprites.GetSet("swordUp"),
				"right": sprites.GetSet("swordRight"),
				"down":  sprites.GetSet("swordDown"),
				"left":  sprites.GetSet("swordLeft"),
			},
			Hitbox: &rooms.HitboxConfig{
				Radius: 20,
			},
			Ignore: true,
		}
	},
	"eyeburrower": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("eyeburrower"),
			},
			Health: 2,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    direction.Down,
				Speed:        1.0,
				MaxSpeed:     1.0,
				HitSpeed:     10.0,
				HitBackMoves: 10,
				MaxMoves:     100,
				PatternName:  "random",
			},
		}
	},
	"heart": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Heart,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Hitbox: &rooms.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("heart"),
			},
		}

	},
	"skeleton": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("skeleton"),
			},
			Health: 2,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    direction.Down,
				Speed:        1.0,
				MaxSpeed:     1.0,
				HitSpeed:     10.0,
				HitBackMoves: 10,
				MaxMoves:     100,
				PatternName:  "random",
			},
		}
	},
	"skull": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("skull"),
			},
			Health: 2,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    direction.Down,
				Speed:        1.0,
				MaxSpeed:     1.0,
				HitSpeed:     10.0,
				HitBackMoves: 10,
				MaxMoves:     100,
				PatternName:  "random",
			},
		}
	},
	"spinner": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("spinner"),
			},
			Invincible: true,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    direction.Right,
				Speed:        1.0,
				MaxSpeed:     1.0,
				HitSpeed:     10.0,
				HitBackMoves: 10,
				MaxMoves:     100,
				PatternName:  "left-right",
			},
		}
	},
	"uiCoin": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Heart,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Hitbox: &rooms.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("uiCoin"),
			},
		}
	},
	"warpStone": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Warp,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			W:        config.TileSize,
			H:        config.TileSize,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("warpStone"),
			},
		}
	},
	"puzzleBox": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.MovableObstacle,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			W:        config.TileSize,
			H:        config.TileSize,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("puzzleBox"),
			},
			Movement: &rooms.MovementConfig{
				Speed:    1.0,
				MaxMoves: int(config.TileSize) / 2,
				MaxSpeed: 2.0,
			},
		}
	},
	"floorSwitch": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.CollisionSwitch,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			W:        config.TileSize,
			H:        config.TileSize,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("floorSwitch"),
			},
			Toggleable: true,
		}
	},
	// this is an impassable obstacle that can be toggled "remotely"
	// it has two visual states that coincide with each toggle state
	"toggleObstacle": func(xTiles, yTiles float64) rooms.EntityConfig {
		// TODO get this working again
		return rooms.EntityConfig{
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			W: config.TileSize,
			H: config.TileSize,
			Animation: rooms.AnimationConfig{
				"default": sprites.GetSet("toggleObstacle"),
			},
			// Impassable: true,
			Toggleable: true,
		}
	},
}

// WarpStone returns an entity config for a warp stone
func WarpStone(X, Y, WarpToRoomID, HitBoxRadius float64) rooms.EntityConfig {
	fmt.Printf("presetWarpStone\n")
	e := entityPresets["warpStone"](X, Y)
	e.WarpToRoomID = 6
	e.Hitbox.Radius = 5
	return e
}
