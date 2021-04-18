package zelduh

import (
	"github.com/faiface/pixel/imdraw"
)

type entityConfigPresetFn = func(xTiles, yTiles float64) EntityConfig

// EntityConfigPresetFnManager contains a map of strings (preset names) to entityConfigPresetFn
// it is used to get an EntityConfig preset
type EntityConfigPresetFnManager struct {
	entityConfigPresetFns map[string]entityConfigPresetFn
}

// NewEntityConfigPresetFnManager returns a new EntityConfigPresetFnManager
func NewEntityConfigPresetFnManager(entityConfigPresetFns map[string]entityConfigPresetFn) EntityConfigPresetFnManager {
	return EntityConfigPresetFnManager{
		entityConfigPresetFns: entityConfigPresetFns,
	}
}

// GetPreset gets an entity config preset function by key
func (m *EntityConfigPresetFnManager) GetPreset(presetName string) entityConfigPresetFn {
	return m.entityConfigPresetFns[presetName]
}

// TODO move this to a higher level configuration location
func BuildEntityConfigPresetFnsMap(tileSize float64) map[string]entityConfigPresetFn {
	return map[string]entityConfigPresetFn{
		"arrow": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryArrow,
				Movement: &MovementConfig{
					Direction: DirectionDown,
					Speed:     0.0,
				},
				W: tileSize,
				H: tileSize,
				X: tileSize * xTiles,
				Y: tileSize * yTiles,
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
		"bomb": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryBomb,
				Movement: &MovementConfig{
					Direction: DirectionDown,
					Speed:     0.0,
				},
				W: tileSize,
				H: tileSize,
				X: tileSize * xTiles,
				Y: tileSize * yTiles,
				Animation: AnimationConfig{
					"default": GetSpriteSet("bomb"),
				},
				Hitbox: &HitboxConfig{
					Radius: 5,
				},
				Ignore: true,
			}
		},
		"coin": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryCoin,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				Animation: AnimationConfig{
					"default": GetSpriteSet("coin"),
				},
			}
		},
		"explosion": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category:   CategoryExplosion,
				Expiration: 12,
				Animation: AnimationConfig{
					"default": GetSpriteSet("explosion"),
				},
			}
		},
		"obstacle": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryObstacle,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
			}
		},
		"player": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryPlayer,
				Health:   3,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				Hitbox: &HitboxConfig{
					Box:                  imdraw.New(nil),
					Radius:               15,
					CollisionWithRectMod: 5,
				},
				Movement: &MovementConfig{
					Direction: DirectionDown,
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
		"sword": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategorySword,
				Movement: &MovementConfig{
					Direction: DirectionDown,
					Speed:     0.0,
				},
				W: tileSize,
				H: tileSize,
				X: tileSize * xTiles,
				Y: tileSize * yTiles,
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
		"eyeburrower": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryEnemy,
				W:        tileSize, H: tileSize, X: tileSize * xTiles, Y: tileSize * yTiles,
				Animation: AnimationConfig{
					"default": GetSpriteSet("eyeburrower"),
				},
				Health: 2,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &MovementConfig{
					Direction:    DirectionDown,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "random",
				},
			}
		},
		"heart": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryHeart,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				Hitbox: &HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: AnimationConfig{
					"default": GetSpriteSet("heart"),
				},
			}

		},
		"skeleton": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryEnemy,
				W:        tileSize, H: tileSize, X: tileSize * xTiles, Y: tileSize * yTiles,
				Animation: AnimationConfig{
					"default": GetSpriteSet("skeleton"),
				},
				Health: 2,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &MovementConfig{
					Direction:    DirectionDown,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "random",
				},
			}
		},
		"skull": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryEnemy,
				W:        tileSize, H: tileSize, X: tileSize * xTiles, Y: tileSize * yTiles,
				Animation: AnimationConfig{
					"default": GetSpriteSet("skull"),
				},
				Health: 2,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &MovementConfig{
					Direction:    DirectionDown,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "random",
				},
			}
		},
		"spinner": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryEnemy,
				W:        tileSize, H: tileSize, X: tileSize * xTiles, Y: tileSize * yTiles,
				Animation: AnimationConfig{
					"default": GetSpriteSet("spinner"),
				},
				Invincible: true,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &MovementConfig{
					Direction:    DirectionRight,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "left-right",
				},
			}
		},
		"uiCoin": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryHeart,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				Hitbox: &HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: AnimationConfig{
					"default": GetSpriteSet("uiCoin"),
				},
			}
		},
		"warpStone": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryWarp,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				W:        tileSize,
				H:        tileSize,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Animation: AnimationConfig{
					"default": GetSpriteSet("warpStone"),
				},
			}
		},
		"puzzleBox": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryMovableObstacle,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				W:        tileSize,
				H:        tileSize,
				Animation: AnimationConfig{
					"default": GetSpriteSet("puzzleBox"),
				},
				Movement: &MovementConfig{
					Speed:    1.0,
					MaxMoves: int(tileSize) / 2,
					MaxSpeed: 2.0,
				},
			}
		},
		"floorSwitch": func(xTiles, yTiles float64) EntityConfig {
			return EntityConfig{
				Category: CategoryCollisionSwitch,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				W:        tileSize,
				H:        tileSize,
				Animation: AnimationConfig{
					"default": GetSpriteSet("floorSwitch"),
				},
				Toggleable: true,
			}
		},
		// this is an impassable obstacle that can be toggled "remotely"
		// it has two visual states that coincide with each toggle state
		"toggleObstacle": func(xTiles, yTiles float64) EntityConfig {
			// TODO get this working again
			return EntityConfig{
				X: tileSize * xTiles,
				Y: tileSize * yTiles,
				W: tileSize,
				H: tileSize,
				Animation: AnimationConfig{
					"default": GetSpriteSet("toggleObstacle"),
				},
				// Impassable: true,
				Toggleable: true,
			}
		},
	}
}
