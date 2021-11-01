package subscription

// Service - экземпляр сервиса.
//
// ID: идентификатор.
//
// Name: наименование.
type Service struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
}

// EventType - тип события экземпляра сервиса.
type EventType uint8

// EventStatus - статус события экземпляра сервиса.
type EventStatus uint8

// Типы событий экземпляра сервиса
//
// Created: сервис создан.
//
// Updated: сервис обновлен.
//
// Removed: сервис удален.
const (
	Created EventType = iota
	Updated
	Removed
)

// Статусы событий экземпляра сервиса
//
// Deferred: событие отложено (ожидает обработки).
//
// Processed: событие обрабатывается.
const (
	Deferred EventStatus = iota
	Processed
)

// ServiceEvent - событие экземпляра сервиса.
//
// ID: идентификатор события.
//
// Type: тип события (EventType).
//
// Status: статус события (EventStatus).
//
// Service: экземпляр сервиса (Service).
type ServiceEvent struct {
	ID      uint64
	Type    EventType
	Status  EventStatus
	Service *Service
}
