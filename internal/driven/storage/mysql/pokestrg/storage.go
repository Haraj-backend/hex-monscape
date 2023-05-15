package pokestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
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

const (
	partner int = 1
	enemy   int = 0
)

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &Storage{sqlClient: cfg.SQLClient}
	return s, nil
}

func (s *Storage) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "PokeStorage: GetAvailablePartners", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	return s.getPokemonsByRole(ctx, partner)
}

func (s *Storage) GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "PokeStorage: GetPossibleEnemies", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	return s.getPokemonsByRole(ctx, enemy)
}

func (s *Storage) GetPartner(ctx context.Context, partnerID string) (*entity.Pokemon, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "PokeStorage: GetPartner", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	var pokemon shared.PokeRow
	query := `
		SELECT
			id,
			name,
			health,
			max_health,
			attack,
			defense,
			speed,
			avatar_url
		FROM pokemons
		WHERE id = ?
	`

	span.SetAttributes(attribute.Key("partner-id").String(partnerID))
	span.SetAttributes(attribute.Key("db.system").String("mysql"))
	span.SetAttributes(attribute.Key("db.statement").String(query))

	if err := s.sqlClient.GetContext(ctx, &pokemon, query, partnerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		span.RecordError(err)
		span.SetAttributes(attribute.Key("error").Bool(true))
		return nil, fmt.Errorf("unable to find partner with id %s: %v", partnerID, err)
	}

	return pokemon.ToPokemon(), nil
}

func (s *Storage) getPokemonsByRole(ctx context.Context, isPartnerable int) ([]entity.Pokemon, error) {
	tr := telemetry.GetTracer()
	ctx, span := tr.Trace(ctx, "PokeStorage: getPokemonsByRole", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	var pokemons shared.PokeRows

	query := `
		SELECT
			id,
			name,
			health,
			max_health,
			attack,
			defense,
			speed,
			avatar_url
		FROM pokemons
		WHERE is_partnerable = ?
	`

	span.SetAttributes(attribute.Key("db.system").String("mysql"))
	span.SetAttributes(attribute.Key("db.statement").String(query))

	if err := s.sqlClient.SelectContext(ctx, &pokemons, query, isPartnerable); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		span.RecordError(err)
		span.SetAttributes(attribute.Key("error").Bool(true))
		return nil, fmt.Errorf("unable to find pokemon's role %d: %v", isPartnerable, err)
	}

	return pokemons.ToPokemons(), nil
}
