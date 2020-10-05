

export GCP_PROJECT=myproject
export CLOUDEVENT_DOMAIN=mydomain.com
export DEBUG=true
export USE_FIRESTORE=true
export USE_PUBSUB=true
export USE_CQRS=false
export GOOGLE_APPLICATION_CREDENTIALS=~/secrets/apiv01-persistenceServices.json

// create a service account that has privileges on pubsub and firestore
// and create the secrets file
gcloud iam service-accounts keys create ${GOOGLE_APPLICATION_CREDENTIALS} --iam-account myserviceaccount@apiv01.iam.gserviceaccount.com

go run cmd/main.go

