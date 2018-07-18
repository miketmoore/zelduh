package entities

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh"
	"github.com/miketmoore/zelduh/config"
	"github.com/miketmoore/zelduh/entityconfig"
	"github.com/miketmoore/zelduh/sprites"
)

// GetPreset gets an entity config preset function by key
func GetPreset(key string) entityConfigPresetFn {
	return entityPresets[key]
}

type entityConfigPresetFn = func(xTiles, yTiles float64) entityconfig.Config

var entityPresets = map[string]entityConfigPresetFn{
	"arrow": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryArrow,
			Movement: &entityconfig.MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"up":    sprites.GetSet("arrowUp"),
				"right": sprites.GetSet("arrowRight"),
				"down":  sprites.GetSet("arrowDown"),
				"left":  sprites.GetSet("arrowLeft"),
			},
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"bomb": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryBomb,
			Movement: &entityconfig.MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("bomb"),
			},
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"coin": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryCoin,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("coin"),
			},
		}
	},
	"explosion": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category:   zelduh.CategoryExplosion,
			Expiration: 12,
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("explosion"),
			},
		}
	},
	"obstacle": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryObstacle,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
		}
	},
	"player": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryPlayer,
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
	"sword": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategorySword,
			Movement: &entityconfig.MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: config.TileSize,
			H: config.TileSize,
			X: config.TileSize * xTiles,
			Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"up":    sprites.GetSet("swordUp"),
				"right": sprites.GetSet("swordRight"),
				"down":  sprites.GetSet("swordDown"),
				"left":  sprites.GetSet("swordLeft"),
			},
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 20,
			},
			Ignore: true,
		}
	},
	"eyeburrower": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryEnemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("eyeburrower"),
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
			Category: zelduh.CategoryHeart,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Hitbox: &entityconfig.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("heart"),
			},
		}

	},
	"skeleton": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryEnemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("skeleton"),
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
			Category: zelduh.CategoryEnemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("skull"),
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
			Category: zelduh.CategoryEnemy,
			W:        config.TileSize, H: config.TileSize, X: config.TileSize * xTiles, Y: config.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("spinner"),
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
			Category: zelduh.CategoryHeart,
			W:        config.TileSize,
			H:        config.TileSize,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			Hitbox: &entityconfig.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("uiCoin"),
			},
		}
	},
	"warpStone": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryWarp,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			W:        config.TileSize,
			H:        config.TileSize,
			Hitbox: &entityconfig.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("warpStone"),
			},
		}
	},
	"puzzleBox": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryMovableObstacle,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			W:        config.TileSize,
			H:        config.TileSize,
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("puzzleBox"),
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
			Category: zelduh.CategoryCollisionSwitch,
			X:        config.TileSize * xTiles,
			Y:        config.TileSize * yTiles,
			W:        config.TileSize,
			H:        config.TileSize,
			Animation: entityconfig.AnimationConfig{
				"default": sprites.GetSet("floorSwitch"),
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
				"default": sprites.GetSet("toggleObstacle"),
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
