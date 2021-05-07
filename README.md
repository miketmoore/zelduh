# Zelduh

[![Go Report Card](https://goreportcard.com/badge/github.com/miketmoore/zelduh)](https://goreportcard.com/report/github.com/miketmoore/zelduh)

https://miketmoore.itch.io/zelduh

Zelduh is a tile based adventure game. 

## Install

```
go get -u github.com/miketmoore/zelduh/...
```

## Run

```
go run cmd/zelduh/*.go
```

## Run in debug mode

```
go run cmd/zelduh/*.go debug
```

## Controls

| Action | Keys |
| ---- | ---- |
| Confirm/Next | Enter |
| Walk | W, A, S, D |
| Sword | F | 
| Arrow | G |
| Dash | F + Space | 


## Notes

Screen is a 15 wide grid (X coordinates 0-14)

## How to add a room

1. Create a TMX file:
    - fixed size: 14 wide by 12 high
    - orientation: orthogonal
    - format: CSV
    - tile render order: top down
2. Put TMX file in `assets/tilemaps/`
3. Add filename without extension to list passed to `BuildMapDrawData`
    - Example, for file name `myMap01.tmx`, add `"myMap01"` to list
4. Create room in level where `Room.TMXFileName` is `"myMap01"` (name without extension)