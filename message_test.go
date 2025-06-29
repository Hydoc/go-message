package message

import (
	"errors"
	"testing"
)

const testMessage = "test message"
const secondTestMessage = "second test message"

func TestSimpleBus_Dispatch(t *testing.T) {
	t.Run("call spy without event", func(t *testing.T) {
		bus := NewBus()

		spyCalled := false

		spy := func(message Message) (*Message, error) {
			spyCalled = true
			return nil, nil
		}

		bus.Register(testMessage, spy)

		err := bus.Dispatch(New(testMessage, nil))

		if err != nil {
			t.Fatal(err)
		}

		if !spyCalled {
			t.Error("Spy was not called")
		}
	})

	t.Run("call spy with event", func(t *testing.T) {
		bus := NewBus()

		timesCalledSpy := 0

		spy := func(message Message) (*Message, error) {
			timesCalledSpy++
			msg := New(secondTestMessage, nil)
			return &msg, nil
		}

		secondSpy := func(message Message) (*Message, error) {
			timesCalledSpy++
			return nil, nil
		}

		bus.Register(testMessage, spy)
		bus.Register(secondTestMessage, secondSpy)

		err := bus.Dispatch(New(testMessage, nil))

		if err != nil {
			t.Fatal(err)
		}

		if timesCalledSpy != 2 {
			t.Error("expected 2 spy calls")
		}
	})

	t.Run("handler returns error", func(t *testing.T) {
		bus := NewBus()

		bus.Register(testMessage, func(message Message) (*Message, error) {
			return nil, errors.New("handler error")
		})

		err := bus.Dispatch(New(testMessage, nil))

		if err == nil {
			t.Error("expected error")
		}

		if err.Error() != "handler error" {
			t.Error("expected handler error")
		}
	})

	t.Run("not found message", func(t *testing.T) {
		bus := NewBus()

		err := bus.Dispatch(Message{})

		if err == nil {
			t.Error("expected an error")
		}

		if err.Error() != "message not found" {
			t.Error("expected message not found error")
		}
	})
}
