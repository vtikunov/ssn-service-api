package subscription

import "time"

// Service - экземпляр сервиса.
//
// ID: идентификатор.
//
// Name: наименование.
type Service struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	IsRemoved bool      `db:"is_removed"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
	ID        uint64
	ServiceID uint64
	Type      EventType
	Status    EventStatus
	Service   *Service
}
