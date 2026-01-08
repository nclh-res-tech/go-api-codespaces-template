package stores

import (
	"context"
	"sync"
	"time"

	"{{MODULE_PATH}}/common/errors"
	"{{MODULE_PATH}}/internal/models"

	"github.com/google/uuid"
)

// ErrNotFound is returned when an item is not found.
const ErrNotFound = errors.Error("item not found")

// ItemStore provides data access for items.
// This is an in-memory implementation. Replace with your database implementation.
type ItemStore struct {
	mu    sync.RWMutex
	items map[string]models.Item
}

// NewItemStore creates a new ItemStore.
func NewItemStore() *ItemStore {
	return &ItemStore{
		items: make(map[string]models.Item),
	}
}

// List returns all items.
func (s *ItemStore) List(ctx context.Context) ([]models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Item, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, item)
	}
	return result, nil
}

// Get retrieves an item by ID.
func (s *ItemStore) Get(ctx context.Context, id string) (*models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &item, nil
}

// Create creates a new item.
func (s *ItemStore) Create(ctx context.Context, req models.CreateItemRequest) (*models.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	item := models.Item{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.items[item.ID] = item
	return &item, nil
}

// Update updates an existing item.
func (s *ItemStore) Update(ctx context.Context, id string, req models.UpdateItemRequest) (*models.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return nil, ErrNotFound
	}

	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	item.UpdatedAt = time.Now()

	s.items[id] = item
	return &item, nil
}

// Delete removes an item.
func (s *ItemStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
