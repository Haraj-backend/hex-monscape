package entity

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/validator.v2"
)

var ErrInvalidState = errors.New("invalid state for given action")

type Battle struct {
	GameID     string     `json:"game_id"`
	State      State      `json:"state"`
	Partner    *Monster   `json:"partner"`
	Enemy      *Monster   `json:"enemy"`
	LastDamage LastDamage `json:"last_damage"`
}

// PartnerAttack is used for executing partner attack. The battle state must
// be PARTNER_TURN otherwise the action will be rejected.
func (b *Battle) PartnerAttack() error {
	// check battle turn
	if b.State != PARTNER_TURN {
		return ErrInvalidState
	}
	// inflict damage to enemy
	damage, err := b.Enemy.InflictDamage(*b.Partner)
	if err != nil {
		return fmt.Errorf("unable to inflict damage to enemy due: %w", err)
	}
	// set enemy last damage
	b.LastDamage.Enemy = damage
	// set battle state accordingly
	if b.Enemy.IsDead() {
		b.State = WIN
	} else {
		b.State = DECIDE_TURN
	}

	return nil
}

// PartnerSurrender is used for executing partner surrender. The battle state must
// be PARTNER_TURN otherwise the action will be rejected.
func (b *Battle) PartnerSurrender() error {
	// check battle turn
	if b.State != PARTNER_TURN {
		return ErrInvalidState
	}
	// set state to lose
	b.State = LOSE
	return nil
}

// EnemyAttack is used for executing enemy attack. The battle state must be
// ENEMY_TURN otherwise the action will be rejected.
func (b *Battle) EnemyAttack() error {
	// check battle state
	if b.State != ENEMY_TURN {
		return ErrInvalidState
	}
	// inflict damage to partner
	damage, err := b.Partner.InflictDamage(*b.Enemy)
	if err != nil {
		return fmt.Errorf("unable to inflict damage to partner due: %w", err)
	}
	// set partner last damage
	b.LastDamage.Partner = damage
	// set battle state accordingly
	if b.Partner.IsDead() {
		b.State = LOSE
	} else {
		b.State = DECIDE_TURN
	}

	return nil
}

// IsEnded returns true when state is either WIN or LOSE
func (b Battle) IsEnded() bool {
	return b.State == WIN || b.State == LOSE
}

// DecideTurn is used for deciding turn in the battle. It calculates turn based
// on speed of both partner & enemy. The battle state must be DECIDE_TURN, otherwise
// the action will be rejected.
func (b *Battle) DecideTurn() (State, error) {
	if b.State != DECIDE_TURN {
		return "", ErrInvalidState
	}
	// define slot for both partner & enemy
	partnerSlot := 0
	enemySlot := 1
	// calculate partner attack chance
	lenSlots := 10
	totalSpeed := b.Enemy.BattleStats.Speed + b.Partner.BattleStats.Speed
	partnerAttackChance := int((float64(b.Partner.BattleStats.Speed) / float64(totalSpeed)) * float64(lenSlots))
	// fill slots
	slots := make([]int, 0, lenSlots)
	for i := 0; i < partnerAttackChance; i++ {
		slots = append(slots, partnerSlot)
	}
	for i := 0; i < lenSlots-partnerAttackChance; i++ {
		slots = append(slots, enemySlot)
	}
	// decide turn
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := rnd.Intn(lenSlots)
	state := ENEMY_TURN
	if slots[idx] == partnerSlot {
		state = PARTNER_TURN
	}
	// update battle internal state
	b.State = state

	return state, nil
}

type BattleConfig struct {
	GameID  string   `validate:"nonzero"`
	Partner *Monster `validate:"nonnil"`
	Enemy   *Monster `validate:"nonnil"`
}

func (c BattleConfig) Validate() error {
	return validator.Validate(c)
}

// NewBattle returns new instance of Battle based on given config.
func NewBattle(cfg BattleConfig) (*Battle, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	b := &Battle{
		GameID:     cfg.GameID,
		State:      DECIDE_TURN,
		Partner:    cfg.Partner,
		Enemy:      cfg.Enemy,
		LastDamage: LastDamage{},
	}
	return b, nil
}

// State represents current battle state
type State string

const (
	DECIDE_TURN  State = "DECIDE_TURN"
	ENEMY_TURN   State = "ENEMY_TURN"
	PARTNER_TURN State = "PARTNER_TURN"
	WIN          State = "WIN"
	LOSE         State = "LOSE"
)

type LastDamage struct {
	Partner int `json:"partner"`
	Enemy   int `json:"enemy"`
}
