package main

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh"
)

const (
	PresetNameArrow            zelduh.PresetName = "arrow"
	PresetNameBomb             zelduh.PresetName = "bomb"
	PresetNameCoin             zelduh.PresetName = "coin"
	PresetNameExplosion        zelduh.PresetName = "explosion"
	PresetNameObstacle         zelduh.PresetName = "obstacle"
	PresetNamePlayer           zelduh.PresetName = "player"
	PresetNameFloorSwitch      zelduh.PresetName = "floorSwitch"
	PresetNameToggleObstacle   zelduh.PresetName = "toggleObstacle"
	PresetNamePuzzleBox        zelduh.PresetName = "puzzleBox"
	PresetNameWarpStone        zelduh.PresetName = "warpStone"
	PresetNameUICoin           zelduh.PresetName = "uiCoin"
	PresetNameEnemySpinner     zelduh.PresetName = "spinner"
	PresetNameEnemySkull       zelduh.PresetName = "skull"
	PresetNameEnemySkeleton    zelduh.PresetName = "skeleton"
	PresetNameHeart            zelduh.PresetName = "heart"
	PresetNameEnemyEyeBurrower zelduh.PresetName = "eyeBurrower"
	PresetNameSword            zelduh.PresetName = "sword"
	PresetNameDialogCorner     zelduh.PresetName = "dialogCorner"
	PresetNameDialogSide       zelduh.PresetName = "dialogSide"
)

type BuildWarpFn func(
	warpToRoomID zelduh.RoomID,
	coordinates zelduh.Coordinates,
	hitboxRadius float64,
) zelduh.EntityConfig

func buildWarpFnFactory(
	tileSize float64,
	dimensions zelduh.Dimensions,
) BuildWarpFn {

	return func(
		warpToRoomID zelduh.RoomID,
		coordinates zelduh.Coordinates,
		hitboxRadius float64,
	) zelduh.EntityConfig {
		return zelduh.EntityConfig{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: warpToRoomID,
			Dimensions:   dimensions,
			Coordinates:  coordinates,
			Hitbox: &zelduh.HitboxConfig{
				Radius: hitboxRadius,
			},
		}
	}
}

type BuildWarpStoneFn func(
	WarpToRoomID zelduh.RoomID,
	coordinates zelduh.Coordinates,
	HitBoxRadius float64,
) zelduh.EntityConfig

func buildWarpStoneFnFactory(
	entityConfigPresetFnManager *zelduh.EntityConfigPresetFnManager,
) BuildWarpStoneFn {
	return func(
		WarpToRoomID zelduh.RoomID,
		coordinates zelduh.Coordinates,
		HitBoxRadius float64,
	) zelduh.EntityConfig {
		presetFn := entityConfigPresetFnManager.GetPreset("warpStone")
		e := presetFn(zelduh.Coordinates{X: coordinates.X, Y: coordinates.Y})
		e.WarpToRoomID = 6
		e.Hitbox.Radius = 5
		return e
	}
}

