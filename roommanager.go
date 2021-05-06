package zelduh

type RoomManager struct {
	current, next RoomID
}

func NewRoomManager(current RoomID) *RoomManager {
	return &RoomManager{
		current: current,
	}
}

func (n *RoomManager) SetNext(next RoomID) {
	n.next = next
}

func (n *RoomManager) Current() RoomID {
	return n.current
}

func (n *RoomManager) Next() RoomID {
	return n.next
}

func (n *RoomManager) MoveToNext() {
	n.current = n.next
}
