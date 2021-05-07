package zelduh

// System is an interface
type System interface {
	Update() error
	AddEntity(Entity)
}
