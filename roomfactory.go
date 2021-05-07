package zelduh

type RoomFactory struct {
	lastID RoomID
}

func NewRoomFactory() *RoomFactory {
	return &RoomFactory{
		lastID: 0,
	}
}

func (r *RoomFactory) newID() RoomID {
	r.lastID++
	return r.lastID
}

// NewRoom builds a new Room
func (r *RoomFactory) NewRoom(tmxFileName RoomName, entityConfigs ...EntityConfig) *Room {
	return &Room{
		ID:             r.newID(),
		TMXFileName:    tmxFileName,
		connectedRooms: &ConnectedRooms{},
		EntityConfigs:  entityConfigs,
	}
}
