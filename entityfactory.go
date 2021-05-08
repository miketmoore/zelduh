package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/core/direction"
	"github.com/miketmoore/zelduh/core/entity"
	"golang.org/x/image/colornames"
)

type EntityFactory struct {
	systemsManager  *SystemsManager
	temporarySystem *TemporarySystem
	movementSystem  *MovementSystem
	tileSize        float64
	frameRate       int
}

func NewEntityFactory(
	systemsManager *SystemsManager,
	temporarySystem *TemporarySystem,
	movementSystem *MovementSystem,
	tileSize float64,
	frameRate int,
) EntityFactory {
	return EntityFactory{
		systemsManager:  systemsManager,
		temporarySystem: temporarySystem,
		movementSystem:  movementSystem,
		tileSize:        tileSize,
		frameRate:       frameRate,
	}
}

func (ef *EntityFactory) NewEntityFromConfig(entityConfig EntityConfig, frameRate int) Entity {
	return ef.buildEntityFromConfig(
		entityConfig,
		ef.systemsManager.NewEntityID(),
		frameRate,
	)
}

func (ef *EntityFactory) buildEntityFromConfig(c EntityConfig, id entity.EntityID, frameRate int) Entity {
	entity := Entity{
		id:                   id,
		Category:             c.Category,
		componentIgnore:      NewComponentIgnore(c.Ignore),
		componentCoordinates: NewComponentCoordinates(c.Coordinates.X, c.Coordinates.Y),
		componentDimensions:  NewComponentDimensions(c.Dimensions.Width, c.Dimensions.Height),
		componentRectangle: NewComponentRectangle(
			c.Coordinates.X,
			c.Coordinates.Y,
			c.Dimensions.Width,
			c.Dimensions.Height,
		),

		// Create default shape and color
		// useful for debugging
		// might want to remove this later... not sure if creating
		// shapes that aren't being used increases heap memory
		componentShape: NewComponentShape(),
		componentColor: NewComponentColor(colornames.Greenyellow),
	}

	if c.Expiration > 0 {
		entity.componentTemporary = NewComponentTemporary(c.Expiration)
	}

	if c.Category == CategoryWarp {
		entity.componentEnabled = NewComponentEnabled(true)
	}

	if c.Health > 0 {
		entity.componentHealth = NewComponentHealth(c.Health)
	}

	if c.Hitbox != nil {
		entity.componentHitbox = NewComponentHitbox(c.Hitbox.Radius, float64(c.Hitbox.CollisionWithRectMod))
	}

	if c.Transform != nil {
		entity.componentRotation = NewComponentRotation(c.Transform.Rotation)
	}

	if c.Toggleable {
		entity.componentToggler = NewComponentToggler(c.Toggled)
	}

	entity.componentInvincible = NewComponentInvincible(c.Invincible)

	if c.Movement != nil {
		entity.componentMovement = NewComponentMovement(
			c.Movement.Direction,
			c.Movement.Speed,
			c.Movement.MaxSpeed,
			c.Movement.MaxMoves,
			c.Movement.RemainingMoves,
			c.Movement.HitSpeed,
			c.Movement.MovingFromHit,
			c.Movement.HitBackMoves,
			c.MovementPatternName,
		)
	}

	if c.Coins {
		entity.componentCoins = NewComponentCoins(0)
	}

	if c.Dash != nil {
		entity.componentDash = NewComponentDash(
			c.Dash.Charge,
			c.Dash.MaxCharge,
			c.Dash.SpeedMod,
		)
	}

	// An animation is a sprite graphic that may have one or more frames
	// so technically it might not be an animation
	if c.Animation != nil {
		entity.componentAnimation = NewComponentAnimation(c.Animation, frameRate)
	}

	return entity
}

func (ef *EntityFactory) CreateCoin(
	v pixel.Vec,
) {
	coordinates := Coordinates{
		X: v.X / ef.tileSize,
		Y: v.Y / ef.tileSize,
	}
	// coin := ef.NewEntityFromPresetName("coin", coordinates, ef.frameRate)
	coinConfig := ef.PresetCoin()(coordinates)
	coin := ef.NewEntityFromConfig(coinConfig, ef.frameRate)
	ef.systemsManager.AddEntity(coin)
}

func (ef *EntityFactory) CreateUICoin() {
	presetFn := ef.PresetUiCoin()
	entityConfig := presetFn(Coordinates{X: 4, Y: 14})
	coin := ef.NewEntityFromConfig(entityConfig, ef.frameRate)
	ef.systemsManager.AddEntity(coin)
}

func (ef *EntityFactory) CreateExplosion(
	entityID entity.EntityID,
) {
	explosion := ef.NewEntityFromConfig(ef.PresetExplosion()(NewCoordinates(0, 0)), ef.frameRate)

	ef.temporarySystem.SetExpiration(
		explosion.ID(),
		len(explosion.componentAnimation.ComponentAnimationByName["default"].Frames),
		func() {
			ef.CreateCoin(explosion.componentRectangle.Rect.Min)
		},
	)

	explosion.componentDimensions = NewComponentDimensions(ef.tileSize, ef.tileSize)
	enemyComponentRectangle, _ := ef.movementSystem.ComponentRectangle(entityID)
	explosion.componentRectangle = &componentRectangle{
		Rect: enemyComponentRectangle.Rect,
	}

	ef.systemsManager.AddEntity(explosion)
}

type BuildWarpFn func(
	warpToRoomID RoomID,
	coordinates Coordinates,
	hitboxRadius float64,
) EntityConfig

