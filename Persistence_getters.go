package persistenceServices

import (
	"errors"
	"fmt"

	"github.com/adeturner/observability"
)

// GetLayer -
func GetLayer() (*PersistenceLayer, error) {
	var err error
	p := PersistenceLayer{}
	err = p.Initialise()
	if err == nil {
		err = p.GetConnection()
	}
	return &p, err
}

// GetPersistenceLayer - deprecated
func GetPersistenceLayer(docType documentType) (*PersistenceLayer, error) {

	observability.Logger("Debug", fmt.Sprintf("persistenceServices.GetPersistenceLayer is deprecated"))

	var err error
	p := PersistenceLayer{}
	p.docType = docType
	p.Initialise()

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
