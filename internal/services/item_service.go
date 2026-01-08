package services

import (
	"context"

	"{{MODULE_PATH}}/internal/models"
	"{{MODULE_PATH}}/internal/stores"
)

// ItemService provides business logic for items.
type ItemService struct {
	store *stores.ItemStore
}

// NewItemService creates a new ItemService.
func NewItemService(store *stores.ItemStore) *ItemService {
	return &ItemService{store: store}
}

// List returns all items.
func (s *ItemService) List(ctx context.Context) ([]models.Item, error) {
	return s.store.List(ctx)
}

// Get retrieves an item by ID.
func (s *ItemService) Get(ctx context.Context, id string) (*models.Item, error) {
	return s.store.Get(ctx, id)
}

// Create creates a new item.
func (s *ItemService) Create(ctx context.Context, req models.CreateItemRequest) (*models.Item, error) {
	return s.store.Create(ctx, req)
}

// Update updates an existing item.
func (s *ItemService) Update(ctx context.Context, id string, req models.UpdateItemRequest) (*models.Item, error) {
	return s.store.Update(ctx, id, req)
}

// Delete removes an item.
func (s *ItemService) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
