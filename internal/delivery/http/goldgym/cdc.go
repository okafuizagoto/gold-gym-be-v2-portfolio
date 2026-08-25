package goldgym

import (
	"context"
	registry "gold-gym-be/internal/registry"
	"log"

	"gold-gym-be/internal/resources"
)

// CDCHandler bertugas menghubungkan event Kafka → handler function
type CDCHandler struct {
	res      *resources.BootResources
	registry map[string]registry.HandlerFunc
}

// NewCDCHandler membuat instance baru CDCHandler
func NewCDCHandler(res *resources.BootResources) *CDCHandler {
	return &CDCHandler{
		res:      res,
		registry: registry.GetRegistry(res),
	}
}

// HandleEvent dipanggil oleh Kafka consumer setiap ada event masuk
func (h *CDCHandler) HandleEvent(ctx context.Context, table, op string, after, before map[string]interface{}) error {
	if handler, ok := h.registry[table]; ok {
		return handler(ctx, op, after, before)
	}
	log.Printf("[CDC] no handler registered for table=%s", table)
	return nil
}
