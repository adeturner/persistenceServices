package persistenceServices

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
)

// GetFirestoreConnection -
func GetFirestoreConnection(gcpProjectID string, collectionStr string) (*FirestoreConnection, error) {

	if debug {
		fmt.Println(fmt.Sprintf("GetFirestoreConnection: starting"))
	}

	f := FirestoreConnection{}

	if debug {
		fmt.Println(fmt.Sprintf("GetFirestoreConnection: set variable"))
	}

	if gcpProjectID == "" {
		return &f, errors.New("Error: GetFirestoreConnection GCP_PROJECT environment variable not set!")
	}

	if debug {
		fmt.Println(fmt.Sprintf("GetFirestoreConnection: getting context"))
	}
	f.ctx = context.Background()

	if debug {
		fmt.Println(fmt.Sprintf("GetFirestoreConnection: have context, new client"))
	}
	client, err := firestore.NewClient(f.ctx, gcpProjectID)
	if debug {
		fmt.Println(fmt.Sprintf("GetFirestoreConnection: new client complete"))
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("GetFirestoreConnection.1: Error %v", err))
		return &f, nil

	} else {

		f.client = *client

		collection := client.Collection(collectionStr)
		f.collection = *collection
		if debug {
			fmt.Println(fmt.Sprintf("GetFirestoreConnection: success"))
		}
	}

	return &f, err
}

// FirestoreAdd -
func (f *FirestoreConnection) FirestoreAdd(docId string, s interface{}) error {

	if debug {
		fmt.Println(fmt.Sprintf("FirestoreAdd.AddSource.1: Starting Firestore write %v", s))
	}

	docRef := f.collection.Doc(docId)

	// wr is a WriteResult, which contains the time at which the document was updated.
	wr, err := docRef.Create(f.ctx, s)

	if err != nil {
		// e.g. it already exists
		fmt.Println(fmt.Sprintf("SourcesapiService.AddSource.2: Error %v", err))
	} else {
		if debug {
			fmt.Println(fmt.Sprintf("SourcesapiService.AddSource.3: Success. %v created at %s", wr, s))
		}
	}

	return err
}

// FirestoreUpdate -
func (f *FirestoreConnection) FirestoreUpdate(docId string, s interface{}) error {

	if debug {
		fmt.Println(fmt.Sprintf("FirestoreUpdate.1: Starting Firestore update  %v", s))
	}

	docRef := f.collection.Doc(docId)

	// wr is a WriteResult, which contains the time at which the document was updated.
	wr, err := docRef.Set(f.ctx, s)

	if err != nil {
		// e.g. it already exists
		fmt.Println(fmt.Sprintf("FirestoreUpdate.2: Error %v", err))
	} else {
		if debug {
			fmt.Println(fmt.Sprintf("FirestoreUpdate.3: Successrv. %v updated %s", wr, s))
		}
	}

	return err
}

// FirestoreDelete -
func (f *FirestoreConnection) FirestoreDelete(docId string, s interface{}) error {

	if debug {
		fmt.Println(fmt.Sprintf("FirestoreDelete.1: Starting Firestore delete id=%v", docId))
	}

	docRef := f.collection.Doc(docId)

	// Note: success even if it doesnt exist
	_, err := docRef.Delete(f.ctx)

	if err != nil {
		fmt.Println(fmt.Sprintf("FirestoreDelete.2: Error deleting source id=%s %v", docId, err))
	} else {
		if debug {
			fmt.Println(fmt.Sprintf("FirestoreDelete.3: Success. deleted id=%s", docId))
		}
	}

	return err
}

// FirestoreFindById -
func (f *FirestoreConnection) FirestoreFindById(key string, values interface{}) (interface{}, error) {

	if debug {
		fmt.Println(fmt.Sprintf("FirestoreConnection.findById.1: About to find id=%s", key))
	}

	docRef := f.collection.Doc(key)
	docsnap, err := docRef.Get(f.ctx)

	if err != nil {
		fmt.Println(fmt.Sprintf("PersistenceLayer.findById.2: Error finding id=%s %v", key, err))
	}

	var dataMap map[string]interface{}

	if err == nil {
		err = docsnap.DataTo(&values)
		dataMap = docsnap.Data()
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("PersistenceLayer.findById.3 Error deserialising into struct %v", dataMap))
	}

	return values, err
}

/*
FirestoreFindByTags -
pass in e.g.
tags = list of values to search for
strlimit = integer limit of the number of rows after which to stop
value = pass in the type, e.g. Source{}
valuesArray, pass in an array of the type, e.g. []Source{}
*/
//func (f *FirestoreConnection) FirestoreFindByTags(tags []string, strlimit string, value interface{}, valuesArray interface{}) (interface{}, error) {
func (f *FirestoreConnection) FirestoreFindByTags(tags []string, strlimit string, value interface{}) (interface{}, error) {

	var err error
	var limit int
	var docSnaps []*firestore.DocumentSnapshot

	/*
		vArray, ok := valuesArray.([]interface{})

		if !ok {
			fmt.Println(fmt.Sprintf("FirestoreConnection.findbyTags.1 Failed to cast valuesArray to array type"))
			err = errors.New("FirestoreConnection.findbyTags.1 Failed to cast valuesArray to array type")

		}
	*/

	limit, err = strconv.Atoi(strlimit)
	if err != nil {
		fmt.Println(fmt.Sprintf("FirestoreConnection.findbyTags.1 Error converting limit to int"))
	}

	if debug {
		fmt.Println(fmt.Sprintf("tags length %d %s", len(tags), tags[0]))
	}

	if err == nil {

		if debug {
			fmt.Println(fmt.Sprintf("FirestoreConnection.findbyTags.2: About to find Sources matching tags=%v limit=%d", tags, limit))
		}

		q := f.collection.Where("Tag", "in", tags).OrderBy("Name", firestore.Asc)

		if len(tags) == 1 && tags[0] == "" {
			if debug {
				fmt.Println(fmt.Sprintf("tags length %d", len(tags)))
			}
			q = f.collection.OrderBy("Name", firestore.Asc)
		}

		iter := q.Documents(f.ctx)
		docSnaps, err = iter.GetAll()
		if err != nil {
			fmt.Println(fmt.Sprintf("FirestoreConnection.findbyTags.4 Error getting documents %v", err))
		} else {
			if debug {
				fmt.Println(fmt.Sprintf("FirestoreConnection.findbyTags.3: Found %d documents", len(docSnaps)))
			}
		}
	}

	var vArray []interface{}

	if err == nil {
		cnt := 0
		for _, ds := range docSnaps {
			err = ds.DataTo(&value)

			if debug {
				//fmt.Println(fmt.Sprintf("FirestoreConnection.findbyTags.5: Adding value=%v to array", value))
			}

			vArray = append(vArray, value)
			cnt++
			if cnt >= limit {
				break
			}
		}
	}

	if debug {
		fmt.Println(fmt.Sprintf("FirestoreConnection.findbyTags.5: Returning vArray of length %d to caller", len(vArray)))
	}

	return vArray, err

}
