package persistenceServices

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type PubsubConnection struct {
	ctx    context.Context
	client pubsub.Client
}
