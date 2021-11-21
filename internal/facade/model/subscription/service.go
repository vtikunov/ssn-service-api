package subscription

// Service - экземпляр сервиса.
//
// ID: идентификатор.
//
// Name: наименование.
//
// Description: описание.
//
// IsRemoved: флаг удаленного сервиса.
//
type Service struct {
	ID          uint64 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	IsRemoved   bool   `db:"is_removed"`
	LastEventID uint64 `db:"last_event_id"`
}

// EventType - тип события экземпляра сервиса.
type EventType string

// EventSubType - субтип события экземпляра сервиса.
type EventSubType string

// Типы событий экземпляра сервиса
//
// Created: сервис создан.
//
// Updated: сервис обновлен.
//
// Removed: сервис удален.
const (
	Created EventType = "CREATED"
	Updated EventType = "UPDATED"
	Removed EventType = "REMOVED"
)

// Субтипы событий экземпляра сервиса
//
// NoneSubType: не определен.
//
// NameSubtype: событие для свойства Name.
//
// DescriptionSubType: событие для свойства Description.
const (
	NoneSubType        EventSubType = "NONE"
	NameSubtype        EventSubType = "NAME"
	DescriptionSubType EventSubType = "DESCRIPTION"
)

// ServiceEvent - событие экземпляра сервиса.
//
// ID: идентификатор события.
//
// ServiceID: идентификатор сервиса.
//
// Type: тип события (EventType).
//
// SubType: подтип события (EventSubType).
//
// Service: экземпляр сервиса (Service).
type ServiceEvent struct {
	ID        uint64
	ServiceID uint64
	Type      EventType
	SubType   EventSubType
	Service   *Service
}
