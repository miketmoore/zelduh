package entities

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh"
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
			Category: zelduh.CategoryArrow,
			Movement: &entityconfig.MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: zelduh.TileSize,
			H: zelduh.TileSize,
			X: zelduh.TileSize * xTiles,
			Y: zelduh.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"up":    zelduh.GetSpriteSet("arrowUp"),
				"right": zelduh.GetSpriteSet("arrowRight"),
				"down":  zelduh.GetSpriteSet("arrowDown"),
				"left":  zelduh.GetSpriteSet("arrowLeft"),
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
			W: zelduh.TileSize,
			H: zelduh.TileSize,
			X: zelduh.TileSize * xTiles,
			Y: zelduh.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("bomb"),
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
			W:        zelduh.TileSize,
			H:        zelduh.TileSize,
			X:        zelduh.TileSize * xTiles,
			Y:        zelduh.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("coin"),
			},
		}
	},
	"explosion": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category:   zelduh.CategoryExplosion,
			Expiration: 12,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("explosion"),
			},
		}
	},
	"obstacle": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryObstacle,
			W:        zelduh.TileSize,
			H:        zelduh.TileSize,
			X:        zelduh.TileSize * xTiles,
			Y:        zelduh.TileSize * yTiles,
		}
	},
	"player": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryPlayer,
			Health:   3,
			W:        zelduh.TileSize,
			H:        zelduh.TileSize,
			X:        zelduh.TileSize * xTiles,
			Y:        zelduh.TileSize * yTiles,
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
				"up":               zelduh.GetSpriteSet("playerUp"),
				"right":            zelduh.GetSpriteSet("playerRight"),
				"down":             zelduh.GetSpriteSet("playerDown"),
				"left":             zelduh.GetSpriteSet("playerLeft"),
				"swordAttackUp":    zelduh.GetSpriteSet("playerSwordUp"),
				"swordAttackRight": zelduh.GetSpriteSet("playerSwordRight"),
				"swordAttackLeft":  zelduh.GetSpriteSet("playerSwordLeft"),
				"swordAttackDown":  zelduh.GetSpriteSet("playerSwordDown"),
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
			W: zelduh.TileSize,
			H: zelduh.TileSize,
			X: zelduh.TileSize * xTiles,
			Y: zelduh.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"up":    zelduh.GetSpriteSet("swordUp"),
				"right": zelduh.GetSpriteSet("swordRight"),
				"down":  zelduh.GetSpriteSet("swordDown"),
				"left":  zelduh.GetSpriteSet("swordLeft"),
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
			W:        zelduh.TileSize, H: zelduh.TileSize, X: zelduh.TileSize * xTiles, Y: zelduh.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("eyeburrower"),
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
			W:        zelduh.TileSize,
			H:        zelduh.TileSize,
			X:        zelduh.TileSize * xTiles,
			Y:        zelduh.TileSize * yTiles,
			Hitbox: &entityconfig.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("heart"),
			},
		}

	},
	"skeleton": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryEnemy,
			W:        zelduh.TileSize, H: zelduh.TileSize, X: zelduh.TileSize * xTiles, Y: zelduh.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("skeleton"),
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
			W:        zelduh.TileSize, H: zelduh.TileSize, X: zelduh.TileSize * xTiles, Y: zelduh.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("skull"),
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
			W:        zelduh.TileSize, H: zelduh.TileSize, X: zelduh.TileSize * xTiles, Y: zelduh.TileSize * yTiles,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("spinner"),
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
			W:        zelduh.TileSize,
			H:        zelduh.TileSize,
			X:        zelduh.TileSize * xTiles,
			Y:        zelduh.TileSize * yTiles,
			Hitbox: &entityconfig.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("uiCoin"),
			},
		}
	},
	"warpStone": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryWarp,
			X:        zelduh.TileSize * xTiles,
			Y:        zelduh.TileSize * yTiles,
			W:        zelduh.TileSize,
			H:        zelduh.TileSize,
			Hitbox: &entityconfig.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("warpStone"),
			},
		}
	},
	"puzzleBox": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryMovableObstacle,
			X:        zelduh.TileSize * xTiles,
			Y:        zelduh.TileSize * yTiles,
			W:        zelduh.TileSize,
			H:        zelduh.TileSize,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("puzzleBox"),
			},
			Movement: &entityconfig.MovementConfig{
				Speed:    1.0,
				MaxMoves: int(zelduh.TileSize) / 2,
				MaxSpeed: 2.0,
			},
		}
	},
	"floorSwitch": func(xTiles, yTiles float64) entityconfig.Config {
		return entityconfig.Config{
			Category: zelduh.CategoryCollisionSwitch,
			X:        zelduh.TileSize * xTiles,
			Y:        zelduh.TileSize * yTiles,
			W:        zelduh.TileSize,
			H:        zelduh.TileSize,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("floorSwitch"),
			},
			Toggleable: true,
		}
	},
	// this is an impassable obstacle that can be toggled "remotely"
	// it has two visual states that coincide with each toggle state
	"toggleObstacle": func(xTiles, yTiles float64) entityconfig.Config {
		// TODO get this working again
		return entityconfig.Config{
			X: zelduh.TileSize * xTiles,
			Y: zelduh.TileSize * yTiles,
			W: zelduh.TileSize,
			H: zelduh.TileSize,
			Animation: entityconfig.AnimationConfig{
				"default": zelduh.GetSpriteSet("toggleObstacle"),
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
