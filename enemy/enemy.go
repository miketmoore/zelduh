package enemy

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/mvmt"
	"github.com/miketmoore/zelduh/palette"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type stateName string

const (
	stateNameAppear stateName = "appear"
	stateNameActive stateName = "active"
	stateNamePause  stateName = "pause"
)

// Enemy represents one non-player character
type Enemy struct {
	// Size is the dimensions (square)
	Size float64
	// Start is the starting vector
	Start pixel.Vec
	// Last is the last vector
	Last pixel.Vec
	// LastDir is the last direction the Enemy was headed in
	LastDir mvmt.Direction
	// Shape is the view
	Shape             *imdraw.IMDraw
	Sprites           map[string]*pixel.Sprite
	WalkCycleCountMax int
	WalkCycleCount    int
	WalkCycleFlag     bool
	// Win is a pointer to the GUI window
	Win *pixelgl.Window

	speed float64

	AttackPower int
	Health      int

	totalMoves         int
	moveCounter        int
	currentState       stateName
	lastState          stateName
	appearDelay        int
	currentAppearDelay int
}

// New returns a new enemy
func New(win *pixelgl.Window, size, x, y, speed float64, health, attackPower int, sprites map[string]*pixel.Sprite) Enemy {
	enemy := Enemy{
		Win:                win,
		Size:               size,
		Start:              pixel.V(x, y),
		Last:               pixel.V(0, 0),
		Shape:              imdraw.New(nil),
		Sprites:            sprites,
		WalkCycleCountMax:  7,
		WalkCycleFlag:      false,
		speed:              speed,
		Health:             health,
		AttackPower:        attackPower,
		currentState:       stateNameAppear,
		appearDelay:        20,
		currentAppearDelay: 20,
	}
	enemy.WalkCycleCount = enemy.WalkCycleCountMax
	return enemy
}

// Reset updates all values to starting values
func (enemy *Enemy) Reset() {
	// TODO reset to random start position
	// need limits defined, probably on enemy constructor
	enemy.Last = enemy.Start
	enemy.currentState = stateNameAppear
	enemy.currentAppearDelay = enemy.appearDelay
}

// Hit registers a hit
func (enemy *Enemy) Hit(attackPower int) {
	enemy.Health -= attackPower
}

func (enemy *Enemy) IsDead() bool {
	return enemy.Health == 0
}

// Draw the enemy's current state
func (enemy *Enemy) Draw(minX, minY, maxX, maxY float64) {
	enemy.Shape.Clear()

	switch enemy.currentState {
	case stateNameAppear:
		enemy.Shape.Color = palette.Map[palette.Lightest]
		if enemy.currentAppearDelay > 0 {
			enemy.currentAppearDelay--
		} else {
			enemy.currentState = stateNameActive
		}
	case stateNameActive:
		enemy.Shape.Color = palette.Map[palette.Darkest]
		// move smoothly into next tile, then orient for next move
		stride := enemy.speed

		if enemy.totalMoves == 0 {
			// number of tiles to move until next change in movement
			maxMoves := 10
			enemy.totalMoves = r.Intn(maxMoves)

			// which direction to move
			directionIndex := r.Intn(4)
			switch directionIndex {
			case 0:
				enemy.LastDir = mvmt.DirectionYPos
			case 1:
				enemy.LastDir = mvmt.DirectionXPos
			case 2:
				enemy.LastDir = mvmt.DirectionYNeg
			case 3:
				enemy.LastDir = mvmt.DirectionXNeg
			}

		} else {
			if enemy.moveCounter > 0 {
				switch enemy.LastDir {
				case mvmt.DirectionYPos:
					if enemy.Last.Y+enemy.Size < maxY {
						enemy.Last = pixel.V(enemy.Last.X, enemy.Last.Y+stride)
					}
				case mvmt.DirectionXPos:
					if enemy.Last.X+enemy.Size < maxX {
						enemy.Last = pixel.V(enemy.Last.X+stride, enemy.Last.Y)
					}
				case mvmt.DirectionYNeg:
					if enemy.Last.Y >= minY {
						enemy.Last = pixel.V(enemy.Last.X, enemy.Last.Y-stride)
					}
				case mvmt.DirectionXNeg:
					if enemy.Last.X >= minX {
						enemy.Last = pixel.V(enemy.Last.X-stride, enemy.Last.Y)
					}
				}
				enemy.moveCounter--
			} else {
				enemy.totalMoves--
				enemy.moveCounter = int(enemy.Size)
			}

		}
	}

	enemy.Shape.Push(enemy.Last)
	enemy.Shape.Push(pixel.V(enemy.Last.X+enemy.Size, enemy.Last.Y+enemy.Size))
	enemy.Shape.Rectangle(0)
	// enemy.Shape.Draw(enemy.Win)

	// Do some walk cycle "math"
	if enemy.WalkCycleCount > 0 {
		enemy.WalkCycleCount--
	} else {
		enemy.WalkCycleCount = enemy.WalkCycleCountMax
		enemy.WalkCycleFlag = !enemy.WalkCycleFlag
	}

	// Figure out which walk cycle sprite to use
	var sprite *pixel.Sprite

	if enemy.LastDir == mvmt.DirectionYNeg {
		if enemy.WalkCycleFlag {
			sprite = enemy.Sprites["downA"]
		} else {
			sprite = enemy.Sprites["downB"]
		}
	} else if enemy.LastDir == mvmt.DirectionYPos {
		if enemy.WalkCycleFlag {
			sprite = enemy.Sprites["upA"]
		} else {
			sprite = enemy.Sprites["upB"]
		}
	} else if enemy.LastDir == mvmt.DirectionXPos {
		if enemy.WalkCycleFlag {
			sprite = enemy.Sprites["rightA"]
		} else {
			sprite = enemy.Sprites["rightB"]
		}
	} else if enemy.LastDir == mvmt.DirectionXNeg {
		if enemy.WalkCycleFlag {
			sprite = enemy.Sprites["leftA"]
		} else {
			sprite = enemy.Sprites["leftB"]
		}
	} else {
		sprite = enemy.Sprites["downA"]
	}

	matrix := pixel.IM.Moved(pixel.V(enemy.Last.X+enemy.Size/2, enemy.Last.Y+enemy.Size/2))
	sprite.Draw(enemy.Win, matrix)
}
