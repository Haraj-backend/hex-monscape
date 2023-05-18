package gamestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "GameStorage: GetGame", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	var game GameRow
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

	span.SetAttributes(attribute.Key("game-id").String(gameID))
	span.SetAttributes(attribute.Key("db.system").String("mysql"))
	span.SetAttributes(attribute.Key("db.statement").String(query))

	if err := s.sqlClient.GetContext(ctx, &game, query, gameID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("unable to find game with id %s: %v", gameID, err)
	}

	return game.ToGame(), nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "GameStorage: SaveGame", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	gameRow := NewGameRow(&game)
	query := `
		REPLACE INTO games (
			id, player_name, created_at, battle_won, scenario, partner_id
		) VALUES (
			:id, :player_name, :created_at, :battle_won, :scenario, :partner_id
		)
	`

	span.SetAttributes(attribute.Key("game-id").String(game.ID))
	span.SetAttributes(attribute.Key("db.system").String("mysql"))
	span.SetAttributes(attribute.Key("db.statement").String(query))

	_, err := s.sqlClient.NamedExecContext(ctx, query, map[string]interface{}{
		"id":          gameRow.ID,
		"player_name": gameRow.PlayerName,
		"created_at":  gameRow.CreatedAt,
		"battle_won":  gameRow.BattleWon,
		"scenario":    gameRow.Scenario,
		"partner_id":  gameRow.Partner.ID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return fmt.Errorf("unable to execute query due: %w", err)
	}
	return nil
}
