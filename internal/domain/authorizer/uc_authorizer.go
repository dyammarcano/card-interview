package authorizer

import (
	"time"

	"stone/cards/authorizer/internal/domain/entities"

	"github.com/google/uuid"
)

const (
	riskAmountLimit      = 10_000.0
	riskAuthorizersLimit = 5
)

type AuthorizerRepository interface {
	InsertAuthorizer(authorizer entities.Authorizer) (uuid.UUID, error)
	CountByCardSince(cardNumber string, sinceTime int64) int
}

type RiskRepository interface {
	InsertRisk(risks entities.Risk)
}

type AuthorizerUC struct {
	authorizerRepo AuthorizerRepository
	riskRepo       RiskRepository
}

func (a AuthorizerUC) Authorize(authorize entities.Authorizer) (uuid.UUID, string, error) {
	warning := ""

	if authorize.Amount > riskAmountLimit {
		warning = entities.RiskHighAmount
		a.riskRepo.InsertRisk(entities.Risk{
			CardNumber: authorize.CardNumber,
			Reason:     entities.RiskHighAmount,
			Timestamp:  authorize.Timestamp,
		})
	}

	since := authorize.Timestamp.Add(-1 * time.Minute).UnixNano()
	count := a.authorizerRepo.CountByCardSince(authorize.CardNumber, since)

	if count >= riskAuthorizersLimit {
		warning = entities.RiskNotStandard
		a.riskRepo.InsertRisk(entities.Risk{
			CardNumber: authorize.CardNumber,
			Reason:     entities.RiskNotStandard,
			Timestamp:  authorize.Timestamp,
		})
	}

	id, err := a.authorizerRepo.InsertAuthorizer(authorize)

	return id, warning, err
}

func NewAuthorizerUC(authRepo AuthorizerRepository, riskRepo RiskRepository) AuthorizerUC {
	return AuthorizerUC{
		authorizerRepo: authRepo,
		riskRepo:       riskRepo,
	}
}
