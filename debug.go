package zelduh

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Debug struct {
	debugGridCellCachePopulated bool
	debugGridCellCache          []*imdraw.IMDraw
	debugTxt                    *text.Text
	activeSpaceRectangle        ActiveSpaceRectangle
	tileSize                    float64
	window                      *pixelgl.Window
}

func NewDebug(
	activeSpaceRectangle ActiveSpaceRectangle,
	tileSize float64,
	window *pixelgl.Window,
) *Debug {
	totalCells := int(activeSpaceRectangle.Width / tileSize * activeSpaceRectangle.Height / tileSize)
	// debugGridCellCachePopulated := false
	// debugGridCellCache := []*imdraw.IMDraw{}
	var debugGridCellCache []*imdraw.IMDraw = make([]*imdraw.IMDraw, totalCells)

	debugTxtOrigin := pixel.V(20, 50)
	debugTxt := text.New(debugTxtOrigin, text.Atlas7x13)

	return &Debug{
		debugGridCellCachePopulated: false,
		debugGridCellCache:          debugGridCellCache,
		debugTxt:                    debugTxt,
		activeSpaceRectangle:        activeSpaceRectangle,
		tileSize:                    tileSize,
		window:                      window,
	}
}

// Draw and overlay representing the virtual grid
func (d *Debug) DrawGrid() {
	// win.Clear(colornames.White)

	actualOriginX := d.activeSpaceRectangle.X
	actualOriginY := d.activeSpaceRectangle.Y

	totalColumns := d.activeSpaceRectangle.Width / d.tileSize
	totalRows := d.activeSpaceRectangle.Height / d.tileSize

	cacheIndex := 0

	var x float64 = 0
	var y float64 = 0

	if !(d.debugGridCellCachePopulated) {
		fmt.Println("building cache")
		for ; x < totalColumns; x++ {
			cellX := actualOriginX + (x * d.tileSize)
			cellY := actualOriginY + (y * d.tileSize)

			rect := pixel.R(cellX, cellY, cellX+d.tileSize, cellY+d.tileSize)

			imdraw := d.buildDebugGridCell(d.window, rect)
			d.debugGridCellCache[cacheIndex] = imdraw
			cacheIndex++

			if (x == (totalColumns - 1)) && (y < (totalRows - 1)) {
				x = -1
				y++
			}
		}
		d.debugGridCellCachePopulated = true
		fmt.Println("cache built")
	} else {
		for _, imdraw := range d.debugGridCellCache {
			imdraw.Draw(d.window)
		}
	}

	d.debugTxt.Clear()
	for ; x < totalColumns; x++ {
		message := fmt.Sprintf("%d,%d", int(x), int(y))
		fmt.Fprintln(d.debugTxt, message)
		matrix := pixel.IM.Moved(
			pixel.V(
				(actualOriginX-18)+(x*d.tileSize),
				(actualOriginY-d.tileSize)+(y*d.tileSize),
			),
		)
		d.debugTxt.Color = colornames.White
		d.debugTxt.Draw(d.window, matrix)
		d.debugTxt.Clear()

		if (x == (totalColumns - 1)) && (y < (totalRows - 1)) {
			x = -1
			y++
		}
	}

}

func (d *Debug) buildDebugGridCell(win *pixelgl.Window, rect pixel.Rect) *imdraw.IMDraw {

	imdraw := imdraw.New(nil)
	imdraw.Color = colornames.Blue

	imdraw.Push(rect.Min)
	imdraw.Push(rect.Max)

	imdraw.Rectangle(1)
	// imdraw.Draw(win)

	return imdraw
}
