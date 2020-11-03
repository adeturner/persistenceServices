package persistenceServices

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/adeturner/observability"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// AddDocument -
func (p *PersistenceLayer) AddDocument(key string, values interface{}) error {

	var err error

	if p.useFirestore && !p.useCQRS {
		// direct database write
		err = p.firestoreConnection.FirestoreAdd(key, values)
	}
	if err == nil && p.usePubsub && p.useCQRS {
		err = p.Publish(EVENT_TYPE_CREATE, key, values)
	}

	return err
}

// UpdateDocument -
func (p *PersistenceLayer) UpdateDocument(key string, values interface{}) error {

	var err error

	if p.useFirestore && !p.useCQRS {
		// direct database write
		err = p.firestoreConnection.FirestoreUpdate(key, values)
	}
	if err == nil && p.usePubsub && p.useCQRS {
		err = p.Publish(EVENT_TYPE_UPDATE, key, values)
	}

	return err
}

// DeleteDocument -
func (p *PersistenceLayer) DeleteDocument(key string, values interface{}) error {

	var err error

	if p.useFirestore && !p.useCQRS {
		// direct database write
		err = p.firestoreConnection.FirestoreDelete(key, values)
	}
	if err == nil && p.usePubsub && p.useCQRS {
		err = p.Publish(EVENT_TYPE_DELETE, key, values)
	}

	return err

}

// Publish -
func (p *PersistenceLayer) Publish(eventType EventType, subject string, values interface{}) error {

	var cloudevent []byte

	payload, err := json.Marshal(values)

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Error Failed to marshall to json topic=Source error=%v", err))

	} else {

		observability.Logger("Info", fmt.Sprintf("Publishing Payload=%v", string(payload)))

		event := cloudevents.NewEvent()
		event.SetSource(p.cloudeventDomain + "/" + p.docType.String())
		event.SetSubject(subject)
		event.SetType(eventType.String())
		event.SetID(observability.GetCorrId())
		event.SetExtension("CausationId", observability.GetCausationId())
		event.SetData(cloudevents.ApplicationJSON, payload)

		cloudevent, err = json.Marshal(event)

		observability.Logger("Info", fmt.Sprintf("Publishing topic=Source bytes=%v", string(cloudevent)))

	}

	if err == nil && p.usePubsub {
		err = p.pubsubConnection.publish(p.docType.Topic(), cloudevent)
		if err != nil {
			observability.Logger("Error", fmt.Sprintf("Error: failed to publish to topic=Source error=%v", err))
		}
	}

	return err
}

// FindById -
func (p *PersistenceLayer) FindById(key string, values interface{}) (interface{}, error) {

	var err error
	var v interface{}

	if p.useFirestore {
		v, err = p.firestoreConnection.FirestoreFindById(key, values)
	}

	return v, err
}

// Find -
func (p *PersistenceLayer) Find(queryParams map[string][]string, value interface{}) (valuesArray interface{}, err error) {

	if p.useFirestore {
		observability.Logger("Info", fmt.Sprintf("firestore find starting"))
		valuesArray, err = p.firestoreConnection.FirestoreFind(queryParams, value)
	}

	return valuesArray, err

}

// Initialise -
func (p *PersistenceLayer) Initialise() error {

	observability.Logger("Info", fmt.Sprintf("Initialising PersistenceLayer"))

	p.gcpProjectID = os.Getenv("GCP_PROJECT")
	if p.gcpProjectID == "" {
		return errors.New("Error: GetPersistenceLayer GCP_PROJECT environment variable not set!")
	}

	p.cloudeventDomain = os.Getenv("CLOUDEVENT_DOMAIN")
	if p.cloudeventDomain == "" {
		return errors.New("Error: GetPersistenceLayer CLOUDEVENT_DOMAIN environment variable not set!")
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

	return nil
}

// GetConnection -
func (p *PersistenceLayer) GetConnection() (err error) {

	if p.useFirestore {
		p.firestoreConnection, err = GetFirestore(p.gcpProjectID)

		if err != nil {
			return errors.New(fmt.Sprintf("Error obtaining firestore connection: %v", err))
		}
	}

	if p.usePubsub {
		p.pubsubConnection, err = GetPubsubConnection(p.gcpProjectID)

		if err != nil {
			return errors.New(fmt.Sprintf("Error obtaining pubsub connection: %v", err))
		}
	}

	return nil
}

// SetDocType - Note that this only impacts Firestore
func (p *PersistenceLayer) SetDocType(docType documentType) {

	// Note that this only impacts Firestore
	if p.useFirestore {
		s := docType.String()
		observability.Logger("Info", fmt.Sprintf("Setting collection to %s", s))
		p.firestoreConnection.SetCollection(s)
	}

}
