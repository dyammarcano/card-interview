package db

import (
	"sync"

	"stone/cards/authorizer/internal/domain/entities"
)

var riskStorage = make([]entities.Risk, 0)

type RiskRepository struct {
	mu sync.Mutex
}

func (r *RiskRepository) InsertRisk(risk entities.Risk) {
	r.mu.Lock()

	riskStorage = append(riskStorage, risk)

	r.mu.Unlock()
}

func NewRiskRepository() *RiskRepository {
	return &RiskRepository{}
}
