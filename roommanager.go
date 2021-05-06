package zelduh

import "fmt"

type RoomManager struct {
	current, next RoomID
}

func NewRoomManager(current RoomID) *RoomManager {
	return &RoomManager{
		current: current,
	}
}

func (n *RoomManager) SetNext(next RoomID) {
	fmt.Printf("RoomManager SetNext next=%d\n", next)
	n.next = next
}

func (n *RoomManager) Current() RoomID {
	return n.current
}

func (n *RoomManager) Next() RoomID {
	return n.next
}

func (n *RoomManager) MoveToNext() {
	fmt.Printf("RoomManager MoveToNext current=%d next=%d", n.current, n.next)
	n.current = n.next
}
