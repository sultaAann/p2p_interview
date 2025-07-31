package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"p2p_interview/internal/domain/models"
	"p2p_interview/internal/repository"
	ce "p2p_interview/internal/repository/errors"
	"time"

	"go.uber.org/zap"
)

var _ repository.SlotRepository = &slotRepository{}

type slotRepository struct {
	db          *sql.DB
	logger      *zap.Logger
	errorParser *ce.ErrorParser
}

func NewSlotRepository(db *sql.DB, logger *zap.Logger, errorParser *ce.ErrorParser) repository.SlotRepository {
	return &slotRepository{db: db, logger: logger, errorParser: errorParser}
}

func (r *slotRepository) GetAvailableSlots() ([]models.Slot, error) {
	r.logger.Info("Fetching all slots")

	res, err := r.db.Query("SELECT id, interviewer_id, interviewee_id, start_time, end_time FROM slots WHERE status = 'available'")
	if err != nil {
		r.logger.Error("Failed to fetch all slots", zap.Error(err))
		return nil, r.errorParser.ParsePostgresError(err, "get_all_slots")
	}

	slots := []models.Slot{}
	for res.Next() {
		slot := models.Slot{}
		err = res.Scan(&slot.ID, &slot.InterviewerID, &slot.IntervieweeID, &slot.StartTime, &slot.EndTime)
		if err != nil {
			r.logger.Error("Failed to scan slot row", zap.Error(err))
			return nil, r.errorParser.ParsePostgresError(err, "scan_slot_row")
		}
		slots = append(slots, slot)
	}

	if err = res.Err(); err != nil {
		r.logger.Error("Error iterating over rows", zap.Error(err))
		return nil, r.errorParser.ParsePostgresError(err, "iterate_slot_rows")
	}

	r.logger.Info("Successfulyy fetched all slots", zap.Int("count", len(slots)))
	return slots, nil
}

// TODO: implement it like get slots from now to available quarter time
func (r *slotRepository) GetStartTimeAndEndTimeOfAvailableSlots() (*time.Time, *time.Time, error) {
	panic("not implemented")
}

func (r *slotRepository) GetSlotById(id string) (*models.Slot, error) {
	r.logger.Info("Fetching slot by ID", zap.String("id", id))

	var slot models.Slot

	err := r.db.QueryRow(
		"SELECT id, interviewer_id, interviewee_id, start_time, end_time, status FROM slots WHERE id = $1", id,
	).Scan(&slot.ID, &slot.InterviewerID, &slot.IntervieweeID, &slot.StartTime, &slot.EndTime, &slot.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Info("Slot not found", zap.String("id", id))
			return nil, &ce.NotFoundError{Resource: "Slot", ID: id}
		}
		r.logger.Error("Database error while fetching slot", zap.Error(err))
		return nil, r.errorParser.ParsePostgresError(err, "get_slot_by_id")
	}

	r.logger.Info("Successfully fetched slot", zap.String("id", id))
	return &slot, nil
}

func (r *slotRepository) CreateSlot(slot *models.Slot) (string, error) {
	r.logger.Info("Creating new slot")

	var id string

	err := r.db.QueryRow(
		"INSERT INTO slots (interviewer_id, start_time, end_time, status) VALUES ($1, $2, $3, $4) RETURNING id",
		slot.InterviewerID, slot.StartTime, slot.EndTime, models.GetStatus(models.StatusAvailable),
	).Scan(&id)

	if err != nil {
		return "", r.errorParser.ParsePostgresError(err, "create_slot")
	}

	r.logger.Info("Successfully created new slot", zap.String("id", id))
	return id, nil
}

func (r *slotRepository) UpdateSlot(slot *models.Slot) error {
	r.logger.Info("Updating slot", zap.String("id", slot.ID))

	query := `
		UPDATE slots 
		SET interviewer_id = $2, 
			interviewee_id = $3, 
			start_time = $4, 
			end_time = $5, 
			status = $6 
		WHERE id = $1
	`
	res, err := r.db.Exec(
		query, slot.ID, slot.InterviewerID, slot.IntervieweeID, slot.StartTime, slot.EndTime, slot.Status,
	)
	if err != nil {
		return r.errorParser.ParsePostgresError(err, "update_slot")
	}
	count, err := res.RowsAffected()
	if err != nil {
		return r.errorParser.ParsePostgresError(err, "update_slot_rows_affected")
	}
	if count == 0 {
		return fmt.Errorf("no slot updated with id: %s", slot.ID)
	}
	r.logger.Info("Successfully Updated slot", zap.String("id", slot.ID))
	return nil
}

func (r *slotRepository) DeleteSlot(id string) error {
	panic("ada")
}
