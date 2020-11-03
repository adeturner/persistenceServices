# persistenceServices

Sharing out of good will and generally unsupported

... I'm still debating its usefulness so please raise an issue if you wish to contribute and we'll discuss

Thanks, Adrian

## Introduction

This module implements a persistence layer, with optional CQRS

- Google Firestore
- Google Pubsub
- [Future] Azure Cosmodb
- [Future] Azure Event Hub
- etc

## Installation
```
go get github.com/adeturner/persistenceServices
go test
```

## Implementation

### Code

Define a docType, must implement documentType:

```go
type documentType interface {
	String() string
	Topic() string
}

type DocType int

// example doc types
const (
	DOCUMENT_TYPE_USERS DocType = iota
	DOCUMENT_TYPE_ORDERS
	DOCUMENT_TYPE_STOCK
)

// Generic name for the Document, used as the Firestore Collection name
func (d DocType) String() string {
	return [...]string{
		"Users",
		"Orders",
		"Stock",
	}[d]
}

// Returns the eventing topic name; topic must be precreated
func (d DocType) Topic() string {
	return [...]string{
		"usersTopic",
		"ordersTopic",
		"stockTopic",
	}[d]
}
```

Instantiate the persistenceLayer

```go
var persistenceLayer *persistenceServices.PersistenceLayer
ver err error
d := DOCUMENT_TYPE_USERS

// deprecated
p, err := GetPersistenceLayer(docType)

// new
p, err = GetLayer()
p.SetDocType(docType)
```

Available functions
```go
func GetPersistenceLayer(docType documentType) (*PersistenceLayer, error) {
func (p *PersistenceLayer) AddDocument(key string, values interface{}) error {
func (p *PersistenceLayer) UpdateDocument(key string, values interface{}) error {
func (p *PersistenceLayer) DeleteDocument(key string, values interface{}) error {
func (p *PersistenceLayer) FindById(key string, values interface{}) (interface{}, error) {
func (p *PersistenceLayer) FindByTags(tags []string, strlimit string, value interface{}, valuesArray interface{}) (interface{}, error) {
```

Example usage

```go
UUID := uuid.New().String()
u := User{Id: UUID, Name: "SomeName", Tag: "SomeTags"}
err := p.AddDocument(UUID, u)
```

### Setup

Environment variables control the data layer access

```go
export CLOUDEVENT_DOMAIN=mydomain.com // Cloud events will have source {CLOUDEVENT_DOMAIN}/{docType.String}
export DEBUG=false            // if true, debug output on
export USE_FIRESTORE=true     // if true reads will come from here; writes will also go here if USE_CQRS=false
export USE_PUBSUB=true        // if true enables a pubsub connection
export USE_CQRS=false         // if true writes are sent to pubsub
export GCP_PROJECT=myproject  // project for GCP connections
export GOOGLE_APPLICATION_CREDENTIALS=~/secrets/persistenceServices.json  // if testing outside of GCP
```

### Tests

```go
./test/localtests.sh
```








