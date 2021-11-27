package subscription

import "fmt"

// Service - экземпляр сервиса.
//
// ID: идентификатор.
//
// Name: наименование.
//
// Description: описание.
//
type Service struct {
	ID          uint64
	Name        string
	Description string
}

func (s *Service) String() string {
	return fmt.Sprintf("%v. %v: %v", s.ID, s.Name, s.Description)
}

// EventType - тип события экземпляра сервиса.
type EventType string

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

// ServiceListChunkGetter - отложенный ленивый загрузчик чанка списка.
type ServiceListChunkGetter func() (*ServiceListChunk, error)

// ServiceListChunk - чанк списка сервисов.
type ServiceListChunk struct {
	Services         []*Service
	PreviousServices ServiceListChunkGetter
	NextServices     ServiceListChunkGetter
}

// Len - количество элементов в списке.
func (c *ServiceListChunk) Len() int {
	return len(c.Services)
}

// IsHasPreviousPage - имеет ли список предыдущую страницу.
func (c *ServiceListChunk) IsHasPreviousPage() bool {
	return c.PreviousServices != nil
}

// IsHasNextPage - имеет ли список следующую страницу.
func (c *ServiceListChunk) IsHasNextPage() bool {
	return c.NextServices != nil
}
