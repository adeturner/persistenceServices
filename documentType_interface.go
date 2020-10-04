package persistenceServices

type documentType interface {
	String() string
	Topic() string
}
