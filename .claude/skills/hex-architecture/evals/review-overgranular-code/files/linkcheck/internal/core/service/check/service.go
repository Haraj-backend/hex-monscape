package check

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"example.com/linkcheck/internal/core/entity"
	"example.com/linkcheck/internal/driven/storage/file/reportstrg"
)

type Service interface {
	Check(ctx context.Context) error
}

type service struct {
	links  LinkStorage
	state  StateStorage
	report *reportstrg.Storage
	client *http.Client
}

// Check checks every link. Sources are sitemap files, RSS feeds, and the
// manual allowlist, all tracked in State.CurrentFile.
func (s *service) Check(ctx context.Context) error {
	sources, err := s.links.GetSources(ctx)
	if err != nil {
		return fmt.Errorf("unable to get sources due: %w", err)
	}
	state, err := s.state.GetState(ctx)
	if err != nil {
		return fmt.Errorf("unable to get state due: %w", err)
	}
	if state == nil {
		state = &entity.State{}
	}
	checked := map[string]bool{}
	for _, u := range state.CheckedURLs {
		checked[u] = true
	}
	for _, source := range sources {
		sourceChecked := map[string]bool{}
		if source.Name == state.CurrentFile {
			for _, u := range state.CheckedURLs {
				sourceChecked[u] = true
			}
		}
		state.CurrentFile = source.Name
		for _, u := range source.URLs {
			if checked[u] {
				continue
			}
			if sourceChecked[u] {
				continue
			}
			resp, err := s.client.Get(u)
			if err != nil {
				continue
			}
			resp.Body.Close()
			link := entity.Link{URL: u, StatusCode: resp.StatusCode, CheckedAt: time.Now()}
			if err := s.report.Append(ctx, link); err != nil {
				return err
			}
			checked[u] = true
			sourceChecked[u] = true
			state.CheckedURLs = append(state.CheckedURLs, u)
			if err := s.state.SaveState(ctx, *state); err != nil {
				return err
			}
		}
	}
	return nil
}
