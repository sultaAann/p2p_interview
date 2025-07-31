package repository

import (
	"p2p_interview/internal/domain/models"
	"time"
)

type SlotRepository interface {
	GetAvailableSlots() ([]models.Slot, error)
	GetStartTimeAndEndTimeOfAvailableSlots() (*time.Time, *time.Time, error)
	GetSlotById(id string) (*models.Slot, error)
	CreateSlot(slot *models.Slot) (string, error)
	UpdateSlot(slot *models.Slot) error
	DeleteSlot(id string) error
}
