package zelduh

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/deanobob/tmxreader"
	"github.com/faiface/pixel"
)

type TmxMapByNameMap map[string]tmxreader.TmxMap

// Load loads the tmx files
func Load(tilemapFiles []string, tilemapDir string) (TmxMapByNameMap, error) {
	tmxMapData := map[string]tmxreader.TmxMap{}
	for _, name := range tilemapFiles {
		path := fmt.Sprintf("%s%s.tmx", tilemapDir, name)
		parsed, err := readAndParseTmxFile(path)
		if err != nil {
			fmt.Println(err)
			return TmxMapByNameMap{}, fmt.Errorf("error reading and parsing file=%s", path)
		}
		tmxMapData[name] = parsed
	}
	return tmxMapData, nil
}

func readAndParseTmxFile(filename string) (tmxreader.TmxMap, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return tmxreader.TmxMap{}, fmt.Errorf("error reading file=%s\n", filename)
	}

	tmxMap, err := tmxreader.Parse(raw)
	if err != nil {
		fmt.Println(err)
		return tmxreader.TmxMap{}, fmt.Errorf("error parsing TMX data")
	}

	return tmxMap, nil
}

type mapDrawData struct {
	Rect     pixel.Rect
	SpriteID int
}

// MapData represents data for one map
type MapData struct {
	Name string
	Data []mapDrawData
}

type MapDrawData map[TMXFileName]MapData

// BuildMapDrawData builds draw data and stores it in a map
func BuildMapDrawData(dir string, files []string, tileSize float64) (MapDrawData, error) {

	// load all TMX file data for each map
	tmxMapData, err := Load(files, dir)
	if err != nil {
		fmt.Println(err)
		return MapDrawData{}, fmt.Errorf("error loading TMX files")
	}

	all := MapDrawData{}

	for mapName, mapData := range tmxMapData {

		md := MapData{
			Name: mapName,
			Data: []mapDrawData{},
		}

		layers := mapData.Layers
		for _, layer := range layers {

			records := ParseCSV(strings.TrimSpace(layer.Data.Value) + ",")
			for row := 0; row <= len(records); row++ {
				if len(records) > row {
					for col := 0; col < len(records[row])-1; col++ {
						y := float64(11-row) * tileSize
						x := float64(col) * tileSize

						record := records[row][col]
						spriteID, err := strconv.Atoi(record)
						if err != nil {
							panic(err)
						}
						mrd := mapDrawData{
							Rect:     pixel.R(x, y, x+tileSize, y+tileSize),
							SpriteID: spriteID,
						}
						md.Data = append(md.Data, mrd)
					}
				}

			}
			all[TMXFileName(mapName)] = md
		}
	}

	return all, nil
}
