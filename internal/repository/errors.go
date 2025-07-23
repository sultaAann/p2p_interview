package repository

import "fmt"

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
