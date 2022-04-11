package shared

import (
	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
)

type PokeRow struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Health    int    `db:"health"`
	MaxHealth int    `db:"max_health"`
	Attack    int    `db:"attack"`
	Defense   int    `db:"defense"`
	Speed     int    `db:"speed"`
	AvatarURL string `db:"avatar_url"`
}

type PokeRows []PokeRow

func (r *PokeRow) ToPokemon() *entity.Pokemon {
	return &entity.Pokemon{
		ID:   r.ID,
		Name: r.Name,
		BattleStats: entity.BattleStats{
			Health:    r.Health,
			MaxHealth: r.MaxHealth,
			Attack:    r.Attack,
			Defense:   r.Defense,
			Speed:     r.Speed,
		},
		AvatarURL: r.AvatarURL,
	}
}

func (r PokeRows) ToPokemons() []entity.Pokemon {
	var pokemons []entity.Pokemon
	for _, row := range r {
		pokemons = append(pokemons, *row.ToPokemon())
	}
	return pokemons
}

type GameRow struct {
	ID         string          `db:"id"`
	PlayerName string          `db:"player_name"`
	CreatedAt  int64           `db:"created_at"`
	BattleWon  int             `db:"battle_won"`
	Scenario   entity.Scenario `db:"scenario"`
	Partner    *PokeRow        `db:"partner"`
}

func (r *GameRow) ToGame() *entity.Game {
	return &entity.Game{
		ID:         r.ID,
		PlayerName: r.PlayerName,
		CreatedAt:  r.CreatedAt,
		BattleWon:  r.BattleWon,
		Scenario:   r.Scenario,
		Partner:    r.Partner.ToPokemon(),
	}
}

func NewGameRow(game *entity.Game) *GameRow {
	return &GameRow{
		ID:         game.ID,
		PlayerName: game.PlayerName,
		CreatedAt:  game.CreatedAt,
		BattleWon:  game.BattleWon,
		Scenario:   game.Scenario,
		Partner:    NewPokeRow(game.Partner),
	}
}

func NewPokeRow(pokemon *entity.Pokemon) *PokeRow {
	return &PokeRow{
		ID:        pokemon.ID,
		Name:      pokemon.Name,
		Health:    pokemon.BattleStats.Health,
		MaxHealth: pokemon.BattleStats.MaxHealth,
		Attack:    pokemon.BattleStats.Attack,
		Defense:   pokemon.BattleStats.Defense,
		Speed:     pokemon.BattleStats.Speed,
		AvatarURL: pokemon.AvatarURL,
	}
}
