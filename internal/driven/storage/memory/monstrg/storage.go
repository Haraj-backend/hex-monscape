package monstrg

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"gopkg.in/validator.v2"
)

type Storage struct {
	partnerMap map[string]entity.Monster
	enemyMap   map[string]entity.Monster
}

func (s *Storage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	var partners []entity.Monster
	for _, partner := range s.partnerMap {
		partners = append(partners, partner)
	}
	return partners, nil
}

func (s *Storage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	partner, ok := s.partnerMap[partnerID]
	if !ok {
		return nil, nil
	}
	return &partner, nil
}

func (s *Storage) GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error) {
	var enemies []entity.Monster
	for _, enemy := range s.enemyMap {
		enemies = append(enemies, enemy)
	}
	return enemies, nil
}

type Config struct {
	Partners []entity.Monster `validate:"min=1"`
	Enemies  []entity.Monster `validate:"min=1"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	partnerMap := map[string]entity.Monster{}
	for _, partner := range cfg.Partners {
		partnerMap[partner.ID] = partner
	}
	enemyMap := map[string]entity.Monster{}
	for _, enemy := range cfg.Enemies {
		enemyMap[enemy.ID] = enemy
	}
	return &Storage{partnerMap: partnerMap, enemyMap: enemyMap}, nil
}
