package monstrg

import (
	"context"
	"encoding/json"
	"fmt"

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
	MonsterData []byte `validate:"nonzero"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	// validate config
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	// parse monster data
	var rows []monsterRow
	err = json.Unmarshal(cfg.MonsterData, &rows)
	if err != nil {
		return nil, fmt.Errorf("unable to parse monster data due: %w", err)
	}
	partnerMap := map[string]entity.Monster{}
	enemyMap := map[string]entity.Monster{}
	for _, monsterRow := range rows {
		// we want every monster to be fightable
		monster := monsterRow.toMonster()
		enemyMap[monster.ID] = monster
		// add to partner map if it is partnerable
		if monsterRow.IsPartnerable {
			partnerMap[monster.ID] = monster
		}
	}
	return &Storage{partnerMap: partnerMap, enemyMap: enemyMap}, nil
}
