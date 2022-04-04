package gamestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	db "github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	var game *entity.Game

	query := `SELECT id, player_name, created_at, battle_won, scenario,
		FROM games
		LEFT JOIN pokemon on partner_id = pokemon.id
		WHERE id = ?`

	if err := mappingGame(s.db.QueryRowContext(ctx, query, gameID), game); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to find game with id %s", gameID)
		}
		return nil, fmt.Errorf("unable to find game with id %s: %v", gameID, err)
	}

	return game, nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	queryGame := `
		INSERT INTO games (id, player_name, created_at, battle_won, scenario, partner_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	queryPartner := `
		INSERT INTO pokemon (id, name, max_health, attack, defense, speed, avatar_url)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, queryGame, game.ID, game.PlayerName, game.CreatedAt, game.BattleWon, game.Scenario, game.Partner.ID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, queryPartner, game.Partner.ID, game.Partner.Name, game.Partner.BattleStats.MaxHealth, game.Partner.BattleStats.Attack, game.Partner.BattleStats.Defense, game.Partner.BattleStats.Speed, game.Partner.AvatarURL); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func mappingGame(row db.RowResultInterface, g *entity.Game) error {
	return row.Scan(
		&g.ID, &g.PlayerName, &g.CreatedAt,
		&g.BattleWon, &g.Scenario,
		&g.Partner.ID, &g.Partner.Name, &g.Partner.AvatarURL,
		&g.Partner.BattleStats.MaxHealth,
		&g.Partner.BattleStats.Attack,
		&g.Partner.BattleStats.Defense,
		&g.Partner.BattleStats.Speed,
	)
}
