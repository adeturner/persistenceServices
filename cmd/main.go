package main

import "github.com/adeturner/persistenceServices"

type DocType int

func (d DocType) String() string {
	return [...]string{
		"FirstDocType",
		"SecondDocType",
	}[d]
}

func (d DocType) Topic() string {
	return [...]string{
		"FirstDocTypeTopic",
		"SecondDocTypeTopic",
	}[d]
}

func main() {

	const (
		DOCUMENT_TYPE_FIRST DocType = iota
		DOCUMENT_TYPE_SECOND
	)

	persistenceServices.LocalEntry(DOCUMENT_TYPE_FIRST)
}
