package persistenceServices

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func GetPubsubConnection(gcpProjectID string) (*PubsubConnection, error) {

	var err error

	p := PubsubConnection{}

	if gcpProjectID == "" {
		err = errors.New("GetPubsubConnection ERROR: GCP_PROJECT not set!")
	}

	if err == nil {

		p.ctx = context.Background()

		client, err := pubsub.NewClient(p.ctx, gcpProjectID)
		p.client = *client

		if err != nil {
			fmt.Println(fmt.Sprintf("GetPubsubConnection.1: Error %v", err))
			return &p, nil
		} else {

		}

	} else {
		fmt.Println(fmt.Sprintf("GetPubsubConnection.2: Error %v", err))
	}

	return &p, err
}

func (p *PubsubConnection) publish(topicID string, msg []byte) error {

	t := p.client.Topic(topicID)
	result := t.Publish(p.ctx, &pubsub.Message{
		Data: msg,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(p.ctx)
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}
	fmt.Println(fmt.Sprintf("Published a message; msg ID: %v\n", id))

	return err
}

/*
func pullMsgs(projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	// Consume 10 messages.
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		msg.Ack()
		received++
		if received == 10 {
			cancel()
		}
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}
*/
