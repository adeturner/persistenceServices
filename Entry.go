package persistenceServices

func LocalEntry(docType documentType) (*PersistenceLayer, error) {

	// deprecated
	p, err := GetPersistenceLayer(docType)

	// new
	p, err = GetLayer()
	p.SetDocType(docType)

	return p, err
}
