package subscription

import "time"

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
	ID          uint64    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	IsRemoved   bool      `db:"is_removed"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// EventType - тип события экземпляра сервиса.
type EventType string

// EventStatus - статус события экземпляра сервиса.
type EventStatus string

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

// Статусы событий экземпляра сервиса
//
// Deferred: событие отложено (ожидает обработки).
//
// Processed: событие обрабатывается.
const (
	Deferred  EventStatus = "DEFERRED"
	Processed EventStatus = "PROCESSED"
)

// ServiceEvent - событие экземпляра сервиса.
//
// ID: идентификатор события.
//
// ServiceID: идентификатор сервиса.
//
// Type: тип события (EventType).
//
// Status: статус события (EventStatus).
//
// Service: экземпляр сервиса (Service).
type ServiceEvent struct {
	ID        uint64      `db:"id"`
	ServiceID uint64      `db:"service_id"`
	Type      EventType   `db:"type"`
	Status    EventStatus `db:"status"`
	Service   *Service    `db:"payload"`
	UpdatedAt time.Time   `db:"updated_at"`
}
