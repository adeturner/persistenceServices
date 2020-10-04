package persistenceServices

// PersistenceLayer implements hexagonal architecture to hide the PaaS functionality from the APIs
type PersistenceLayer struct {
	docType             documentType
	firestoreConnection *FirestoreConnection
	pubsubConnection    *PubsubConnection

	// ENVIRONMENT VARIABLES
	// The following are initialised when the persistenceLayer is created
	gcpProjectID     string // export GCP_PROJECT=myproject
	cloudeventDomain string // export CLOUDEVENT_DOMAIN=mydomain.com (for use in cloudevent headers)
	debug            bool   // export DEBUG=true
	useFirestore     bool   // export USE_FIRESTORE=true
	usePubsub        bool   // export USE_PUBSUB=true
	useCQRS          bool   // export USE_CQRS=true
}
