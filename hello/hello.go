package hello

import (
	"context"

	"encore.dev/pubsub"
)

type Response struct {
	Message string
}

//encore:api public path=/hello/:name
func World(ctx context.Context, name string) (*Response, error) {
	msg := "Hello, " + name + "!"

	ev := &GreetingEvent{Name: name, Message: msg}
	if _, err := Greetings.Publish(ctx, ev); err != nil {
		return nil, err
	}

	return &Response{Message: msg}, nil
}

type GreetingEvent struct {
	Name    string
	Message string
}

var Greetings = pubsub.NewTopic[*GreetingEvent]("greetings", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})
