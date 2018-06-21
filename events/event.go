package events

import "sync"

type EventType uint16

const(
	EventSaveBlock 		=	iota
	EventReplyTx
	Eventxx
)

type EventFunc func(v interface{})
type Subscriber struct {
	fun 	EventFunc
	listener []string
}
type Event struct{
	m		sync.RWMutex
	subscribers map[EventType]Subscriber
}

var event *Event
var once sync.Once

func NewEvent() *Event{
	once.Do(func() {
		event = &Event{subscribers:make(map[EventType]Subscriber)}
	})
	return event
}

func (e *Event) Subscribe(eventType EventType,eventFun EventType)bool{
	e.m.Lock()
	defer e.m.Unlock()

}