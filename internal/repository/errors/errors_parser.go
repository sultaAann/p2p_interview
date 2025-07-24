package errors

import (
	"strings"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

// PostgreSQL error ErrorParser
type ErrorParser struct {
	logger *zap.Logger
}

func NewErrorParser(logger *zap.Logger) *ErrorParser {
	return &ErrorParser{
		logger: logger,
	}
}

func (p *ErrorParser) ParsePostgresError(err error, operation string) error {
	if err == nil {
		return nil
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		p.logError(operation, err, "non-postgres error")
		return &DatabaseError{
			Operation: operation,
			Details:   err.Error(),
		}
	}

	p.logPostgresError(operation, pqErr)

	switch pqErr.Code {
	case "23505": 
		return p.parseUniqueConstraintError(pqErr, operation)
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

func (p *ErrorParser) parseUniqueConstraintError(pqErr *pq.Error, operation string) error {
	field, resource := p.parseConstraintName(pqErr.Constraint)
	value := p.extractValueFromDetail(pqErr.Detail)

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

func (p *ErrorParser) parseConstraintName(constraint string) (field, resource string) {
	constraintMap := map[string]struct {
		field    string
		resource string
	}{
		"users_email_key":   {"email", "user"},
		"users_login_key":   {"login", "user"},
		"users_phone_key":   {"phone", "user"},
		"products_name_key": {"name", "product"},
		"orders_number_key": {"number", "order"},
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

func (p *ErrorParser) extractValueFromDetail(detail string) string {
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

func (p *ErrorParser) logError(operation string, err error, errorType string) {
	if p.logger != nil {
		p.logger.Error("Database error occurred",
			zap.String("operation", operation),
			zap.String("error_type", errorType),
			zap.Error(err),
		)
	}
}

func (p *ErrorParser) logPostgresError(operation string, pqErr *pq.Error) {
	if p.logger != nil {
		p.logger.Error("PostgreSQL error occurred",
			zap.String("operation", operation),
			zap.String("code", string(pqErr.Code)),
			zap.String("constraint", pqErr.Constraint),
			zap.String("detail", pqErr.Detail),
			zap.String("message", pqErr.Message),
		)
	}
}
