package battlestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type Storage struct {
	sqlClient *sqlx.DB
}

type Config struct {
	SQLClient *sqlx.DB `validate:"nonnil"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &Storage{sqlClient: cfg.SQLClient}
	return s, nil
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*battle.Battle, error) {
	var row battleRow
	query := `
		SELECT
			b.game_id,
			b.partner_last_damage,
			b.enemy_last_damage,
			b.state,
			p.pokemon_id as 'partner.id',
			p.name as 'partner.name',
			p.max_health as 'partner.max_health',
			p.health as 'partner.health',
			p.attack as 'partner.attack',
			p.defense as 'partner.defense',
			p.speed as 'partner.speed',
			p.avatar_url as 'partner.avatar_url',
			e.pokemon_id as 'enemy.id',
			e.name as 'enemy.name',
			e.max_health as 'enemy.max_health',
			e.health as 'enemy.health',
			e.attack as 'enemy.attack',
			e.defense as 'enemy.defense',
			e.speed as 'enemy.speed',
			e.avatar_url as 'enemy.avatar_url'
		FROM battles b
		LEFT JOIN pokemon_battle_states p on b.partner_state_id = p.id
		LEFT JOIN pokemon_battle_states e on b.enemy_state_id = e.id
		WHERE game_id = ?
	`
	if err := s.sqlClient.GetContext(ctx, &row, query, gameID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to find battle with id %s: %v", gameID, err)
	}

	return row.ToBattle(), nil
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
	tx, err := s.sqlClient.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction when update battle due: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			err = fmt.Errorf("unable to commit transaction update battle due: %w", err)
			return
		}
	}()

	query := `
		UPDATE battles
		SET
			partner_last_damage = :partner_last_damage,
			enemy_last_damage = :enemy_last_damage,
			state = :state
		WHERE game_id = :game_id
	`

	queryPoke := `
		UPDATE pokemon_battle_states
		SET
			health = :health
		WHERE id = :id
	`

	_, err = tx.NamedExecContext(ctx, query, map[string]interface{}{
		"partner_last_damage": b.LastDamage.Partner,
		"enemy_last_damage":   b.LastDamage.Enemy,
		"state":               b.State,
		"game_id":             b.GameID,
	})
	if err != nil {
		return fmt.Errorf("unable to update battle with id %s: %v", b.GameID, err)
	}

	_, err = tx.NamedExecContext(ctx, queryPoke, map[string]interface{}{
		"health": b.Partner.BattleStats.Health,
		"id":     partnerStateId,
	})
	if err != nil {
		return fmt.Errorf("unable to update partner battle state with id %s: %v", b.GameID, err)
	}

	_, err = tx.NamedExecContext(ctx, queryPoke, map[string]interface{}{
		"health": b.Enemy.BattleStats.Health,
		"id":     enemyStateId,
	})
	if err != nil {
		return fmt.Errorf("unable to update enemy battle state with id %s: %v", b.GameID, err)
	}

	return nil
}

func (s *Storage) checkBattleExists(ctx context.Context, gameID string) (bool, int, int, error) {
	query := `
		SELECT
			1,
			partner_state_id,
			enemy_state_id
		FROM battles
		WHERE game_id = ?
	`

	var count int
	var partnerStateId int
	var enemyStateId int
	if err := s.sqlClient.QueryRowContext(ctx, query, gameID).Scan(&count, &partnerStateId, &enemyStateId); err != nil {
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

	tx, err := s.sqlClient.BeginTx(ctx, nil)
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
