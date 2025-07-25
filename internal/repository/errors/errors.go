package errors

import (
	"fmt"
)

type DatabaseError struct {
	Operation string
	Details   string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %s", e.Operation, e.Details)
}

func (e *DatabaseError) Is(target error) bool {
	_, ok := target.(*DatabaseError)
	return ok
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

type DuplicateKeyError struct {
	Field    string
	Value    string
	Resource string
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("%s with %s '%s' already exists", e.Resource, e.Field, e.Value)
}

func (e *DuplicateKeyError) Is(target error) bool {
	_, ok := target.(*DuplicateKeyError)
	return ok
}

type ConstraintViolationError struct {
	Constraint string
	Details    string
	Operation  string
}

func (e *ConstraintViolationError) Error() string {
	return fmt.Sprintf("constraint violation during %s: %s (%s)", e.Operation, e.Details, e.Constraint)
}

func (e *ConstraintViolationError) Is(target error) bool {
	_, ok := target.(*ConstraintViolationError)
	return ok
}
