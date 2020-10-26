package persistenceServices

import (
	"errors"
	"fmt"
	"os"

	"github.com/adeturner/observability"
)

func GetPersistenceLayer(docType documentType) (*PersistenceLayer, error) {

	var err error

	p := PersistenceLayer{}

	p.docType = docType

	p.gcpProjectID = os.Getenv("GCP_PROJECT")
	if p.gcpProjectID == "" {
		return nil, errors.New("Error: GetPersistenceLayer GCP_PROJECT environment variable not set!")
	}

	p.cloudeventDomain = os.Getenv("CLOUDEVENT_DOMAIN")
	if p.cloudeventDomain == "" {
		return nil, errors.New("Error: GetPersistenceLayer CLOUDEVENT_DOMAIN environment variable not set!")
	}

	if os.Getenv("DEBUG") == "true" {
		p.debug = true
		observability.Logger("Info", fmt.Sprintf("GetPersistenceLayer : DEBUG on"))
	} else {
		p.debug = false
		observability.Logger("Info", fmt.Sprintf("GetPersistenceLayer : DEBUG off"))
	}

	if os.Getenv("USE_FIRESTORE") == "true" {
		p.useFirestore = true
		observability.Logger("Info", fmt.Sprintf("GetPersistenceLayer : USE_FIRESTORE on"))
	} else {
		p.useFirestore = false
	}

	if os.Getenv("USE_PUBSUB") == "true" {
		p.usePubsub = true
		observability.Logger("Info", fmt.Sprintf("GetPersistenceLayer : USE_PUBSUB on"))
	} else {
		p.usePubsub = false
	}

	if os.Getenv("USE_CQRS") == "true" {
		p.useCQRS = true
		observability.Logger("Info", fmt.Sprintf("GetPersistenceLayer : USE_CQRS on"))
	} else {
		p.useCQRS = false
	}

	if p.useFirestore {
		p.firestoreConnection, err = GetFirestoreConnection(p.gcpProjectID, docType.String())

		if err != nil {
			return nil, errors.New(fmt.Sprintf("GetPersistenceLayer.1: Error obtaining firestore connection: %v", err))
		}
	}

	if p.usePubsub {
		p.pubsubConnection, err = GetPubsubConnection(p.gcpProjectID)

		if err != nil {
			return nil, errors.New(fmt.Sprintf("GetPersistenceLayer.1: Error obtaining pubsub connection: %v", err))
		}
	}

	return &p, err
}

func (p *PersistenceLayer) CloudEventDomain() string {
	return p.cloudeventDomain
}
