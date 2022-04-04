package battlestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	db "github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*battle.Battle, error) {
	var battle battle.Battle

	query := `
	SELECT game_id, partner_last_damage, enemy_last_damage, state,
	p.pokemon_id as partner_id, p.name as partner_name, p.max_health as partner_max_health, p.health as partner_health,
	p.attack as partner_attack, p.defense as partner_defense, p.speed as partner_speed, p.avatar_url as partner_avatar_url,
	e.pokemon_id as enemy_id, e.name as enemy_name, e.max_health as enemy_max_health, e.health as enemy_health,
	e.attack as enemy_attack, e.defense as enemy_defense, e.speed as enemy_speed, e.avatar_url as enemy_avatar_url,
	FROM battles
	LEFT JOIN pokemon_battle_states p on partner_state_id = p.id
	LEFT JOIN pokemon_battle_states e on enemy_state_id = e.id
	WHERE game_id = ?
		`

	if err := mappingBattle(s.db.QueryRowContext(ctx, query, gameID), &battle); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to find game with id %s", gameID)
		}
		return nil, fmt.Errorf("unable to find game with id %s: %v", gameID, err)
	}

	return &battle, nil
}

func (s *Storage) SaveBattle(ctx context.Context, b battle.Battle) error {
	queryGame := `
		INSERT INTO battles (game_id, partner_last_damage, enemy_last_damage, state)
		VALUES (?, ?, ?, ?)
	`
	queryPoke := `
		INSERT INTO pokemon_battle_states (id, name, max_health, health, attack, defense, speed, avatar_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, queryGame, b.GameID, b.LastDamage.Partner, b.LastDamage.Enemy, b.State); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, queryPoke, b.Partner.ID, b.Partner.Name, b.Partner.BattleStats.MaxHealth, b.Partner.BattleStats.Health, b.Partner.BattleStats.Attack, b.Partner.BattleStats.Defense, b.Partner.BattleStats.Speed, b.Partner.AvatarURL); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, queryPoke, b.Enemy.ID, b.Enemy.Name, b.Enemy.BattleStats.MaxHealth, b.Enemy.BattleStats.Health, b.Enemy.BattleStats.Attack, b.Enemy.BattleStats.Defense, b.Enemy.BattleStats.Speed, b.Enemy.AvatarURL); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func mappingBattle(row db.RowResultInterface, b *battle.Battle) error {
	return row.Scan(
		&b.GameID, &b.LastDamage.Partner, &b.LastDamage.Enemy, &b.State,
		&b.Partner.ID, &b.Partner.Name, &b.Partner.BattleStats.MaxHealth, &b.Partner.BattleStats.Health,
		&b.Partner.BattleStats.Attack, &b.Partner.BattleStats.Defense, &b.Partner.BattleStats.Speed, &b.Partner.AvatarURL,
		&b.Enemy.ID, &b.Enemy.Name, &b.Enemy.BattleStats.MaxHealth, &b.Enemy.BattleStats.Health,
		&b.Enemy.BattleStats.Attack, &b.Enemy.BattleStats.Defense, &b.Enemy.BattleStats.Speed, &b.Enemy.AvatarURL,
	)
}
