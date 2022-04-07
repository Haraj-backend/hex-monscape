package battlestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	db "github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*battle.Battle, error) {
	var battle battle.Battle
	var partner, enemy entity.Pokemon
	battle.Partner = &partner
	battle.Enemy = &enemy

	query := `
	SELECT b.game_id, b.partner_last_damage, b.enemy_last_damage, b.state,
	p.pokemon_id as partner_id, p.name as partner_name, p.max_health as partner_max_health, p.health as partner_health,
	p.attack as partner_attack, p.defense as partner_defense, p.speed as partner_speed, p.avatar_url as partner_avatar_url,
	e.pokemon_id as enemy_id, e.name as enemy_name, e.max_health as enemy_max_health, e.health as enemy_health,
	e.attack as enemy_attack, e.defense as enemy_defense, e.speed as enemy_speed, e.avatar_url as enemy_avatar_url
	FROM battles b
	LEFT JOIN pokemon_battle_states p on b.partner_state_id = p.id
	LEFT JOIN pokemon_battle_states e on b.enemy_state_id = e.id
	WHERE game_id = ?
		`

	if err := mappingBattle(s.db.QueryRowContext(ctx, query, gameID), &battle); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to find battle with id %s: %v", gameID, err)
	}

	return &battle, nil
}

func (s *Storage) SaveBattle(ctx context.Context, b battle.Battle) error {
	exists, partnerStateId, enemyStateId, err := s.checkBattleExists(ctx, b.GameID)

	if err != nil {
		return err
	}

	if exists {
		return s.updateBattle(ctx, &b, partnerStateId, enemyStateId)
	}

	return s.insertBattle(ctx, b)
}

func (s *Storage) updateBattle(ctx context.Context, b *battle.Battle, partnerStateId, enemyStateId int) error {
	query := `
	UPDATE battles
	SET partner_last_damage = ?, enemy_last_damage = ?, state = ?
	WHERE game_id = ?
		`

	queryPoke := `
		UPDATE pokemon_battle_states
		SET health = ?
		WHERE id = ?
		`

	if _, err := s.db.ExecContext(ctx, query, b.LastDamage.Partner, b.LastDamage.Enemy, b.State, b.GameID); err != nil {
		return fmt.Errorf("unable to update battle with id %s: %v", b.GameID, err)
	}

	if _, err := s.db.ExecContext(ctx, queryPoke, b.Partner.BattleStats.Health, partnerStateId); err != nil {
		return fmt.Errorf("unable to update partner battle state with id %s: %v", b.GameID, err)
	}

	if _, err := s.db.ExecContext(ctx, queryPoke, b.Enemy.BattleStats.Health, enemyStateId); err != nil {
		return fmt.Errorf("unable to update enemy battle state with id %s: %v", b.GameID, err)
	}

	return nil
}

func (s *Storage) checkBattleExists(ctx context.Context, gameID string) (bool, int, int, error) {
	query := `
	SELECT 1, partner_state_id, enemy_state_id
	FROM battles
	WHERE game_id = ?
		`

	var count int
	var partnerStateId int
	var enemyStateId int
	if err := s.db.QueryRowContext(ctx, query, gameID).Scan(&count, &partnerStateId, &enemyStateId); err != nil {
		if err == sql.ErrNoRows {
			return false, 0, 0, nil
		}
		return false, 0, 0, fmt.Errorf("unable to check battle with id %s: %v", gameID, err)
	}

	return count > 0, partnerStateId, enemyStateId, nil
}

func (s *Storage) insertBattle(ctx context.Context, b battle.Battle) error {
	queryGame := `
		INSERT INTO battles (game_id, partner_last_damage, enemy_last_damage, state, partner_state_id, enemy_state_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	queryPoke := `
		INSERT INTO pokemon_battle_states (pokemon_id, name, max_health, health, attack, defense, speed, avatar_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	partner, err := tx.ExecContext(ctx, queryPoke, b.Partner.ID, b.Partner.Name, b.Partner.BattleStats.MaxHealth, b.Partner.BattleStats.Health, b.Partner.BattleStats.Attack, b.Partner.BattleStats.Defense, b.Partner.BattleStats.Speed, b.Partner.AvatarURL)
	if err != nil {
		tx.Rollback()
		return err
	}

	partnerId, err := partner.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	enemy, err := tx.ExecContext(ctx, queryPoke, b.Enemy.ID, b.Enemy.Name, b.Enemy.BattleStats.MaxHealth, b.Enemy.BattleStats.Health, b.Enemy.BattleStats.Attack, b.Enemy.BattleStats.Defense, b.Enemy.BattleStats.Speed, b.Enemy.AvatarURL)

	if err != nil {
		tx.Rollback()
		return err
	}

	enemyId, err := enemy.LastInsertId()

	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, queryGame, b.GameID, b.LastDamage.Partner, b.LastDamage.Enemy, b.State, partnerId, enemyId); err != nil {
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
