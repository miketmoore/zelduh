package entities

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/config"
	"github.com/miketmoore/zelduh/entityconfig"
)

// GetPreset gets an entity config preset function by key
func GetPreset(key string) entityConfigPresetFn {
	return entityPresets[key]
}

type entityConfigPresetFn = func(xTiles, yTiles float64) entityconfig.Config

var entityPresets = map[string]entityConfigPresetFn{
	"arrow": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Arrow,
			Movement: &entityconfig.MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"up":    terraform2d.GetSet("arrowUp"),
				"right": terraform2d.GetSet("arrowRight"),
				"down":  terraform2d.GetSet("arrowDown"),
				"left":  terraform2d.GetSet("arrowLeft"),
			},
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"bomb": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Bomb,
			Movement: &entityconfig.MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("bomb"),
			},
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"coin": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Coin,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("coin"),
			},
		}
	},
	"explosion": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category:   categories.Explosion,
			Expiration: 12,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("explosion"),
			},
		}
	},
	"obstacle": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Obstacle,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
		}
	},
	"player": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Player,
			Health:   3,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Hitbox: &entityconfig.HitboxConfig{
				Box:                  imdraw.New(nil),
				Radius:               15,
				CollisionWithRectMod: 5,
			},
			Movement: &entityconfig.MovementConfig{
				Direction: terraform2d.DirectionDown,
				MaxSpeed:  7.0,
				Speed:     0.0,
			},
			Coins: true,
			Dash: &entityconfig.DashConfig{
				Charge:    0,
				MaxCharge: 50,
				SpeedMod:  7,
			},
			Animation: entityconfig.AnimationConfig{
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
	"sword": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Sword,
			Movement: &entityconfig.MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"up":    terraform2d.GetSet("swordUp"),
				"right": terraform2d.GetSet("swordRight"),
				"down":  terraform2d.GetSet("swordDown"),
				"left":  terraform2d.GetSet("swordLeft"),
			},
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 20,
			},
			Ignore: true,
		}
	},
	"eyeburrower": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Enemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("eyeburrower"),
			},
			Health: 2,
			Hitbox: &entityconfig.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &entityconfig.MovementConfig{
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
	"heart": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Heart,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Hitbox: &entityconfig.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("heart"),
			},
		}

	},
	"skeleton": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Enemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("skeleton"),
			},
			Health: 2,
			Hitbox: &entityconfig.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &entityconfig.MovementConfig{
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
	"skull": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Enemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("skull"),
			},
			Health: 2,
			Hitbox: &entityconfig.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &entityconfig.MovementConfig{
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
	"spinner": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Enemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("spinner"),
			},
			Invincible: true,
			Hitbox: &entityconfig.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &entityconfig.MovementConfig{
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
	"uiCoin": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Heart,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Hitbox: &entityconfig.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("uiCoin"),
			},
		}
	},
	"warpStone": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.Warp,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			W:        config.TileSize,
			H:        config.TileSize,
			Hitbox: &entityconfig.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("warpStone"),
			},
		}
	},
	"puzzleBox": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.MovableObstacle,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			W:        config.TileSize,
			H:        config.TileSize,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("puzzleBox"),
			},
			Movement: &entityconfig.MovementConfig{
				Speed:    1.0,
				MaxMoves: int(config.TileSize) / 2,
				MaxSpeed: 2.0,
			},
		}
	},
	"floorSwitch": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: categories.CollisionSwitch,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			W:        config.TileSize,
			H:        config.TileSize,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("floorSwitch"),
			},
			Toggleable: true,
		}
	},
	// this is an impassable obstacle that can be toggled "remotely"
	// it has two visual states that coincide with each toggle state
	"toggleObstacle": func(xTiles, yTiles float64) entityconfig.Config {
		// TODO get this working again
		return entityconfig.Config{
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			W: config.TileSize,
			H: config.TileSize,
			Animation: entityconfig.AnimationConfig{
				"default": terraform2d.GetSet("toggleObstacle"),
			},
			// Impassable: true,
			Toggleable: true,
		}
	},
}

// WarpStone returns an entity config for a warp stone
func WarpStone(X, Y, WarpToRoomID, HitBoxRadius float64) entityconfig.Config {
	e := entityPresets["warpStone"](X, Y)
	e.WarpToRoomID = 6
	e.Hitbox.Radius = 5
	return e
}
