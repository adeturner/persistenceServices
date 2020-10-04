package persistenceServices

func LocalEntry(docType documentType) (*PersistenceLayer, error) {
	p, err := GetPersistenceLayer(docType)
	return p, err
}
