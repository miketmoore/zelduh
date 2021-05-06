package zelduh

import "github.com/miketmoore/zelduh/core/entity"

type RoomWarps map[entity.EntityID]EntityConfig

func NewRoomWarps() RoomWarps {
	return RoomWarps{}
}
