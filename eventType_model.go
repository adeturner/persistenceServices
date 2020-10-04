package persistenceServices

type EventType int

const (
	EVENT_TYPE_CREATE EventType = iota
	EVENT_TYPE_UPDATE
	EVENT_TYPE_DELETE
	EVENT_TYPE_NOTIFICATION
	EVENT_TYPE_ERROR
)

func (t EventType) String() string {
	return [...]string{"Create", "Update", "Delete", "Notification"}[t]
}
