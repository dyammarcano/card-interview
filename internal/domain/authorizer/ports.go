package authorizer

import (
	"stone/cards/authorizer/internal/domain/entities"

	"github.com/google/uuid"
)

type AuthorizerRepository interface {
	InsertAuthorizer(authorizer entities.Authorizer) (uuid.UUID, error)
	CountByCardSince(cardNumber string, sinceTime int64) int
}

type RiskRepository interface {
	InsertRisk(risks entities.Risk)
}