func (ef *EntityFactory) buildWarpFnFactory(
	dimensions Dimensions,
) BuildWarpFn {

	return func(
		warpToRoomID RoomID,
		coordinates Coordinates,
		hitboxRadius float64,
	) EntityConfig {
		return EntityConfig{
			Category:     CategoryWarp,
			WarpToRoomID: warpToRoomID,
			Dimensions:   dimensions,
			Coordinates:  coordinates,
			Hitbox: &HitboxConfig{
				Radius: hitboxRadius,
			},
		}
	}
}

type BuildWarpStoneFn func(
	WarpToRoomID RoomID,
	coordinates Coordinates,
	HitBoxRadius float64,
) EntityConfig

func (ef *EntityFactory) buildWarpStoneFnFactory() BuildWarpStoneFn {
	return func(
		WarpToRoomID RoomID,
		coordinates Coordinates,
		HitBoxRadius float64,
	) EntityConfig {
		presetFn := ef.PresetWarpStone()
		e := presetFn(Coordinates{X: coordinates.X, Y: coordinates.Y})
		e.WarpToRoomID = 6
		e.Hitbox.Radius = 5
		return e
	}
}

func (ef *EntityFactory) buildCoordinates(coordinates Coordinates) Coordinates {
	return Coordinates{
		X: ef.tileSize * coordinates.X,
		Y: ef.tileSize * coordinates.Y,
	}
}

func (ef *EntityFactory) PresetArrow() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryArrow,
			Movement: &MovementConfig{
				Direction: direction.DirectionDown,
				Speed:     0.0,
			},
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
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
	}
}

func (ef *EntityFactory) PresetBomb() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryBomb,
			Movement: &MovementConfig{
				Direction: direction.DirectionDown,
				Speed:     0.0,
			},
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
			Animation: AnimationConfig{
				"default": GetSpriteSet("bomb"),
			},
			Hitbox: &HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	}
}

func (ef *EntityFactory) PresetCoin() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryCoin,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
			Animation: AnimationConfig{
				"default": GetSpriteSet("coin"),
			},
		}
	}
}

func (ef *EntityFactory) PresetExplosion() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category:   CategoryExplosion,
			Expiration: 12,
			Animation: AnimationConfig{
				"default": GetSpriteSet("explosion"),
			},
		}
	}
}

func (ef *EntityFactory) PresetObstacle() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryObstacle,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
		}
	}
}

func (ef *EntityFactory) PresetPlayer() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryPlayer,
			Health:   3,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
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
	}
}

func (ef *EntityFactory) PresetSword() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategorySword,
			Movement: &MovementConfig{
				Direction: direction.DirectionDown,
				Speed:     0.0,
			},
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
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
	}
}

func (ef *EntityFactory) PresetEyeBurrower() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryEnemy,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
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
	}
}

func (ef *EntityFactory) PresetHeart() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryHeart,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
			Hitbox: &HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: AnimationConfig{
				"default": GetSpriteSet("heart"),
			},
		}

	}
}

func (ef *EntityFactory) PresetSkeleton() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryEnemy,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
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
	}
}

func (ef *EntityFactory) PresetSkull() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryEnemy,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
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
	}
}

func (ef *EntityFactory) PresetSpinner() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryEnemy,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
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
	}
}

func (ef *EntityFactory) PresetUiCoin() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category: CategoryHeart,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
			Hitbox: &HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: AnimationConfig{
				"default": GetSpriteSet("uiCoin"),
			},
		}
	}
}

func (ef *EntityFactory) PresetDialogCorner(degrees float64) EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		config := EntityConfig{
			Category: CategoryIgnore,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
			Animation: AnimationConfig{
				// TODO document the "default" animation key
				"default": GetSpriteSet("dialogCorner"),
			},
		}
		if degrees != 0 {
			config.Transform = &Transform{
				Rotation: degrees,
			}
		}
		return config
	}
}

func (ef *EntityFactory) PresetDialogSide(degrees float64) EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		config := EntityConfig{
			Category: CategoryIgnore,
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Coordinates: ef.buildCoordinates(coordinates),
			Animation: AnimationConfig{
				// TODO document the "default" animation key
				"default": GetSpriteSet("dialogSide"),
			},
		}
		if degrees != 0 {
			config.Transform = &Transform{
				Rotation: degrees,
			}
		}
		return config
	}
}

func (ef *EntityFactory) PresetWarpStone() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category:    CategoryWarp,
			Coordinates: ef.buildCoordinates(coordinates),
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Hitbox: &HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Animation: AnimationConfig{
				"default": GetSpriteSet("warpStone"),
			},
		}
	}
}

func (ef *EntityFactory) PresetPuzzleBox() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category:    CategoryMovableObstacle,
			Coordinates: ef.buildCoordinates(coordinates),
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Animation: AnimationConfig{
				"default": GetSpriteSet("puzzleBox"),
			},
			Movement: &MovementConfig{
				Speed:    1.0,
				MaxMoves: int(ef.tileSize) / 2,
				MaxSpeed: 2.0,
			},
		}
	}
}

func (ef *EntityFactory) PresetFloorSwitch() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		return EntityConfig{
			Category:    CategoryCollisionSwitch,
			Coordinates: ef.buildCoordinates(coordinates),
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Animation: AnimationConfig{
				"default": GetSpriteSet("floorSwitch"),
			},
			Toggleable: true,
		}
	}
}

func (ef *EntityFactory) PresetToggleObstacle() EntityConfigPresetFn {
	return func(coordinates Coordinates) EntityConfig {
		// TODO get this working again
		return EntityConfig{
			Coordinates: ef.buildCoordinates(coordinates),
			Dimensions: Dimensions{
				Width:  ef.tileSize,
				Height: ef.tileSize,
			},
			Animation: AnimationConfig{
				"default": GetSpriteSet("toggleObstacle"),
			},
			// Impassable: true,
			Toggleable: true,
		}
	}
}
