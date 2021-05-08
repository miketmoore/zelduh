package zelduh

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/core/direction"
)

const (
	PresetNameArrow            PresetName = "arrow"
	PresetNameBomb             PresetName = "bomb"
	PresetNameCoin             PresetName = "coin"
	PresetNameExplosion        PresetName = "explosion"
	PresetNameObstacle         PresetName = "obstacle"
	PresetNamePlayer           PresetName = "player"
	PresetNameFloorSwitch      PresetName = "floorSwitch"
	PresetNameToggleObstacle   PresetName = "toggleObstacle"
	PresetNamePuzzleBox        PresetName = "puzzleBox"
	PresetNameWarpStone        PresetName = "warpStone"
	PresetNameUICoin           PresetName = "uiCoin"
	PresetNameEnemySpinner     PresetName = "spinner"
	PresetNameEnemySkull       PresetName = "skull"
	PresetNameEnemySkeleton    PresetName = "skeleton"
	PresetNameHeart            PresetName = "heart"
	PresetNameEnemyEyeBurrower PresetName = "eyeBurrower"
	PresetNameSword            PresetName = "sword"
	PresetNameDialogCorner     PresetName = "dialogCorner"
	PresetNameDialogSide       PresetName = "dialogSide"
)

type BuildWarpFn func(
	warpToRoomID RoomID,
	coordinates Coordinates,
	hitboxRadius float64,
) EntityConfig

// TODO move this to a higher level configuration location
func BuildEntityConfigPresetFnsMap(tileSize float64) map[PresetName]EntityConfigPresetFn {

	dimensions := Dimensions{
		Width:  tileSize,
		Height: tileSize,
	}

	buildCoordinates := func(coordinates Coordinates) Coordinates {
		return Coordinates{
			X: tileSize * coordinates.X,
			Y: tileSize * coordinates.Y,
		}
	}

	return map[PresetName]EntityConfigPresetFn{
		PresetNameArrow: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category: CategoryArrow,
				Movement: &MovementConfig{
					Direction: direction.DirectionDown,
					Speed:     0.0,
				},
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
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
		PresetNameBomb: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category: CategoryBomb,
				Movement: &MovementConfig{
					Direction: direction.DirectionDown,
					Speed:     0.0,
				},
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: AnimationConfig{
					"default": GetSpriteSet("bomb"),
				},
				Hitbox: &HitboxConfig{
					Radius: 5,
				},
				Ignore: true,
			}
		},
		PresetNameCoin: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryCoin,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: AnimationConfig{
					"default": GetSpriteSet("coin"),
				},
			}
		},
		PresetNameExplosion: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:   CategoryExplosion,
				Expiration: 12,
				Animation: AnimationConfig{
					"default": GetSpriteSet("explosion"),
				},
			}
		},
		PresetNameObstacle: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryObstacle,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
			}
		},
		PresetNamePlayer: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryPlayer,
				Health:      3,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Hitbox: &HitboxConfig{
					Box:                  imdraw.New(nil),
					Radius:               15,
					CollisionWithRectMod: 5,
				},
				Movement: &MovementConfig{
					Direction: direction.DirectionDown,
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
		PresetNameSword: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category: CategorySword,
				Movement: &MovementConfig{
					Direction: direction.DirectionDown,
					Speed:     0.0,
				},
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
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
		PresetNameEnemyEyeBurrower: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: AnimationConfig{
					"default": GetSpriteSet("eyeburrower"),
				},
				Health: 2,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &MovementConfig{
					Direction:           direction.DirectionDown,
					Speed:               1.0,
					MaxSpeed:            1.0,
					HitSpeed:            10.0,
					HitBackMoves:        10,
					MaxMoves:            100,
					MovementPatternName: "random",
				},
			}
		},
		PresetNameHeart: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryHeart,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Hitbox: &HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: AnimationConfig{
					"default": GetSpriteSet("heart"),
				},
			}

		},
		PresetNameEnemySkeleton: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: AnimationConfig{
					"default": GetSpriteSet("skeleton"),
				},
				Health: 2,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &MovementConfig{
					Direction:           direction.DirectionDown,
					Speed:               1.0,
					MaxSpeed:            1.0,
					HitSpeed:            10.0,
					HitBackMoves:        10,
					MaxMoves:            100,
					MovementPatternName: "random",
				},
			}
		},
		PresetNameEnemySkull: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: AnimationConfig{
					"default": GetSpriteSet("skull"),
				},
				Health: 2,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &MovementConfig{
					Direction:           direction.DirectionDown,
					Speed:               1.0,
					MaxSpeed:            1.0,
					HitSpeed:            10.0,
					HitBackMoves:        10,
					MaxMoves:            100,
					MovementPatternName: "random",
				},
			}
		},
		PresetNameEnemySpinner: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: AnimationConfig{
					"default": GetSpriteSet("spinner"),
				},
				Invincible: true,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &MovementConfig{
					Direction:           direction.DirectionRight,
					Speed:               1.0,
					MaxSpeed:            1.0,
					HitSpeed:            10.0,
					HitBackMoves:        10,
					MaxMoves:            100,
					MovementPatternName: "left-right",
				},
			}
		},
		PresetNameUICoin: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryHeart,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Hitbox: &HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: AnimationConfig{
					"default": GetSpriteSet("uiCoin"),
				},
			}
		},
		PresetNameDialogCorner: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryIgnore,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: AnimationConfig{
					// TODO document the "default" animation key
					"default": GetSpriteSet("dialogCorner"),
				},
			}
		},
		PresetNameDialogSide: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryIgnore,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: AnimationConfig{
					// TODO document the "default" animation key
					"default": GetSpriteSet("dialogSide"),
				},
			}
		},
		// "square": func(coordinates Coordinates) EntityConfig {
		// 	return EntityConfig{
		// 		Category:    CategoryIgnore,
		// 		Dimensions:  dimensions,
		// 		Coordinates: buildCoordinates(coordinates),
		// 		// TODO draw a square!
		// 	}
		// },
		PresetNameWarpStone: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryWarp,
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Hitbox: &HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Animation: AnimationConfig{
					"default": GetSpriteSet("warpStone"),
				},
			}
		},
		PresetNamePuzzleBox: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryMovableObstacle,
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
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
		PresetNameFloorSwitch: func(coordinates Coordinates) EntityConfig {
			return EntityConfig{
				Category:    CategoryCollisionSwitch,
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Animation: AnimationConfig{
					"default": GetSpriteSet("floorSwitch"),
				},
				Toggleable: true,
			}
		},
		// this is an impassable obstacle that can be toggled "remotely"
		// it has two visual states that coincide with each toggle state
		PresetNameToggleObstacle: func(coordinates Coordinates) EntityConfig {
			// TODO get this working again
			return EntityConfig{
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Animation: AnimationConfig{
					"default": GetSpriteSet("toggleObstacle"),
				},
				// Impassable: true,
				Toggleable: true,
			}
		},
	}
}
