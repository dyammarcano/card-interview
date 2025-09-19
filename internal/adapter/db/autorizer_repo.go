package db

import (
	"sync"
	"time"

	"stone/cards/authorizer/internal/domain/entities"

	"github.com/google/uuid"
)

// databaseMap is a concurrent map storing authorizer keyed by UUID.
var databaseMap sync.Map

type AuthorizerRepository struct{}

func (r *AuthorizerRepository) InsertAuthorizer(authorizer entities.Authorizer) (uuid.UUID, error) {
	uid := uuid.New()
	databaseMap.Store(uid, authorizer)

	return uid, nil
}

// CountByCardSince returns the number of transactions for the given card with Timestamp after the provided sinceTime (unix nano)
func (r *AuthorizerRepository) CountByCardSince(cardNumber string, sinceTime int64) int {
	since := time.Unix(0, sinceTime)
	count := 0

	databaseMap.Range(func(key, value any) bool {
		v, ok := value.(entities.Authorizer)
		if ok {
			if v.CardNumber == cardNumber && v.Timestamp.After(since) {
				count++
			}

			return true
		}

		return false
	})

	return count
}

func NewAuthorizerRepository() *AuthorizerRepository {
	return &AuthorizerRepository{}
}
