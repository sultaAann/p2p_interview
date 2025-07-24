package errors

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
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

func ParsePostgresError(err error, operation string) error {
	if err == nil {
		return nil
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		return &DatabaseError{
			Operation: operation,
			Details:   err.Error(),
		}
	}

	switch pqErr.Code {
	case "23505":
		return parseUniqueConstraintError(pqErr, operation)
	case "23503":
		return &ConstraintViolationError{
			Constraint: pqErr.Constraint,
			Details:    pqErr.Detail,
			Operation:  operation,
		}
	case "23502":
		return &ConstraintViolationError{
			Constraint: pqErr.Constraint,
			Details:    pqErr.Detail,
			Operation:  operation,
		}
	default:
		return &DatabaseError{
			Operation: operation,
			Details:   pqErr.Message,
		}
	}
}

func parseUniqueConstraintError(pqErr *pq.Error, operation string) error {
	field, resource := parseConstraintName(pqErr.Constraint)
	value := extractValueFromDetail(pqErr.Detail)

	if field != "" && resource != "" {
		return &DuplicateKeyError{
			Field:    field,
			Value:    value,
			Resource: resource,
		}
	}

	return &ConstraintViolationError{
		Constraint: pqErr.Constraint,
		Details:    pqErr.Detail,
		Operation:  operation,
	}
}

func parseConstraintName(constraint string) (field, resource string) {
	constraintMap := map[string]struct {
		field    string
		resource string
	}{
		"users_email_key": {"email", "user"},
		"users_login_key": {"login", "user"},
	}

	if info, exists := constraintMap[constraint]; exists {
		return info.field, info.resource
	}

	parts := strings.Split(constraint, "_")
	if len(parts) >= 3 && parts[len(parts)-1] == "key" {
		resource = parts[0]
		field = strings.Join(parts[1:len(parts)-1], "_")
		return field, resource
	}

	return "", ""
}

func extractValueFromDetail(detail string) string {
	start := strings.Index(detail, "=(")
	if start == -1 {
		return "unknown"
	}
	start += 2

	end := strings.Index(detail[start:], ")")
	if end == -1 {
		return "unknown"
	}
	end += start

	return detail[start:end]
}
