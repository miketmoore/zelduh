package zelduh

type RoomWarps map[EntityID]EntityConfig

func NewRoomWarps() RoomWarps {
	return RoomWarps{}
}
