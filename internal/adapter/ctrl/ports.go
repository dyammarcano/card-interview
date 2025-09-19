package ctrl

import (
	"stone/cards/authorizer/internal/domain/entities"

	"github.com/google/uuid"
)

// AuthorizerUseCase defines the application boundary for authorizing transactions.
type AuthorizerUseCase interface {
	Authorize(authorizer entities.Authorizer) (uuid.UUID, string, error)
}
