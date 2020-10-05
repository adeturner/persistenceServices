# persistenceServices

## Introduction

Implements a persistence layer

- Google Firestore
- Google Pubsub
- [Future] Azure Cosmodb
- [Future] Azure Event Hub

```
go get github.com/adeturner/persistenceServices
```

## Implementation

### Code

Define a docType, must implement documentType:

```
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

func (d DocType) String() string {
	return [...]string{
		"Users",
		"Orders",
		"Stock",
	}[d]
}

func (d DocType) Topic() string {
	return [...]string{
		"usersTopic",
		"ordersTopic",
		"stockTopic",
	}[d]
}
```

Instantiate the persistenceLayer

```
persistenceLayer *persistenceServices.PersistenceLayer

persistenceLayer, err := persistenceServices.GetPersistenceLayer(docType)
```

Available functions
```
func GetPersistenceLayer(docType documentType) (*PersistenceLayer, error) {
func (p *PersistenceLayer) AddDocument(key string, values interface{}) error {
func (p *PersistenceLayer) UpdateDocument(key string, values interface{}) error {
func (p *PersistenceLayer) DeleteDocument(key string, values interface{}) error {
func (p *PersistenceLayer) FindById(key string, values interface{}) (interface{}, error) {
func (p *PersistenceLayer) FindByTags(tags []string, strlimit string, value interface{}, valuesArray interface{}) (interface{}, error) {
```

Example usage

```
UUID := uuid.New().String()
u := User{Id: UUID, Name: "SomeName", Tag: "SomeTags"}
err := p.AddDocument(UUID, u)
```

### Setup

Environment variables control the data layer access

```
export CLOUDEVENT_DOMAIN=mydomain.com // Cloud events will have source {CLOUDEVENT_DOMAIN}/{docType.String}
export DEBUG=false            // if true, debug output on
export USE_FIRESTORE=true     // if true reads will come from here
export USE_PUBSUB=true        // if true writes will be sent to the topic
export USE_CQRS=false         // if true writes are sent to the topic
export GCP_PROJECT=myproject  // project for GCP connections
export GOOGLE_APPLICATION_CREDENTIALS=~/secrets/persistenceServices.json  // if testing outside of GCP
```

### Tests

```
./test/localtests.sh
```








