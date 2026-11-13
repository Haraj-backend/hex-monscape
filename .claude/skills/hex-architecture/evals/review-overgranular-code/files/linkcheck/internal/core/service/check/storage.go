package check

import (
	"context"

	"example.com/linkcheck/internal/core/entity"
)

type Source struct {
	Name string
	URLs []string
}

type LinkStorage interface {
	GetSources(ctx context.Context) ([]Source, error)
}

type StateStorage interface {
	GetState(ctx context.Context) (*entity.State, error)
	SaveState(ctx context.Context, state entity.State) error
}
