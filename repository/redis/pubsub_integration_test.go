package redis

import (
	"context"
	"testing"
	"time"
)

func Test_PubSubRepository(t *testing.T) {
	repo := NewPubSubRepository(client)
	ctx := context.Background()

	channel := "testChannel"
	message := "testMessage"

	pubsub := repo.Subscribe(ctx, channel)
	defer pubsub.Close()

	time.Sleep(5 * time.Second) // PublishとSubscribeの間に少し遅延を入れます

	err := repo.Publish(ctx, channel, message)
	ValidateErr(t, err, nil)

	ch := pubsub.Channel()
	select {
	case msg := <-ch:
		if msg.Payload != message {
			t.Errorf("Subscribe() \n got = %v,\n want = %v", msg.Payload, message)
		}
	case <-time.After(10 * time.Second):
		t.Error("Timeout waiting for message")
	}
}
