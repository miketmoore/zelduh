package message

//A Handler is used to dispatch a message to the subscribed handler.
type Handler func(msg Message)

// A Message is used to send messages within the MessageManager
type Message interface {
	Type() string
}

// Manager manages messages and subscribed handlers
type Manager struct {
	listeners map[string][]Handler
}

// Dispatch sends a message to all subscribed handlers of the message's type
func (mm *Manager) Dispatch(message Message) {
	handlers := mm.listeners[message.Type()]

	for _, handler := range handlers {
		handler(message)
	}
}

// Listen subscribes to the specified message type and calls the specified handler when fired
func (mm *Manager) Listen(messageType string, handler Handler) {
	if mm.listeners == nil {
		mm.listeners = make(map[string][]Handler)
	}
	mm.listeners[messageType] = append(mm.listeners[messageType], handler)
}
