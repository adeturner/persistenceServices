package persistenceServices

import (
	"encoding/json"
	"fmt"

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
func (p *PersistenceLayer) Publish(eventType EventType, key string, values interface{}) error {

	var cloudevent []byte

	payload, err := json.Marshal(values)

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Error Failed to marshall to json topic=Source error=%v", err))

	} else {

		observability.Logger("Info", fmt.Sprintf("Publishing Payload=%v", string(payload)))

		event := cloudevents.NewEvent()
		event.SetSource(p.cloudeventDomain + "/" + p.docType.String())
		event.SetType(eventType.String())
		event.SetID(observability.GetCorrId())
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
