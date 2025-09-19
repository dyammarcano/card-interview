package ctrl

import (
	"encoding/json"
	"time"

	"stone/cards/authorizer/internal/adapter/ctrl/schema"
	"stone/cards/authorizer/internal/domain/entities"

	"github.com/google/uuid"
)

type AuthorizerUseCase interface {
	Authorize(authorizer entities.Authorizer) (uuid.UUID, string, error)
}

type AuthorizerCtrl struct {
	authorizerUC AuthorizerUseCase
}

func (a AuthorizerCtrl) Authorize(payload json.RawMessage) schema.AuthorizerResponse {
	type input struct {
		CardNumber string  `json:"card_number"`
		Amount     float64 `json:"amount"`
		Currency   string  `json:"currency"`
		Merchant   string  `json:"merchant"`
		Timestamp  string  `json:"timestamp"`
	}

	in := input{}
	if err := json.Unmarshal(payload, &in); err != nil {
		return schema.AuthorizerResponse{Status: "rejected", Error: "invalid payload"}
	}

	timestamp, err := time.Parse(time.RFC3339, in.Timestamp)
	if err != nil {
		return schema.AuthorizerResponse{Status: "rejected", Error: "timestamp not valid"}
	}

	if timestamp.After(time.Now()) {
		return schema.AuthorizerResponse{Status: "rejected", Error: "timestamp on future"}
	}

	authorizer := entities.Authorizer{
		CardNumber: in.CardNumber,
		Amount:     in.Amount,
		Currency:   in.Currency,
		Merchant:   in.Merchant,
		Timestamp:  timestamp,
	}

	id, warning, err := a.authorizerUC.Authorize(authorizer)
	if err != nil {
		return schema.AuthorizerResponse{Status: "rejected", Error: "invalid payload"}
	}

	resp := schema.AuthorizerResponse{AuthorizeID: id.String()}
	if warning != "" {
		resp.Status = "approved_with_warning"
		resp.Warning = "transaction marked as suspicious: " + warning

		return resp
	}

	resp.Status = "approved"

	return resp
}

func NewAuthorizerCtrl(authorizerUC AuthorizerUseCase) AuthorizerCtrl {
	return AuthorizerCtrl{
		authorizerUC: authorizerUC,
	}
}
