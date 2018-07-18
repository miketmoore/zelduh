package zelduh

// GetSpriteSet returns a sprite set by key
func GetSpriteSet(key string) []int {
	return spriteSets[key]
}

var spriteSets = map[string][]int{
	"eyeburrower": []int{50, 50, 50, 91, 91, 91, 92, 92, 92, 93, 93, 93, 92, 92, 92},
	"explosion": []int{
		122, 122, 122,
		123, 123, 123,
		124, 124, 124,
		125, 125, 125,
	},
	"uiCoin":           []int{20},
	"skeleton":         []int{31, 32},
	"skull":            []int{36, 37, 38, 39},
	"spinner":          []int{51, 52},
	"puzzleBox":        []int{63},
	"warpStone":        []int{61},
	"playerUp":         []int{4, 195},
	"playerRight":      []int{3, 194},
	"playerDown":       []int{1, 192},
	"playerLeft":       []int{2, 193},
	"playerSwordUp":    []int{165},
	"playerSwordRight": []int{164},
	"playerSwordLeft":  []int{179},
	"playerSwordDown":  []int{180},
	"floorSwitch":      []int{112, 127},
	"toggleObstacle":   []int{144, 114},
	"swordUp":          []int{70},
	"swordRight":       []int{67},
	"swordDown":        []int{68},
	"swordLeft":        []int{69},
	"arrowUp":          []int{101},
	"arrowRight":       []int{100},
	"arrowDown":        []int{103},
	"arrowLeft":        []int{102},
	"bomb":             []int{138, 139, 140, 141},
	"coin":             []int{5, 5, 6, 6, 21, 21},
	"heart":            []int{106},
}
