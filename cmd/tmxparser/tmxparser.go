package main

import (
	"fmt"
	"io/ioutil"

	"github.com/deanobob/tmxreader"
)

// type Reader interface {
// 	ReadFile(string) ([]byte, error)
// }

type Reader struct{}

func (r Reader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func main() {

	reader := Reader{}
	filename := "assets/tilemaps/overworld.tmx"
	raw, err := reader.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	tmxMap, err := tmxreader.Parse(raw)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%v\n", tmxMap)
	fmt.Printf("Version: %v\n", tmxMap.Version)
	fmt.Printf("Orientation: %v\n", tmxMap.Orientation)
	fmt.Printf("Width: %v\n", tmxMap.Width)
	fmt.Printf("Height: %v\n", tmxMap.Height)
	fmt.Printf("TileWidth: %v\n", tmxMap.TileWidth)
	fmt.Printf("TileHeight: %v\n", tmxMap.TileHeight)
	fmt.Printf("Total Properties: %d\n", len(tmxMap.Properties))
	fmt.Printf("Total Tilesets: %d\n", len(tmxMap.Tilesets))
	fmt.Printf("Total Layers: %d\n", len(tmxMap.Layers))
	fmt.Printf("Total ObjectGroups: %d\n", len(tmxMap.ObjectGroups))

	for i, tileset := range tmxMap.Tilesets {
		fmt.Printf("TmxTileset #%d\n", i)
		fmt.Printf("FirstGid: %v\n", tileset.FirstGid)
		fmt.Printf("Name: %v\n", tileset.Name)
		fmt.Printf("TileWidth: %v\n", tileset.TileWidth)
		fmt.Printf("TileHeight: %v\n", tileset.TileHeight)
		fmt.Printf("Images: %v\n", tileset.Images)
	}

	for i, layer := range tmxMap.Layers {
		fmt.Printf("TmxLayer #%d\n", i)
		fmt.Printf("Name: %v\n", layer.Name)
		fmt.Printf("Width: %v\n", layer.Width)
		fmt.Printf("Height: %v\n", layer.Height)
		fmt.Printf("Data.Encoding: %s\n", layer.Data.Encoding)
		fmt.Printf("Data.Value:\n")
		fmt.Println(layer.Data.Value)
	}

}
