package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"gopkg.in/validator.v2"
)

type MonsterStorage struct {
	partnerMap map[string]entity.Monster
	enemyMap   map[string]entity.Monster
}

type MonsterStorageConfig struct {
	MonsterData []byte `validate:"nonzero"`
}

func (c MonsterStorageConfig) Validate() error {
	return validator.Validate(c)
}

func NewMonsterStorage(cfg MonsterStorageConfig) (*MonsterStorage, error) {
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

	// split monsters into partners & enemies
	partnerMap := map[string]entity.Monster{}
	enemyMap := map[string]entity.Monster{}
	for _, row := range rows {
		// we want every monster to be fightable
		monster := row.toMonster()
		enemyMap[monster.ID] = monster

		// add to partner map if it is partnerable
		if row.IsPartnerable {
			partnerMap[monster.ID] = monster
		}
	}

	return &MonsterStorage{partnerMap: partnerMap, enemyMap: enemyMap}, nil
}

func (s *MonsterStorage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	var partners []entity.Monster
	for _, partner := range s.partnerMap {
		partners = append(partners, partner)
	}

	return partners, nil
}

func (s *MonsterStorage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	partner, ok := s.partnerMap[partnerID]
	if !ok {
		return nil, nil
	}

	return &partner, nil
}

func (s *MonsterStorage) GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error) {
	var enemies []entity.Monster
	for _, enemy := range s.enemyMap {
		enemies = append(enemies, enemy)
	}

	return enemies, nil
}
