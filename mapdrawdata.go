package zelduh

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/faiface/pixel"
	"github.com/miketmoore/zelduh/core/csv"
	"github.com/miketmoore/zelduh/core/tmx"
)

type mapDrawData struct {
	Rect     pixel.Rect
	SpriteID int
}

// MapData represents data for one map
type MapData struct {
	Name tmx.TMXFileName
	Data []mapDrawData
}

type MapDrawData map[tmx.TMXFileName]MapData

// BuildMapDrawData builds draw data and stores it in a map
func BuildMapDrawData(dir string, files []string, tileSize float64) (MapDrawData, error) {

	// load all TMX file data for each map
	tmxMapData, err := tmx.ReadAndParse(files, dir)
	if err != nil {
		fmt.Println(err)
		return MapDrawData{}, fmt.Errorf("error loading TMX files")
	}

	all := MapDrawData{}

	for mapName, rawMapData := range tmxMapData {

		mapData := MapData{
			Name: mapName,
			Data: []mapDrawData{},
		}

		layers := rawMapData.Layers
		for _, layer := range layers {

			records := csv.ParseCSV(strings.TrimSpace(layer.Data.Value) + ",")
			for row := 0; row <= len(records); row++ {
				if len(records) > row {
					for col := 0; col < len(records[row])-1; col++ {
						y := float64(11-row) * tileSize
						x := float64(col) * tileSize

						record := records[row][col]
						spriteID, err := strconv.Atoi(record)
						if err != nil {
							fmt.Println(err)
							return MapDrawData{}, fmt.Errorf("error parsing string to int; string=%s", record)
						}
						mrd := mapDrawData{
							Rect:     pixel.R(x, y, x+tileSize, y+tileSize),
							SpriteID: spriteID,
						}
						mapData.Data = append(mapData.Data, mrd)
					}
				}

			}
			all[tmx.TMXFileName(mapName)] = mapData
		}
	}

	return all, nil
}
