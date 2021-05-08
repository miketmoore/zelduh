package tmx

import (
	"fmt"
	"io/ioutil"

	"github.com/deanobob/tmxreader"
)

type TMXFileName string

type tmxMapByNameMap map[TMXFileName]tmxreader.TmxMap

// ReadAndParse reads *.tmx files and parses them
func ReadAndParse(tilemapFiles []string, tilemapDir string) (tmxMapByNameMap, error) {
	tmxMapData := tmxMapByNameMap{}
	for _, name := range tilemapFiles {
		path := fmt.Sprintf("%s%s.tmx", tilemapDir, name)
		parsed, err := readAndParseTmxFile(path)
		if err != nil {
			fmt.Println(err)
			return tmxMapByNameMap{}, fmt.Errorf("error reading and parsing file=%s", path)
		}
		tmxMapData[TMXFileName(name)] = parsed
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
