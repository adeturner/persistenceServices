package persistenceServices

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/adeturner/observability"
)

// GetFirestoreConnection - deprecated
func GetFirestoreConnection(gcpProjectID string, collectionStr string) (*FirestoreConnection, error) {

	observability.Logger("Debug", fmt.Sprintf("firestore_service.GetFirestoreConnection is deprecated"))

	f := FirestoreConnection{}

	if gcpProjectID == "" {
		return &f, errors.New("Error: GCP_PROJECT environment variable not set!")
	}

	observability.Logger("Debug", fmt.Sprintf("getting context"))
	f.ctx = context.Background()

	observability.Logger("Debug", fmt.Sprintf("have context, new client"))

	client, err := firestore.NewClient(f.ctx, gcpProjectID)

	observability.Logger("Debug", fmt.Sprintf("new client complete"))

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Error %v", err))
		return &f, nil

	} else {

		f.client = *client

		collection := client.Collection(collectionStr)
		f.collection = *collection
		observability.Logger("Info", fmt.Sprintf("success"))
	}

	return &f, err
}

// GetFirestore -
func GetFirestore(gcpProjectID string) (*FirestoreConnection, error) {

	observability.Logger("Debug", fmt.Sprintf("starting"))

	f := FirestoreConnection{}

	if gcpProjectID == "" {
		return &f, errors.New("Error: GCP_PROJECT environment variable not set!")
	}

	observability.Logger("Debug", fmt.Sprintf("getting context"))
	f.ctx = context.Background()

	observability.Logger("Debug", fmt.Sprintf("have context, new client"))

	client, err := firestore.NewClient(f.ctx, gcpProjectID)

	observability.Logger("Debug", fmt.Sprintf("new client complete"))

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Error %v", err))
		return &f, nil

	} else {
		f.client = *client
		observability.Logger("Info", fmt.Sprintf("success"))
	}

	return &f, err
}

func (f *FirestoreConnection) SetCollection(collectionStr string) {
	c := f.client.Collection(collectionStr)
	f.collection = *c
}

// FirestoreAdd -
func (f *FirestoreConnection) FirestoreAdd(docId string, s interface{}) error {

	observability.Logger("Debug", fmt.Sprintf("Starting Firestore write %v", s))

	docRef := f.collection.Doc(docId)

	// wr is a WriteResult, which contains the time at which the document was updated.
	wr, err := docRef.Create(f.ctx, s)

	if err != nil {
		// e.g. it already exists
		observability.Logger("Error", fmt.Sprintf("Error %v", err))
	} else {
		observability.Logger("Info", fmt.Sprintf("Success. %v created at %s", wr, s))
	}

	return err
}

// FirestoreUpdate -
func (f *FirestoreConnection) FirestoreUpdate(docId string, s interface{}) error {

	observability.Logger("Debug", fmt.Sprintf("Starting Firestore update  %v", s))

	docRef := f.collection.Doc(docId)

	// wr is a WriteResult, which contains the time at which the document was updated.
	wr, err := docRef.Set(f.ctx, s)

	if err != nil {
		// e.g. it already exists
		observability.Logger("Error", fmt.Sprintf("Error %v", err))
	} else {
		observability.Logger("Info", fmt.Sprintf("Success. %v updated %s", wr, s))
	}

	return err
}

// FirestoreDelete -
func (f *FirestoreConnection) FirestoreDelete(docId string, s interface{}) error {

	observability.Logger("Debug", fmt.Sprintf("Starting Firestore delete id=%v", docId))

	docRef := f.collection.Doc(docId)

	// Note: success even if it doesnt exist
	_, err := docRef.Delete(f.ctx)

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Error deleting source id=%s %v", docId, err))
	} else {
		observability.Logger("Info", fmt.Sprintf("Success. deleted id=%s", docId))
	}

	return err
}

// FirestoreFindById -
func (f *FirestoreConnection) FirestoreFindById(key string, values interface{}) (interface{}, error) {

	observability.Logger("Debug", fmt.Sprintf("About to find id=%s", key))

	docRef := f.collection.Doc(key)
	docsnap, err := docRef.Get(f.ctx)

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Error finding id=%s %v", key, err))
	}

	var dataMap map[string]interface{}

	if err == nil {
		err = docsnap.DataTo(&values)
		dataMap = docsnap.Data()
	}

	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Error deserialising into struct %v", dataMap))
	} else {
		observability.Logger("Info", fmt.Sprintf("Found values=%v", values))
	}

	return values, err
}

/*
FirestoreFind -
queryParams -  map["db.field"][]values
value = pass in a struct, e.g. Source{}
loop through the map building a query like
db.field1 in ("value[0]", "value[1]", .. "value[n]") && db.field2 in ("value[0]", "value[1]", .. "value[n]")
Limitation: create a separate query for each OR condition and merge the query results in your app.
*/
func (f *FirestoreConnection) FirestoreFind(queryParams map[string][]string, value interface{}) (interface{}, error) {

	var err error
	var docSnaps []*firestore.DocumentSnapshot

	caller := observability.Caller{}
	c := caller.Get(4)

	// 5 row default limit
	limit := 5

	var q firestore.Query
	var orderBy []string
	limitSet := false
	var queryStr string

	for key, element := range queryParams {
		observability.Logger("Debug", fmt.Sprintf("%s Key: %v => Element: %v", c, key, element))

		if key == "orderBy" {
			orderBy = element
		} else if key == "limit" {
			limitSet = true
			limit, err = strconv.Atoi(element[0])
			if err != nil {
				observability.Logger("Error", fmt.Sprintf("Error converting limit to int"))
			}
		} else {
			q = f.collection.Where(key, "in", element)
			queryStr = queryStr + fmt.Sprintf("WHERE %s in '%v' ", key, element)
		}
	}

	// add limit and orderby
	if err == nil {
		if len(orderBy) > 0 {
			for i := 0; i < len(orderBy); i++ {
				q = q.OrderBy(orderBy[i], firestore.Asc)
			}
		}

		if limitSet {
			q = q.Limit(limit)
		}
	}

	if err == nil {
		observability.Logger("Debug", fmt.Sprintf("About to find Sources matching orderBy=%v limit=%d", orderBy, limit))

		iter := q.Documents(f.ctx)
		docSnaps, err = iter.GetAll()
		if err != nil {
			observability.Logger("Error", fmt.Sprintf("Error getting documents %v", err))
		} else {
			observability.Logger("Info", fmt.Sprintf("Found %d documents", len(docSnaps)))
		}
	}

	var vArray []interface{}

	if err == nil {
		cnt := 0
		for _, ds := range docSnaps {
			err = ds.DataTo(&value)

			if err != nil {
				observability.Logger("Error", fmt.Sprintf("Failed to unmarshal %v", &value))
			}

			observability.Logger("Debug", fmt.Sprintf("Adding value=%v to array", value))

			vArray = append(vArray, value)
			cnt++
			if cnt >= limit {
				break
			}
		}
	}

	observability.Logger("Info", fmt.Sprintf("Returning vArray of length %d to caller", len(vArray)))

	return vArray, err
}
