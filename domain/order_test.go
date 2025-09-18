package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	userID := uuid.New()
	itemID := "intel-basic"
	quantity := 2

	orderAggregate := NewOrder(userID, itemID, quantity)

	assert.NotNil(t, orderAggregate)
	assert.NotNil(t, orderAggregate.Order)
	assert.Equal(t, userID, orderAggregate.Order.UserID)
	assert.Equal(t, itemID, orderAggregate.Order.ItemID)
	assert.Equal(t, quantity, orderAggregate.Order.Quantity)
	assert.Equal(t, OrderStatusPending, orderAggregate.Order.Status)
	assert.NotEmpty(t, orderAggregate.Order.ID)
	assert.False(t, orderAggregate.Order.CreatedAt.IsZero())
	assert.False(t, orderAggregate.Order.UpdatedAt.IsZero())
}

func TestOrderAggregate_Confirm(t *testing.T) {
	userID := uuid.New()
	orderAggregate := NewOrder(userID, "intel-basic", 1)
	originalUpdatedAt := orderAggregate.Order.UpdatedAt

	time.Sleep(1 * time.Millisecond)
	orderAggregate.Confirm()

	assert.Equal(t, OrderStatusConfirmed, orderAggregate.Order.Status)
	assert.True(t, orderAggregate.Order.UpdatedAt.After(originalUpdatedAt))
}

func TestOrderAggregate_Complete(t *testing.T) {
	userID := uuid.New()
	orderAggregate := NewOrder(userID, "intel-basic", 1)
	originalUpdatedAt := orderAggregate.Order.UpdatedAt

	time.Sleep(1 * time.Millisecond)
	orderAggregate.Complete()

	assert.Equal(t, OrderStatusCompleted, orderAggregate.Order.Status)
	assert.True(t, orderAggregate.Order.UpdatedAt.After(originalUpdatedAt))
}

func TestOrderAggregate_Cancel(t *testing.T) {
	userID := uuid.New()
	orderAggregate := NewOrder(userID, "intel-basic", 1)
	originalUpdatedAt := orderAggregate.Order.UpdatedAt

	time.Sleep(1 * time.Millisecond)
	orderAggregate.Cancel()

	assert.Equal(t, OrderStatusCancelled, orderAggregate.Order.Status)
	assert.True(t, orderAggregate.Order.UpdatedAt.After(originalUpdatedAt))
}