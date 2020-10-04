package persistenceServices

import (
	"context"

	"cloud.google.com/go/firestore"
)

type FirestoreConnection struct {
	ctx        context.Context
	client     firestore.Client
	collection firestore.CollectionRef
}
