package event

import "github.com/go-gl/glfw/v3.2/glfw"

type Event interface{}

type EventListener interface {
	Notify(Event)
}

type KeyEvent struct {
	Key    glfw.Key
	Action glfw.Action
}

type EventHandler struct {
	eventListeners []EventListener
}

func (e *EventHandler) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	for _, el := range e.eventListeners {
		el.Notify(KeyEvent{key, action})
	}
}

func (e *EventHandler) Subscribe(el EventListener) {
	e.eventListeners = append(e.eventListeners, el)
}

func CreateEventHandler() *EventHandler {
	return &EventHandler{}
}
