package gamestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	var game shared.GameRow
	query := `
		SELECT
			g.id,
			g.player_name,
			g.created_at,
			g.battle_won,
			g.scenario,
			p.id as 'partner.id',
			p.name as 'partner.name',
			p.avatar_url as 'partner.avatar_url',
			p.health as 'partner.health',
			p.max_health as 'partner.max_health',
			p.attack as 'partner.attack',
			p.defense as 'partner.defense',
			p.speed as 'partner.speed'
		FROM games g
		LEFT JOIN pokemons p on partner_id = p.id
		WHERE g.id = ?
	`

	if err := s.db.GetContext(ctx, &game, query, gameID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to find game with id %s", gameID)
		}
		return nil, fmt.Errorf("unable to find game with id %s: %v", gameID, err)
	}

	return game.ToGame(), nil
}

func (s *Storage) checkGameExists(ctx context.Context, gameID string) (bool, error) {
	query := `
		SELECT
			id
		FROM games
		WHERE id = ?
	`

	var id string
	if err := s.db.QueryRowContext(ctx, query, gameID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("unable to find game with id %s: %v", gameID, err)
	}

	return true, nil
}

func (s *Storage) updateGame(ctx context.Context, g *entity.Game) error {
	query := `
		UPDATE games
		SET
			player_name = :player_name,
			created_at = :created_at,
			battle_won = :battle_won,
			scenario = :scenario,
			partner_id = :partner_id
		WHERE id = :id
	`

	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":          g.ID,
		"player_name": g.PlayerName,
		"created_at":  g.CreatedAt,
		"battle_won":  g.BattleWon,
		"scenario":    g.Scenario,
		"partner_id":  g.Partner.ID,
	})

	return err
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	exists, err := s.checkGameExists(ctx, game.ID)
	if err != nil {
		return err
	}

	if exists {
		return s.updateGame(ctx, &game)
	}

	return s.insertGame(ctx, game)
}

func (s *Storage) insertGame(ctx context.Context, game entity.Game) error {
	queryGame := `
		INSERT INTO games
			(id, player_name, created_at, battle_won, scenario, partner_id)
		VALUES
			(:id, :player_name, :created_at, :battle_won, :scenario, :partner_id)
	`

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.NamedExecContext(ctx, queryGame, map[string]interface{}{
		"id":          game.ID,
		"player_name": game.PlayerName,
		"created_at":  game.CreatedAt,
		"battle_won":  game.BattleWon,
		"scenario":    game.Scenario,
		"partner_id":  game.Partner.ID,
	}); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