// TODO move this to a higher level configuration location
func BuildEntityConfigPresetFnsMap(tileSize float64) map[zelduh.PresetName]zelduh.EntityConfigPresetFn {

	dimensions := zelduh.Dimensions{
		Width:  tileSize,
		Height: tileSize,
	}

	buildCoordinates := func(coordinates zelduh.Coordinates) zelduh.Coordinates {
		return zelduh.Coordinates{
			X: tileSize * coordinates.X,
			Y: tileSize * coordinates.Y,
		}
	}

	return map[zelduh.PresetName]zelduh.EntityConfigPresetFn{
		PresetNameArrow: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryArrow,
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					Speed:     0.0,
				},
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"up":    zelduh.GetSpriteSet("arrowUp"),
					"right": zelduh.GetSpriteSet("arrowRight"),
					"down":  zelduh.GetSpriteSet("arrowDown"),
					"left":  zelduh.GetSpriteSet("arrowLeft"),
				},
				Hitbox: &zelduh.HitboxConfig{
					Radius: 5,
				},
				Ignore: true,
			}
		},
		PresetNameBomb: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryBomb,
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					Speed:     0.0,
				},
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("bomb"),
				},
				Hitbox: &zelduh.HitboxConfig{
					Radius: 5,
				},
				Ignore: true,
			}
		},
		PresetNameCoin: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryCoin,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("coin"),
				},
			}
		},
		PresetNameExplosion: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:   zelduh.CategoryExplosion,
				Expiration: 12,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("explosion"),
				},
			}
		},
		PresetNameObstacle: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryObstacle,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
			}
		},
		PresetNamePlayer: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryPlayer,
				Health:      3,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Hitbox: &zelduh.HitboxConfig{
					Box:                  imdraw.New(nil),
					Radius:               15,
					CollisionWithRectMod: 5,
				},
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					MaxSpeed:  7.0,
					Speed:     0.0,
				},
				Coins: true,
				Dash: &zelduh.DashConfig{
					Charge:    0,
					MaxCharge: 50,
					SpeedMod:  7,
				},
				Animation: zelduh.AnimationConfig{
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
		PresetNameSword: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategorySword,
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					Speed:     0.0,
				},
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"up":    zelduh.GetSpriteSet("swordUp"),
					"right": zelduh.GetSpriteSet("swordRight"),
					"down":  zelduh.GetSpriteSet("swordDown"),
					"left":  zelduh.GetSpriteSet("swordLeft"),
				},
				Hitbox: &zelduh.HitboxConfig{
					Radius: 20,
				},
				Ignore: true,
			}
		},
		PresetNameEnemyEyeBurrower: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("eyeburrower"),
				},
				Health: 2,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &zelduh.MovementConfig{
					Direction:    zelduh.DirectionDown,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "random",
				},
			}
		},
		PresetNameHeart: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryHeart,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Hitbox: &zelduh.HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("heart"),
				},
			}

		},
		PresetNameEnemySkeleton: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("skeleton"),
				},
				Health: 2,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &zelduh.MovementConfig{
					Direction:    zelduh.DirectionDown,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "random",
				},
			}
		},
		PresetNameEnemySkull: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("skull"),
				},
				Health: 2,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &zelduh.MovementConfig{
					Direction:    zelduh.DirectionDown,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "random",
				},
			}
		},
		PresetNameEnemySpinner: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("spinner"),
				},
				Invincible: true,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &zelduh.MovementConfig{
					Direction:    zelduh.DirectionRight,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "left-right",
				},
			}
		},
		PresetNameUICoin: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryHeart,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Hitbox: &zelduh.HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("uiCoin"),
				},
			}
		},
		PresetNameDialogCorner: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryIgnore,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					// TODO document the "default" animation key
					"default": zelduh.GetSpriteSet("dialogCorner"),
				},
			}
		},
		PresetNameDialogSide: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryIgnore,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					// TODO document the "default" animation key
					"default": zelduh.GetSpriteSet("dialogSide"),
				},
			}
		},
		"square": func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryIgnore,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				// TODO draw a square!
			}
		},
		PresetNameWarpStone: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryWarp,
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("warpStone"),
				},
			}
		},
		PresetNamePuzzleBox: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryMovableObstacle,
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("puzzleBox"),
				},
				Movement: &zelduh.MovementConfig{
					Speed:    1.0,
					MaxMoves: int(tileSize) / 2,
					MaxSpeed: 2.0,
				},
			}
		},
		PresetNameFloorSwitch: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryCollisionSwitch,
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("floorSwitch"),
				},
				Toggleable: true,
			}
		},
		// this is an impassable obstacle that can be toggled "remotely"
		// it has two visual states that coincide with each toggle state
		PresetNameToggleObstacle: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			// TODO get this working again
			return zelduh.EntityConfig{
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("toggleObstacle"),
				},
				// Impassable: true,
				Toggleable: true,
			}
		},
	}
}
