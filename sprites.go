package zelduh

import (
	"fmt"
	"image"
	"math"
	"os"

	"github.com/faiface/pixel"
)

type Spritesheet map[int]*pixel.Sprite

func loadPicture(path string) pixel.Picture {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Could not open the picture:")
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Could not decode the picture:")
		fmt.Println(err)
		os.Exit(1)
	}
	return pixel.PictureDataFromImage(img)
}

// LoadAndBuildSpritesheet this is a map of pixel engine sprites
func LoadAndBuildSpritesheet(path string, tileSize float64) map[int]*pixel.Sprite {
	pic := loadPicture(path)

	cols := pic.Bounds().W() / tileSize
	rows := pic.Bounds().H() / tileSize

	maxIndex := (rows * cols) - 1.0

	index := maxIndex
	id := maxIndex + 1
	spritesheet := map[int]*pixel.Sprite{}
	for row := (rows - 1); row >= 0; row-- {
		for col := (cols - 1); col >= 0; col-- {
			x := col
			y := math.Abs(rows-row) - 1
			spritesheet[int(id)] = pixel.NewSprite(pic, pixel.R(
				x*tileSize,
				y*tileSize,
				x*tileSize+tileSize,
				y*tileSize+tileSize,
			))
			index--
			id--
		}
	}
	return spritesheet
}

// GetSpriteSet returns a sprite set by key
func GetSpriteSet(key string) []int {
	return spriteSets[key]
}

// spritesheet is 15 tiles wid
// each tile is 48 pixels square

var spriteSets = map[string][]int{
	"eyeburrower": {50, 50, 50, 91, 91, 91, 92, 92, 92, 93, 93, 93, 92, 92, 92},
	"explosion": {
		122, 122, 122,
		123, 123, 123,
		124, 124, 124,
		125, 125, 125,
	},
	"uiCoin":           {20},
	"skeleton":         {31, 32},
	"skull":            {36, 37, 38, 39},
	"spinner":          {51, 52},
	"puzzleBox":        {63},
	"warpStone":        {61},
	"playerUp":         {4, 195},
	"playerRight":      {3, 194},
	"playerDown":       {1, 192},
	"playerLeft":       {2, 193},
	"playerSwordUp":    {165},
	"playerSwordRight": {164},
	"playerSwordLeft":  {179},
	"playerSwordDown":  {180},
	"floorSwitch":      {112, 127},
	"toggleObstacle":   {144, 114},
	"swordUp":          {70},
	"swordRight":       {67},
	"swordDown":        {68},
	"swordLeft":        {69},
	"arrowUp":          {101},
	"arrowRight":       {100},
	"arrowDown":        {103},
	"arrowLeft":        {102},
	"bomb":             {138, 139, 140, 141},
	"coin":             {5, 5, 6, 6, 21, 21},
	"heart":            {106},
	"dialogCorner":     {11},
	"dialogSide":       {12},
}
