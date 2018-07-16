package entities

import (
	"fmt"

	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/config"
	"github.com/miketmoore/zelduh/rooms"
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
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"up":    terraform2d.GetSet("arrowUp"),
				"right": terraform2d.GetSet("arrowRight"),
				"down":  terraform2d.GetSet("arrowDown"),
				"left":  terraform2d.GetSet("arrowLeft"),
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
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"default": terraform2d.GetSet("bomb"),
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
				"default": terraform2d.GetSet("coin"),
			},
		}
	},
	"explosion": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category:   categories.Explosion,
			Expiration: 12,
			Animation: rooms.AnimationConfig{
				"default": terraform2d.GetSet("explosion"),
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
				Direction: terraform2d.DirectionDown,
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
				"up":               terraform2d.GetSet("playerUp"),
				"right":            terraform2d.GetSet("playerRight"),
				"down":             terraform2d.GetSet("playerDown"),
				"left":             terraform2d.GetSet("playerLeft"),
				"swordAttackUp":    terraform2d.GetSet("playerSwordUp"),
				"swordAttackRight": terraform2d.GetSet("playerSwordRight"),
				"swordAttackLeft":  terraform2d.GetSet("playerSwordLeft"),
				"swordAttackDown":  terraform2d.GetSet("playerSwordDown"),
			},
		}
	},
	"sword": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Sword,
			Movement: &rooms.MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"up":    terraform2d.GetSet("swordUp"),
				"right": terraform2d.GetSet("swordRight"),
				"down":  terraform2d.GetSet("swordDown"),
				"left":  terraform2d.GetSet("swordLeft"),
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
				"default": terraform2d.GetSet("eyeburrower"),
			},
			Health: 2,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    terraform2d.DirectionDown,
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
				"default": terraform2d.GetSet("heart"),
			},
		}

	},
	"skeleton": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: rooms.AnimationConfig{
				"default": terraform2d.GetSet("skeleton"),
			},
			Health: 2,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    terraform2d.DirectionDown,
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
				"default": terraform2d.GetSet("skull"),
			},
			Health: 2,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    terraform2d.DirectionDown,
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
				"default": terraform2d.GetSet("spinner"),
			},
			Invincible: true,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    terraform2d.DirectionRight,
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
				"default": terraform2d.GetSet("uiCoin"),
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
				"default": terraform2d.GetSet("warpStone"),
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
				"default": terraform2d.GetSet("puzzleBox"),
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
				"default": terraform2d.GetSet("floorSwitch"),
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
				"default": terraform2d.GetSet("toggleObstacle"),
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
