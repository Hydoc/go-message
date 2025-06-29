package message

import "errors"

type Message struct {
	Type    string
	Payload any
}

type Bus interface {
	Dispatch(message Message) error
	Register(message string, handler func(message Message) (*Message, error))
}

type SimpleBus struct {
	handlers map[string]func(message Message) (*Message, error)
}

func (b *SimpleBus) Dispatch(message Message) error {
	handler, ok := b.handlers[message.Type]
	if !ok {
		return errors.New("message type not found")
	}

	event, err := handler(message)
	if err != nil {
		return err
	}

	if event != nil {
		return b.Dispatch(*event)
	}

	return nil
}

func (b *SimpleBus) Register(message string, handler func(message Message) (*Message, error)) {
	b.handlers[message] = handler
}

func NewBus() Bus {
	return &SimpleBus{
		handlers: make(map[string]func(message Message) (*Message, error)),
	}
}
