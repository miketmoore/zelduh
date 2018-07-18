package zelduh

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/terraform2d"
)

// GetPreset gets an entity config preset function by key
func GetPreset(key string) entityConfigPresetFn {
	return entityPresets[key]
}

type entityConfigPresetFn = func(xTiles, yTiles float64) Config

var entityPresets = map[string]entityConfigPresetFn{
	"arrow": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryArrow,
			Movement: &MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: TileSize,
			H: TileSize,
			X: TileSize * xTiles,
			Y: TileSize * yTiles,
			Animation: AnimationConfig{
				"up":    GetSpriteSet("arrowUp"),
				"right": GetSpriteSet("arrowRight"),
				"down":  GetSpriteSet("arrowDown"),
				"left":  GetSpriteSet("arrowLeft"),
			},
			Hitbox: &HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"bomb": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryBomb,
			Movement: &MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: TileSize,
			H: TileSize,
			X: TileSize * xTiles,
			Y: TileSize * yTiles,
			Animation: AnimationConfig{
				"default": GetSpriteSet("bomb"),
			},
			Hitbox: &HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"coin": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryCoin,
			W:        TileSize,
			H:        TileSize,
			X:        TileSize * xTiles,
			Y:        TileSize * yTiles,
			Animation: AnimationConfig{
				"default": GetSpriteSet("coin"),
			},
		}
	},
	"explosion": func(xTiles, yTiles float64) Config {
		return Config{
			Category:   CategoryExplosion,
			Expiration: 12,
			Animation: AnimationConfig{
				"default": GetSpriteSet("explosion"),
			},
		}
	},
	"obstacle": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryObstacle,
			W:        TileSize,
			H:        TileSize,
			X:        TileSize * xTiles,
			Y:        TileSize * yTiles,
		}
	},
	"player": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryPlayer,
			Health:   3,
			W:        TileSize,
			H:        TileSize,
			X:        TileSize * xTiles,
			Y:        TileSize * yTiles,
			Hitbox: &HitboxConfig{
				Box:                  imdraw.New(nil),
				Radius:               15,
				CollisionWithRectMod: 5,
			},
			Movement: &MovementConfig{
				Direction: terraform2d.DirectionDown,
				MaxSpeed:  7.0,
				Speed:     0.0,
			},
			Coins: true,
			Dash: &DashConfig{
				Charge:    0,
				MaxCharge: 50,
				SpeedMod:  7,
			},
			Animation: AnimationConfig{
				"up":               GetSpriteSet("playerUp"),
				"right":            GetSpriteSet("playerRight"),
				"down":             GetSpriteSet("playerDown"),
				"left":             GetSpriteSet("playerLeft"),
				"swordAttackUp":    GetSpriteSet("playerSwordUp"),
				"swordAttackRight": GetSpriteSet("playerSwordRight"),
				"swordAttackLeft":  GetSpriteSet("playerSwordLeft"),
				"swordAttackDown":  GetSpriteSet("playerSwordDown"),
			},
		}
	},
	"sword": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategorySword,
			Movement: &MovementConfig{
				Direction: terraform2d.DirectionDown,
				Speed:     0.0,
			},
			W: TileSize,
			H: TileSize,
			X: TileSize * xTiles,
			Y: TileSize * yTiles,
			Animation: AnimationConfig{
				"up":    GetSpriteSet("swordUp"),
				"right": GetSpriteSet("swordRight"),
				"down":  GetSpriteSet("swordDown"),
				"left":  GetSpriteSet("swordLeft"),
			},
			Hitbox: &HitboxConfig{
				Radius: 20,
			},
			Ignore: true,
		}
	},
	"eyeburrower": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryEnemy,
			W:        TileSize, H: TileSize, X: TileSize * xTiles, Y: TileSize * yTiles,
			Animation: AnimationConfig{
				"default": GetSpriteSet("eyeburrower"),
			},
			Health: 2,
			Hitbox: &HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &MovementConfig{
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
	"heart": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryHeart,
			W:        TileSize,
			H:        TileSize,
			X:        TileSize * xTiles,
			Y:        TileSize * yTiles,
			Hitbox: &HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: AnimationConfig{
				"default": GetSpriteSet("heart"),
			},
		}

	},
	"skeleton": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryEnemy,
			W:        TileSize, H: TileSize, X: TileSize * xTiles, Y: TileSize * yTiles,
			Animation: AnimationConfig{
				"default": GetSpriteSet("skeleton"),
			},
			Health: 2,
			Hitbox: &HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &MovementConfig{
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
	"skull": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryEnemy,
			W:        TileSize, H: TileSize, X: TileSize * xTiles, Y: TileSize * yTiles,
			Animation: AnimationConfig{
				"default": GetSpriteSet("skull"),
			},
			Health: 2,
			Hitbox: &HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &MovementConfig{
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
	"spinner": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryEnemy,
			W:        TileSize, H: TileSize, X: TileSize * xTiles, Y: TileSize * yTiles,
			Animation: AnimationConfig{
				"default": GetSpriteSet("spinner"),
			},
			Invincible: true,
			Hitbox: &HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &MovementConfig{
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
	"uiCoin": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryHeart,
			W:        TileSize,
			H:        TileSize,
			X:        TileSize * xTiles,
			Y:        TileSize * yTiles,
			Hitbox: &HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: AnimationConfig{
				"default": GetSpriteSet("uiCoin"),
			},
		}
	},
	"warpStone": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryWarp,
			X:        TileSize * xTiles,
			Y:        TileSize * yTiles,
			W:        TileSize,
			H:        TileSize,
			Hitbox: &HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Animation: AnimationConfig{
				"default": GetSpriteSet("warpStone"),
			},
		}
	},
	"puzzleBox": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryMovableObstacle,
			X:        TileSize * xTiles,
			Y:        TileSize * yTiles,
			W:        TileSize,
			H:        TileSize,
			Animation: AnimationConfig{
				"default": GetSpriteSet("puzzleBox"),
			},
			Movement: &MovementConfig{
				Speed:    1.0,
				MaxMoves: int(TileSize) / 2,
				MaxSpeed: 2.0,
			},
		}
	},
	"floorSwitch": func(xTiles, yTiles float64) Config {
		return Config{
			Category: CategoryCollisionSwitch,
			X:        TileSize * xTiles,
			Y:        TileSize * yTiles,
			W:        TileSize,
			H:        TileSize,
			Animation: AnimationConfig{
				"default": GetSpriteSet("floorSwitch"),
			},
			Toggleable: true,
		}
	},
	// this is an impassable obstacle that can be toggled "remotely"
	// it has two visual states that coincide with each toggle state
	"toggleObstacle": func(xTiles, yTiles float64) Config {
		// TODO get this working again
		return Config{
			X: TileSize * xTiles,
			Y: TileSize * yTiles,
			W: TileSize,
			H: TileSize,
			Animation: AnimationConfig{
				"default": GetSpriteSet("toggleObstacle"),
			},
			// Impassable: true,
			Toggleable: true,
		}
	},
}

// WarpStone returns an entity config for a warp stone
func WarpStone(X, Y, WarpToRoomID, HitBoxRadius float64) Config {
	e := entityPresets["warpStone"](X, Y)
	e.WarpToRoomID = 6
	e.Hitbox.Radius = 5
	return e
}
