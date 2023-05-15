package battlestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
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

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*battle.Battle, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "BattleStorage: GetBattle", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	query := `SELECT * FROM battles WHERE game_id = ?`

	span.SetAttributes(attribute.Key("game-id").String(gameID))
	span.SetAttributes(attribute.Key("db.system").String("mysql"))
	span.SetAttributes(attribute.Key("db.statement").String(query))

	var row battleRow
	if err := s.sqlClient.GetContext(ctx, &row, query, gameID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		span.RecordError(err)
		span.SetAttributes(attribute.Key("error").Bool(true))
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}
	return row.ToBattle(), nil
}

func (s *Storage) SaveBattle(ctx context.Context, b battle.Battle) error {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "BattleStorage: SaveBattle", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	battleRow := newBattleRow(b)
	query := `
		REPLACE INTO battles (
			game_id, state, partner_pokemon_id, 
			partner_name, partner_max_health, partner_health, 
			partner_attack, partner_defense, partner_speed, 
			partner_avatar_url, partner_last_damage, enemy_pokemon_id, 
			enemy_name, enemy_max_health, enemy_health,
			enemy_attack, enemy_defense, enemy_speed, 
			enemy_avatar_url, enemy_last_damage
		) VALUES (
			:game_id, :state, :partner_pokemon_id, 
			:partner_name, :partner_max_health, :partner_health, 
			:partner_attack, :partner_defense, :partner_speed, 
			:partner_avatar_url, :partner_last_damage, :enemy_pokemon_id, 
			:enemy_name, :enemy_max_health, :enemy_health,
			:enemy_attack, :enemy_defense, :enemy_speed, 
			:enemy_avatar_url, :enemy_last_damage
		)
	`

	span.SetAttributes(attribute.Key("game-id").String(battleRow.GameID))
	span.SetAttributes(attribute.Key("db.system").String("mysql"))
	span.SetAttributes(attribute.Key("db.statement").String(query))

	_, err := s.sqlClient.NamedExecContext(ctx, query, battleRow)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("error").Bool(true))

		return fmt.Errorf("unable to execute query due: %w", err)
	}
	return nil
}
