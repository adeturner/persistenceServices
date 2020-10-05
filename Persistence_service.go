package persistenceServices

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
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
		fmt.Println(fmt.Sprintf("GetPersistenceLayer : DEBUG on"))
	} else {
		p.debug = false
		fmt.Println(fmt.Sprintf("GetPersistenceLayer : DEBUG off"))
	}

	if os.Getenv("USE_FIRESTORE") == "true" {
		p.useFirestore = true
		fmt.Println(fmt.Sprintf("GetPersistenceLayer : USE_FIRESTORE on"))
	} else {
		p.useFirestore = false
	}

	if os.Getenv("USE_PUBSUB") == "true" {
		p.usePubsub = true
		fmt.Println(fmt.Sprintf("GetPersistenceLayer : USE_PUBSUB on"))
	} else {
		p.usePubsub = false
	}

	if os.Getenv("USE_CQRS") == "true" {
		p.useCQRS = true
		fmt.Println(fmt.Sprintf("GetPersistenceLayer : USE_CQRS on"))
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

// Publish to Pubsub
func (p *PersistenceLayer) Publish(eventType EventType, key string, values interface{}) error {

	// Producers MUST ensure that source + id is unique for each distinct event.
	// Consumers MAY assume that Events with identical source and id are duplicates.

	var cloudevent []byte

	payload, err := json.Marshal(values)

	if err != nil {
		fmt.Println(fmt.Sprintf("Source.PublishPubsub : Error Failed to marshall to json topic=Source error=%v", err))

	} else {

		fmt.Println(fmt.Sprintf("Source.PublishPubsub : Publishing Payload=%v", string(payload)))

		event := cloudevents.NewEvent()
		event.SetSource(p.cloudeventDomain + "/" + p.docType.String())
		event.SetType(eventType.String())
		event.SetID(key)
		// map[string]string{"hello": "world"}
		event.SetData(cloudevents.ApplicationJSON, payload)

		cloudevent, err = json.Marshal(event)

		fmt.Println(fmt.Sprintf("Publishing topic=Source bytes=%v", string(cloudevent)))

	}

	if err == nil {
		err = p.pubsubConnection.publish(p.docType.Topic(), cloudevent)
		if err != nil {
			fmt.Println(fmt.Sprintf("Source.PublishPubsub : Error: failed to publish to topic=Source error=%v", err))
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

// FindByTags -
func (p *PersistenceLayer) FindByTags(tags []string, strlimit string, value interface{}, valuesArray interface{}) (interface{}, error) {

	var err error

	if p.useFirestore {
		fmt.Println(fmt.Sprintf("Source.FindByTags : firestore find starting"))
		//valuesArray, err = p.firestoreConnection.FirestoreFindByTags(tags, strlimit, value, valuesArray)
		valuesArray, err = p.firestoreConnection.FirestoreFindByTags(tags, strlimit, value)
	}

	return valuesArray, err

}
